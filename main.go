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
	"strings"
	"time"
)

type Transaction struct {
	Id       string `json:"id" form:"id""`
	DateInfo string `json:"date" form:"date"`
	Amount   int    `json:"amount" form:"amount"`
	Type     string `json:"type" form:"type"`
	Account  string `json:"account" form:"account"`
	User     string `json:"user" form:"user"`
	Note     string `json:"note" form:"note"`
}

type BankAccount struct {
	Id      string `json:"id" form:"id""`
	Name    string `json:"name" form:"name""`
	Holder  string `json:"holder" form:"holder""`
	Balance int    `json:"balance" form:"balance""`
}

type Trends struct {
	MonthlySpending     string `json:"month" form:"month""`
	HalfMonthlyTrending string `json:"day" form:"day""`
}

type Users struct {
	Username  string `json:"username" form:"username"`
	Password  string `json:"password" form:"password"`
	Role      string `json:"role" form:"role"`
	InviteKey string `json:"inviteKey" form:"inviteKey"`
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

	router.GET("account/getUserAmount", func(c *gin.Context) {

		res := queryNumOfUsers()
		c.JSON(http.StatusOK, gin.H{
			"nUser": res,
		})
	})

	router.POST("account/adminRegister", func(c *gin.Context) {
		username := c.Request.FormValue("username")
		password := c.Request.FormValue("password")

		role := "admin"
		inviteKey := "null"

		user := Users{
			Username:  username,
			Password:  password,
			Role:      role,
			InviteKey: inviteKey,
		}

		res := registerAccount(user)

		c.JSON(http.StatusOK, gin.H{
			"admin_registeration": res,
		})
	})

	router.POST("account/userRegister", func(c *gin.Context) {
		username := c.Request.FormValue("username")
		password := c.Request.FormValue("password")
		role := "user"
		inviteKey := c.Request.FormValue("invitekey")
		fmt.Println(password)

		user := Users{
			Username:  username,
			Password:  password,
			Role:      role,
			InviteKey: inviteKey,
		}

		res := registerForNormalUsers(user)

		c.JSON(http.StatusOK, gin.H{
			"result": res,
		})
	})

	router.POST("account/verifyAccount", func(c *gin.Context) {
		username := c.Request.FormValue("username")
		password := c.Request.FormValue("password")

		user := Users{
			Username:  username,
			Password:  password,
			Role:      "",
			InviteKey: "",
		}

		res := verifyLoginInfo(user)

		c.JSON(http.StatusOK, gin.H{
			"result": res,
		})

	})

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
		dateinfo := c.Request.FormValue("date")
		amount := c.Request.FormValue("amount")
		typeinfo := c.Request.FormValue("type")
		account := c.Request.FormValue("account")
		user := c.Request.FormValue("user")
		note := c.Request.FormValue("note")

		amount_int, _ := strconv.Atoi(amount)

		transaction := Transaction{
			Id:       id,
			DateInfo: dateinfo,
			Amount:   amount_int,
			Type:     typeinfo,
			Account:  account,
			User:     user,
			Note:     note,
		}

		lastID := addTransaction(transaction)
		msg := fmt.Sprintf("Insert successfully %d", lastID)
		c.JSON(http.StatusOK, gin.H{
			"status": msg,
		})

	})

	router.POST("/transactions/del", func(c *gin.Context) {
		id := c.Request.FormValue("id")
		token := c.Request.FormValue("token")
		isTokenVaild := verifyToken(token)

		if isTokenVaild {
			res := delTransaction(id)
			if res != 0 {
				c.JSON(-1, gin.H{
					"status": "deletion failed",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status": "successfully deleted",
				})
			}
		}else{
			c.JSON(http.StatusForbidden, gin.H{
				"Forbidden": "Invaild Token",
			})
		}
	})

	router.GET("bankAccount/query", func(c *gin.Context) {

		token := c.Query("token")
		isTokenVaild := verifyToken(token)

		if isTokenVaild {
			res := queryAllBankAccount()
			c.JSON(http.StatusOK, gin.H{
				"bank_account": res,
			})
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"Forbidden": "Invaild Token",
			})
		}

	})

	router.POST("bankAccount/addAccount", func(c *gin.Context) {

		id := c.Request.FormValue("id")
		name := c.Request.FormValue("account_name")
		balance := c.Request.FormValue("account_balance")
		holder := c.Request.FormValue("holder")

		balance_int, _ := strconv.Atoi(balance)

		newAccount := BankAccount{
			Name:    name,
			Balance: balance_int,
			Holder:  holder,
			Id:      id,
		}

		lastid := addNewBankAccount(newAccount)
		msg := fmt.Sprintf("Insert successfully %d", lastid)
		c.JSON(http.StatusOK, gin.H{
			"status": msg,
		})

	})

	router.GET("stats/trends", func(c *gin.Context) {
		res := queryAllTrends()
		c.JSON(http.StatusOK, gin.H{
			"trends": res,
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

	db, err = sql.Open("mysql", "juizeffs:root@tcp(127.0.0.1)/juizeffs_catusMoney")
	if err != nil {
		log.Fatalln(err)
	}

}

func queryNumOfUsers() int {
	data, err := db.Query("SELECT COUNT(*) FROM `user` WHERE 1 ")
	if err != nil {
		log.Fatalln(err)
	}

	counter := 0
	for data.Next() {
		data.Scan(&counter)
	}

	return counter

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
		data.Scan(&transaction.Id, &transaction.DateInfo, &transaction.Amount, &transaction.Type, &transaction.Account, &transaction.User, &transaction.Note)
		fmt.Println(transaction)
		transactions = append(transactions, transaction)
	}

	data.Close()
	return
}

func addTransaction(transaction Transaction) int64 {
	data, err := db.Exec("INSERT INTO `transaction`(`id`, `date`, `amount`, `type`, `account`, `user`, `note`) VALUES (?,?,?,?,?,?,?)",
		transaction.Id, transaction.DateInfo, transaction.Amount, transaction.Type, transaction.Account, transaction.User, transaction.Note)
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
	_, err := db.Exec("DELETE FROM `transaction` WHERE id=?", id)
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
		data.Scan(&account.Name, &account.Holder, &account.Balance, &account.Id)
		accounts = append(accounts, account)

	}

	return

}

func addNewBankAccount(newAccount BankAccount) int64 {
	data, err := db.Exec("INSERT INTO `bank_account`(`account_name`, `account_holder`, `account_balance`, `account_id`) VALUES (?,?,?,?)",
		newAccount.Name, newAccount.Holder, newAccount.Balance, newAccount.Id)

	if err != nil {
		log.Fatalln(err)
	}

	id, err := data.LastInsertId()
	if err != nil {
		log.Fatalln(err)
	}

	return id
}

func queryAllTrends() (trendsData []Trends) {
	data, err := db.Query("SELECT * FROM `trend` WHERE 1")
	if err != nil {
		log.Fatalln(err)
	}

	for data.Next() {
		trend := Trends{}
		data.Scan(&trend.MonthlySpending, &trend.HalfMonthlyTrending)
		trendsData = append(trendsData, trend)

	}

	return
}

func registerForNormalUsers(user Users) int {
	data, err := db.Query("SELECT `password` FROM `user` WHERE `username`= 'admin'")
	admin_token := ""
	result := 0
	if err != nil {
		log.Fatalln(err)
	}

	for data.Next() {
		data.Scan(&admin_token)
	}

	if strings.Compare(user.InviteKey, admin_token) == 0 {
		result = 1
		registerAccount(user)
	}
	return result

}

func verifyToken(token string) bool {
	data, err := db.Query("SELECT COUNT(`password`) FROM `user` WHERE `password`= ?", token)

	counter := 0
	result := false
	if err != nil {
		log.Fatalln(err)
	}

	for data.Next() {
		data.Scan(&counter)
	}

	if counter > 0 {
		result = true
	}

	return result

}

func registerAccount(user Users) int64 {
	data, err := db.Exec("INSERT INTO `user`(`username`, `password`, `role`) VALUES (?,?,?)", user.Username, user.Password, user.Role)
	if err != nil {
		fmt.Println(err)
	}
	id, err := data.LastInsertId()
	if err != nil {
		log.Fatalln(err)
	}
	return id

}

func verifyLoginInfo(user Users) int {
	data, err := db.Query("SELECT `password` FROM `user` WHERE `username`= ?", user.Username)
	if err != nil {
		log.Fatalln(err)
	}

	password := ""
	result := 0
	for data.Next() {
		data.Scan(&password)
	}

	if strings.Compare(user.Password, password) == 0 {
		result = 1
	}
	return result
}
