package database

import (
	"database/sql"
	"fmt"

	"github.com/noisersup/dashboard-backend-finance/models"
)

type Database struct{
	db 	*sql.DB
}

func ConnectToDatabase(user, password, dbName string) (*Database,error){
	payload := fmt.Sprintf("user=%s password=%s dbname=%s",user,password,dbName)

	db, err := sql.Open("postgres",payload)
	if err != nil {return nil,err}
	
	return &Database{db},nil
}

func (db *Database) CreateExpense(title string, cost float64) (*models.Expense, error) {
	var expense models.Expense
	expense.Title = title
	expense.Cost = cost

	err := db.db.QueryRow("INSERT INTO expenses(title, cost) VALUES ($1,$2) RETURNING ID",
						  title,cost).Scan(&expense.Id)
	if err != nil {return nil,err}

	return &expense,nil
}