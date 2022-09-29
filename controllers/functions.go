package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	orm "github.com/go-pg/pg/v9/orm"
)

type Bank struct {
	TableName struct{} `json:"table_name" pg:"bank"`
	BankID    int      `json:"bank_id" pg:"bank_id,pk"`
	BankName  string   `json:"bank_name"`
}

type Acc struct {
	TableName     struct{} `json:"table_name" pg:"accs"`
	Name          string   `json:"name"`
	AccountID     int64    `json:"account_id" pg:"account_id,pk, type:serial"`
	AccountNumber int64    `json:"account_number" pg:"account_number,type:serial"`
	BankID        int      `json:"bank_id" pg:"bank_id,type:serial references bank(bank_id) on delete set NULL on update cascade"`
	CustomerID    int64    `json:"customer_id" pg:"customer_id,type:int references customers(customer_id) on delete set NULL on update cascade"`
	Balance       int64    `json:"balance"`
}

type credit string

type debit string

type TransType struct {
	Savings credit
	Current debit
}

type Transactions struct {
	TableName       struct{}  `json:"transactions" pg:"transactions"`
	TransactionID   int64     `json:"transaction_id"`
	Amount          int64     `json:"amount"`
	TransactionType TransType `json:"transaction_type"`
	Date            time.Time `json:"date"`
	AccountID       int       `json:"account_id" pg:"account_id,type:int references accs(account_id) on delete set NULL on update cascade"`
}

type Customers struct {
	TableName  struct{} `json:"table_name" pg:"customers"`
	CustomerID int64    `json:"customer_id" pg:"customer_id,pk"`
	Name       string   `json:"name"`
	BankID     int      `json:"bank_id" pg:"bank_id,type:serial references bank(bank_id) on delete set NULL on update cascade"`
}

// Create User Table, WORKING
func CreateAccountTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		IfNotExists: true,
	}
	createError := db.CreateTable(&Acc{}, opts)
	if createError != nil {
		log.Printf("Error while creating Accounts table, Reason: %v\n", createError)
		return createError
	}
	log.Printf("Accounts table created")
	return nil
}

func CreateCustomerTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		IfNotExists: true,
	}
	createError := db.CreateTable(&Customers{}, opts)
	if createError != nil {
		log.Printf("Error while creating Accounts table, Reason: %v\n", createError)
		return createError
	}
	log.Printf("Accounts table created")
	return nil
}

// INITIALIZE DB CONNECTION (TO AVOID TOO MANY CONNECTION)
var dbConnect *pg.DB

func InitiateDB(db *pg.DB) {
	dbConnect = db
}

func GetAllAccounts(c *gin.Context) {
	var temps []Acc
	err := dbConnect.Model(&temps).Select()

	if err != nil {
		log.Printf("Error while getting all Accounts, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Accounts",
		"data":    temps,
	})
	return
}

func CreateAccount(c *gin.Context) {

	var temp Acc
	c.BindJSON(&temp)
	bank_id := temp.BankID
	name := temp.Name

	customer_id := temp.CustomerID
	balance := temp.Balance

	insertError := dbConnect.Insert(&Acc{
		TableName: struct{}{},
		Name:      name,
		//AccountID:     account_id,
		//AccountNumber: account_number,
		BankID:     bank_id,
		CustomerID: customer_id,
		Balance:    balance,
	})
	if insertError != nil {
		log.Printf("Error while inserting new account into db, Reason: %v\n", insertError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusOK,
		"message": "Account created Successfully",
	})
	return
}

// Function to create customer
func CreateCustomer(c *gin.Context) {
	var temp Customers
	c.BindJSON(&temp)
	customer_id := temp.CustomerID
	name := temp.Name
	bank_id := temp.BankID

	insertError := dbConnect.Insert(&Customers{
		TableName:  struct{}{},
		CustomerID: customer_id,
		Name:       name,
		BankID:     bank_id,
	})
	if insertError != nil {
		log.Printf("Error while inserting new account into db, Reason: %v\n", insertError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Customer created Successfully",
	})
	return
}

func GetSingleAccount(c *gin.Context) {
	varName := c.Param("varId")
	temp := &Acc{Name: varName}
	err := dbConnect.Select(temp)

	if err != nil {
		log.Printf("Error while getting a single Account, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Account not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single Account",
		"data":    temp,
	})
	return
}


func EditAccount(c *gin.Context) {

	accountID := c.Query("id")
	var updateAccount Acc
	c.BindJSON(&updateAccount)
	_, err := dbConnect.Model(&updateAccount).Where("account_id=?", accountID).Update()

	if err != nil {
		log.Printf("Error, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Account Edited Successfully",
	})
	return
}

func DeleteAccount(c *gin.Context) {

	accountID := c.Query("id")
	var deleteAccount Acc
	_, err := dbConnect.Model(&deleteAccount).Where("account_id=?", accountID).Delete()

	if err != nil {
		log.Printf("Error while deleting a single account, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Account deleted successfully",
	})
	return
}
