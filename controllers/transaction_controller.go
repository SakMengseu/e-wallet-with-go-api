package controllers

import (
	"e-wallet/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SendMoney(c *gin.Context, db *mongo.Database) {
	userId := c.MustGet("user_id").(string)

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var input struct {
		ReceiverEmail string  `json:"receiver_email" binding:"required,email"`
		Amount        float64 `json:"amount" binding:"required,gt=0"`
	}

	err := c.ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
		return
	}

	walletColl := db.Collection("wallets")
	txColl := db.Collection("transactions")
	userColl := db.Collection("users")

	// Find sender and receiver wallets
	var senderWallet models.Wallet
	var senderUser models.User
	err = userColl.FindOne(c, bson.M{"_id": userId}).Decode(&senderUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sender user not found"})
		return
	}

	err = walletColl.FindOne(c, bson.M{"user_id": userId}).Decode(&senderWallet)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sender wallet not found"})
		return
	}

	if senderWallet.Balance < input.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	// Find receiver wallet by email
	var receiverWallet models.Wallet
	var receiverUser models.User

	err = userColl.FindOne(c, bson.M{"email": input.ReceiverEmail}).Decode(&receiverUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receiver user not found"})
		return
	}

	//print is receiver user id
	fmt.Println("Receiver User ID:", receiverUser.ID)
	err = walletColl.FindOne(c, bson.M{"user_id": receiverUser.ID}).Decode(&receiverWallet)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receiver wallet not found"})
		return
	}
	// Update balances
	senderWallet.Balance -= input.Amount

	receiverWallet.Balance += input.Amount

	_, err = walletColl.UpdateOne(c, bson.M{"user_id": userId}, bson.M{"$set": bson.M{"balance": senderWallet.Balance}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sender wallet"})
		return
	}

	_, err = walletColl.UpdateOne(c, bson.M{"user_id": receiverWallet.ID}, bson.M{"$set": bson.M{"balance": receiverWallet.Balance}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update receiver wallet"})
		return
	}
	//find sender and receiver user

	// Log transactions for both sender and receiver
	newTxSender := models.Transaction{
		ID:         uuid.NewString(),
		UserID:     userId,
		SenderID:   userId,
		ReceiverID: receiverUser.ID,
		Amount:     input.Amount,
		Type:       "send",
		Status:     "success",
		CreatedAt:  time.Now(),
	}

	_, err = txColl.InsertOne(c, newTxSender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert transaction"})
		return
	}

	newTxReceiver := models.Transaction{
		ID:         uuid.NewString(),
		UserID:     receiverUser.ID,
		SenderID:   userId,
		ReceiverID: receiverUser.ID,
		Amount:     input.Amount,
		Type:       "receive",
		Status:     "success",
		CreatedAt:  time.Now(),
	}

	_, err = txColl.InsertOne(c, newTxReceiver)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Money sent successfully",
		"sender": gin.H{
			"id":      senderUser.ID,
			"name":    senderUser.Name,
			"email":   senderUser.Email,
			"balance": input.Amount,
			"date":    time.Now(),
		},
		"receiver": gin.H{
			"id":      receiverUser.ID,
			"name":    receiverUser.Name,
			"email":   receiverUser.Email,
			"balance": input.Amount,
			"date":    time.Now(),
		},
	})
}

func TransactionHistories(c *gin.Context, db *mongo.Database) {
	userId := c.MustGet("user_id").(string)

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	txColl := db.Collection("transactions")

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "user_id", Value: userId},
			{Key: "$or", Value: bson.A{
				bson.D{{Key: "sender_id", Value: userId}},
				bson.D{{Key: "receiver_id", Value: userId}},
			}},
		}}},

		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},           // target collection
			{Key: "localField", Value: "sender_id"}, // field in transactions
			{Key: "foreignField", Value: "_id"},     // field in users
			{Key: "as", Value: "sender"},            // output array
		}}},
		
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},             // target collection
			{Key: "localField", Value: "receiver_id"}, // field in transactions
			{Key: "foreignField", Value: "_id"},       // field in users
			{Key: "as", Value: "receiver"},            // output array
		}}},
		{{Key: "$unwind", Value: "$sender"}},   // flatten array into object
		{{Key: "$unwind", Value: "$receiver"}}, // flatten array into object
		{{Key: "$project", Value: bson.D{
			{Key: "amount", Value: 1},
			{Key: "type", Value: 1},
			{Key: "status", Value: 1},
			{Key: "created_at", Value: 1},
			{Key: "sender.name", Value: 1},
			{Key: "sender.email", Value: 1},
			{Key: "receiver.name", Value: 1},
			{Key: "receiver.email", Value: 1},
		}}},
	}

	cursor, err := txColl.Aggregate(c, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction history"})
		return
	}
	defer cursor.Close(c)

	var result []bson.M

	err = cursor.All(c, &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode transaction history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"history": result,
	})
}

func TransactionHistory(c *gin.Context, db *mongo.Database) {
	userID := c.MustGet("user_id").(string)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	txID := c.Param("transaction_id")
	if txID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	txColl := db.Collection("transactions")

	pipeline := mongo.Pipeline{
		// Match transaction by ID and where user is sender or receiver
		{{Key: "$match", Value: bson.D{
			{Key: "_id", Value: txID},
			{Key: "$or", Value: bson.A{
				bson.D{{Key: "sender_id", Value: userID}},
				bson.D{{Key: "receiver_id", Value: userID}},
			}},
		}}},

		// Lookup sender info
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "sender_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "sender"},
		}}},

		// Lookup receiver info
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "receiver_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "receiver"},
		}}},

		// Unwind sender and receiver arrays
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$sender"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$receiver"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},

		// Project only needed fields
		{{Key: "$project", Value: bson.D{
			{Key: "amount", Value: 1},
			{Key: "type", Value: 1},
			{Key: "status", Value: 1},
			{Key: "created_at", Value: 1},
			{Key: "sender.name", Value: 1},
			{Key: "sender.email", Value: 1},
			{Key: "receiver.name", Value: 1},
			{Key: "receiver.email", Value: 1},
		}}},
	}

	cursor, err := txColl.Aggregate(c, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction history"})
		return
	}
	defer cursor.Close(c)

	var result []bson.M
	if err := cursor.All(c, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode transaction history"})
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "success",
		"transaction": result[0], // single transaction
	})
}
