package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/noisersup/dashboard-backend-finance/models"
)

type Database struct{
	db 	*sql.DB
}

func ConnectToDatabase(host string ,port int,user, password, dbName string) (*Database,error){
	payload := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	host, port,user,password,dbName)
	
	log.Printf("Connecting to database")
	db, err := sql.Open("postgres",payload)
	if err != nil {return nil,err}
	log.Printf("Connected to database")

	log.Printf("Checking for groups table")
	err = CreateGroupsTable(db)
	if err != nil {return nil,err}

	log.Printf("Checking for expenses table")	
	err = CreateExpensesTable(db)
	if err != nil {return nil,err}
	
	return &Database{db},nil
}

func CreateGroupsTable(db *sql.DB) error{
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS groups (
							id				SERIAL,
							title			TEXT NOT NULL UNIQUE,
							max_expenses		NUMERIC(10,2) NOT NULL,
							curr_expenses	NUMERIC(10,2) NOT NULL DEFAULT 0.00,
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

func (db *Database) CreateGroup(title string,maxExpenses float64,) (*models.Group, error) {
	group := models.Group{
		Title: title,
		MaxExpenses: maxExpenses,
	}

	err := db.db.QueryRow("INSERT INTO groups(title, max_expenses) VALUES ($1,$2) RETURNING ID",
						  title,maxExpenses).Scan(&group.Id)

	if err != nil {return nil,err}

	return &group,nil
}

func (db *Database) GetGroups() ([]models.Group,error){
	var groups []models.Group
	
	rows, err := db.db.Query("SELECT * from groups")
	if err != nil {return nil,err}

	defer rows.Close()
	
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.Id,&group.Title,&group.MaxExpenses,&group.CurrExpenses) ; err != nil {
			return nil, err
		}
		
		group.Expenses,err = db.GetExpensesInGroup(group.Id)
		if err != nil { return nil, err }

		groups = append(groups, group)
	}
	return groups, nil
}

func (db *Database) GetExpensesInGroup(id int) ([]models.Expense,error){
	rows, err := db.db.Query("SELECT id,title,cost from expenses WHERE group_id = $1",id)
	if err != nil {return nil,err}
	defer rows.Close()

	var expenses []models.Expense

	for rows.Next() {
		var expense models.Expense
		if err := rows.Scan(&expense.Id,&expense.Title,&expense.Cost) ; err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func (db *Database) ResetDatabase() error {
	if _,err := db.db.Exec("DELETE FROM expenses") ; err != nil { return err}
	if _,err := db.db.Exec("DELETE FROM groups") ; err != nil { return err}
	if _,err := db.db.Exec("ALTER SEQUENCE expenses_id_seq RESTART WITH 1") ; err != nil { return err}
	if _,err := db.db.Exec("ALTER SEQUENCE groups_id_seq RESTART WITH 1") ; err != nil { return err}
	return nil
} 

func (db *Database) CloseDatabase() error{
	return db.db.Close()
}