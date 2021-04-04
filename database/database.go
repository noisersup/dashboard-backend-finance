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
	
	err = CreateGroupsTable(db)
	if err != nil {return nil,err}
	
	err = CreateExpensesTable(db)
	if err != nil {return nil,err}
	
	return &Database{db},nil
}

func CreateGroupsTable(db *sql.DB) error{
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS groups (
							id				SERIAL,
							title			TEXT NOT NULL UNIQUE,
							maxExpenses		NUMERIC(10,2) NOT NULL,
							currExpenses	NUMERIC(10,2) NOT NULL,
							PRIMARY KEY (id)
						)`)
	return err
}

func CreateExpensesTable(db *sql.DB) error{
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS expenses (
							id			SERIAL,
							group_id	int NOT NULL,
							title		TEXT NOT NULL,
							cost		NUMERIC(10,2) NOT NULL,
							PRIMARY KEY (id),
							FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
						)`)
	return err
}

func (db *Database) CreateExpense(title string,cost float64, group int) (*models.Expense, error) {
	var expense models.Expense
	expense.Title = title
	expense.Cost = cost

	err := db.db.QueryRow("INSERT INTO expenses(title, group_id, cost) VALUES ($1,$2,$3) RETURNING ID",
						  title,group,cost).Scan(&expense.Id)
	if err != nil {return nil,err}

	return &expense,nil
}
