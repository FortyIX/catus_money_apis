package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Transaction struct {
	Id       string `json:"id" form:"id""`
	DateInfo string `json:"date" form:"date"`
	Amount   int    `json:"amount" form:"amount"`
	Type     string `json:"type" form:"type"`
	Account  string `json:"account" form:"account"`
	User 	 string `json:"user" form:"user"`
	Note     string `json:"note" form:"note"`
}


type BankAccount struct {
	Id string `json:"id" form:"id""`
	Name string `json:"name" form:"name""`
	Holder string `json:"holder" form:"holder""`
	Balance int `json:"balance" form:"balance""`

}

var db *sql.DB

func main() {

	initDbConnection()
	router := gin.Default()


	/// Set heaeders fro CROS
	router.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))



	///Query all transactions
	router.GET("/transactions/query", func(c *gin.Context) {

		res, _ := queryAllTransaction()
		c.JSON(http.StatusOK, gin.H{
			"transactions": res,
		})
	})

	/// Add  a new transaction
	router.POST("/transactions/add", func(c *gin.Context) {

		id := c.Request.FormValue("id")
		dateinfo := c.PostForm("date")
		amount := c.PostForm("amount")
		typeinfo := c.PostForm("type")
		account := c.PostForm("account")
		user := c.PostForm("user")
		note := c.PostForm("note")

		amount_int, _ := strconv.Atoi(amount)

		transaction := Transaction{
			Id:       id,
			DateInfo: dateinfo,
			Amount:   amount_int,
			Type:     typeinfo,
			Account:  account,
			User:	  user,
			Note:     note,
		}

		lastID := addTransaction(transaction)
		msg := fmt.Sprintf("Insert successfully %d", lastID)
		c.JSON(http.StatusOK, gin.H{
			"status": msg,
		})

	})

	router.POST("/transactions/del", func(c *gin.Context){
		id := c.Request.FormValue("id")

		res := delTransaction(id);
		if res != 0 {
			c.JSON(-1,gin.H{
				"status": "deletion failed",
			})
		}else{
			c.JSON(http.StatusOK,gin.H{
				"status":"successfully deleted",
			})
		}
	})


	router.GET("bankAccount/query",func(c *gin.Context){
		res := queryAllBankAccount()
		c.JSON(http.StatusOK,gin.H{
			"bank_account": res,
		})
	})

	router.POST("bankAccount/addAccount",func(c *gin.Context){

		id := c.Request.FormValue("id");
		name := c.Request.FormValue("account_name");
		balance := c.Request.FormValue("account_balance");
		holder := c.Request.FormValue("holder");

		balance_int,_ := strconv.Atoi(balance)

		newAccount := BankAccount{
			Name: name,
			Balance: balance_int,
			Holder: holder,
			Id: id,
		}


		lastid := addNewBankAccount(newAccount);
		msg := fmt.Sprintf("Insert successfully %d", lastid)
		c.JSON(http.StatusOK, gin.H{
			"status": msg,
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
		data.Scan(&transaction.Id, &transaction.DateInfo, &transaction.Amount, &transaction.Type, &transaction.Account, &transaction.User,&transaction.Note)
		fmt.Println(transaction)
		transactions = append(transactions, transaction)
	}

	data.Close()
	return
}

func addTransaction(transaction Transaction) int64 {
	data, err := db.Exec("INSERT INTO `transaction`(`id`, `date`, `amount`, `type`, `account`, `user`, `note`) VALUES (?,?,?,?,?,?,?)",
		transaction.Id, transaction.DateInfo, transaction.Amount, transaction.Type, transaction.Account, transaction.User,transaction.Note)
	if err != nil {
		log.Fatalln(err)
	}

	id, err := data.LastInsertId()
	if err != nil {
		log.Fatalln(err)
	}
	return id
}


func delTransaction(id string) int64 {
	_,err := db.Exec("DELETE FROM `transaction` WHERE id=?",id)
	if err != nil {
		log.Fatalln(err)
		return -1
	}
	return 0
}


func queryAllBankAccount() (accounts []BankAccount) {
	data, err := db.Query("SELECT * FROM `bank_account` WHERE 1")
	if err != nil {
		log.Fatalln(err)
	}

	for data.Next() {
		account := BankAccount{}
		data.Scan(&account.Name,&account.Holder,&account.Balance,&account.Id)
		accounts = append(accounts,account)

	}

	return

}


func addNewBankAccount(newAccount BankAccount) int64 {
	data,err := db.Exec("INSERT INTO `bank_account`(`account_name`, `account_holder`, `account_balance`, `account_id`) VALUES (?,?,?,?)",
		newAccount.Name,newAccount.Holder,newAccount.Balance,newAccount.Id)

	if err != nil {
		log.Fatalln(err)
	}

	id,err := data.LastInsertId()
	if err != nil {
		log.Fatalln(err)
	}

	return id
}

