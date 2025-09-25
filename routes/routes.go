package routes

import (
	"e-wallet/controllers"
	"e-wallet/middlewares"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupRoutes(db *mongo.Database) *gin.Engine {

	r := gin.Default()

	// api := r.Group("/api/v1")

	r.POST("/api/v1/register", func(c *gin.Context) { controllers.Register(c, db) })
	r.POST("/api/v1/login", func(c *gin.Context) { controllers.Login(c, db) })

	auth := r.Group("/api/v1")
	auth.Use(middlewares.AuthMiddleware())
	{
		/*
		* Wallet
		 */
		auth.GET("/wallet", func(c *gin.Context) { controllers.GetWallet(c, db) })
		auth.POST("/wallet/deposit", func(c *gin.Context) { controllers.Deposit(c, db) })
		auth.POST("/wallet/withdraw", func(c *gin.Context) { controllers.Withdraw(c, db) })

		/*
		* Transactions
		 */
		auth.POST("/transactions/send", func(c *gin.Context) { controllers.SendMoney(c, db) })
		auth.GET("/transactions/histories", func(c *gin.Context) { controllers.TransactionHistories(c, db) })
		auth.GET("/transactions/history/:transaction_id", func(c *gin.Context) { controllers.TransactionHistory(c, db) })
	}

	return r
}
