package controllers

import (
	"e-wallet/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetWallet(c *gin.Context, db *mongo.Database) {
	userId := c.MustGet("user_id").(string)

	walletColl := db.Collection("wallets")
	var wallet models.Wallet
	err := walletColl.FindOne(c, bson.M{"user_id": userId}).Decode(&wallet)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"wallet":  wallet,
	})
}

func Deposit(c *gin.Context, db *mongo.Database) {
	userId := c.MustGet("user_id").(string)

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var input struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
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

	var wallet models.Wallet
	err = walletColl.FindOne(c, bson.M{"user_id": userId}).Decode(&wallet)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	wallet.Balance += input.Amount
	_, err = walletColl.UpdateOne(c, bson.M{"user_id": userId}, bson.M{"$set": bson.M{"balance": wallet.Balance}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet"})
		return
	}

	newTx := models.Transaction{
		ID:        uuid.NewString(),
		SenderID:  userId,
		UserID:    userId,
		Amount:    input.Amount,
		Type:      "deposit",
		Status:    "success",
		CreatedAt: time.Now(),
	}
	_, err = txColl.InsertOne(c, newTx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Deposit successful",
		"wallet":  wallet.Balance,
	})

}

func Withdraw(c *gin.Context, db *mongo.Database) {
	userId := c.MustGet("user_id").(string)

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var input struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
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

	var wallet models.Wallet
	err = walletColl.FindOne(c, bson.M{"user_id": userId}).Decode(&wallet)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	if wallet.Balance < input.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	wallet.Balance -= input.Amount
	_, err = walletColl.UpdateOne(c, bson.M{"user_id": userId}, bson.M{"$set": bson.M{"balance": wallet.Balance}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet"})
		return
	}

	newTx := models.Transaction{
		ID:        uuid.NewString(),
		UserID:    userId,
		SenderID:  uuid.NewString(),
		Amount:    input.Amount,
		Type:      "withdraw",
		Status:    "success",
		CreatedAt: time.Now(),
	}
	_, err = txColl.InsertOne(c, newTx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Withdraw successful",
		"wallet":  wallet.Balance,
	})

}
