package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

type Transaction struct {
	Id       string `json:"id" form:"id""`
	DateInfo string `json:"data" form:"date"`
	Amount   int    `json:"amount" form:"amount"`
	Type     string `json:"type" form:"type"`
	Account  string `json:"account" form:"account"`
	Note     string `json:"note" form:"note"`
}

var db *sql.DB

func main() {

	initDbConnection()
	router := gin.Default()

	router.GET("/transactions", func(c *gin.Context) {
		res, _ := queryAllTransaction()
		c.JSON(http.StatusOK, gin.H{
			"transactions": res,
		})
	})

	err := router.Run(":3990")
	if err != nil {
		return
	}
}

/**
initialize the database connection
*/
func initDbConnection() {

	var err error

	db, err = sql.Open("mysql", "juizeffs:dzQbQ473f0@tcp(43.240.31.71)/juizeffs_catusMoney")
	if err != nil {
		log.Fatalln(err)
	}

}

/**
  Function that get all entries in the transaction table
*/
func queryAllTransaction() (transactions []Transaction, errMsg error) {

	data, err := db.Query("SELECT * FROM `transaction` WHERE 1")
	if err != nil {
		log.Fatalln(err)
	}
	for data.Next() {
		transaction := Transaction{}
		data.Scan(&transaction.Id, &transaction.DateInfo, &transaction.Amount, &transaction.Type, &transaction.Account, &transaction.Note)
		fmt.Println(transaction)
		transactions = append(transactions, transaction)
	}

	data.Close()
	return
}
