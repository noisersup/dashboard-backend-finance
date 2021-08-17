package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/noisersup/dashboard-backend-finance/models"
)

var lastErr string

func LogErr(err error) {
	if err.Error() != lastErr {
		log.Print(err)
		lastErr = err.Error()
	}
}

type Database struct {
	db *sql.DB
}

func ConnectToDatabase(host string, port int, user, password, dbName string, timeout time.Duration) (*Database, error) {
	payload := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	log.Printf("Connecting to database...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan *sql.DB)

	var db *sql.DB

	go func() {
		for {
			db, err := sql.Open("postgres", payload)
			if err != nil {
				log.Print(err)
			} else {
				done <- db
				return
			}
		}
	}()

	select {
	case ndb := <-done:
		db = ndb
		cancel()
		break
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	log.Printf("	Checking for groups table...")
	err := CreateGroupsTable(db, timeout)
	if err != nil {
		return nil, err
	}

	log.Printf("	Checking for expenses table...")
	err = CreateExpensesTable(db, timeout)
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to database!")
	return &Database{db}, nil
}

func CreateGroupsTable(db *sql.DB, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	done := make(chan bool)

	go func() {
		for {
			_, err := db.Exec(`CREATE TABLE IF NOT EXISTS groups (
							id				SERIAL,
							title			TEXT NOT NULL UNIQUE,
							max_expenses		NUMERIC(10,2) NOT NULL,
							curr_expenses	NUMERIC(10,2) NOT NULL DEFAULT 0.00,
							PRIMARY KEY (id)
						)`)
			if err != nil {
				LogErr(err)
			} else {
				done <- true
				return
			}
		}
	}()
	select {
	case <-done:
		cancel()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func CreateExpensesTable(db *sql.DB, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	done := make(chan bool)

	go func() {
		_, err := db.Exec(`CREATE TABLE IF NOT EXISTS expenses (
							id			SERIAL,
							group_id	int NOT NULL,
							title		TEXT NOT NULL,
							cost		NUMERIC(10,2) NOT NULL,
							PRIMARY KEY (id),
							FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
						)`)
		if err != nil {
			LogErr(err)
		} else {
			done <- true
			return
		}

	}()
	select {
	case <-done:
		cancel()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (db *Database) CreateExpense(title string, cost float64, group int, timeout time.Duration) (*models.Expense, error) {

	if title == "" || cost <= 0 || group == 0 {
		return nil, errors.New("CreateExpense: title is empty, cost value is not valid or group is not specified.")
	}
	var expense models.Expense
	expense.Title = title
	expense.Cost = cost

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	done := make(chan bool)

	go func() {
		for {
			err := db.db.QueryRow("INSERT INTO expenses(title, group_id, cost) VALUES ($1,$2,$3) RETURNING ID",
				title, group, cost).Scan(&expense.Id)
			if err != nil {
				LogErr(err)
			} else {
				done <- true
				return
			}
		}
	}()

	select {
	case <-done:
		cancel()
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return &expense, nil
}

func (db *Database) CreateGroup(title string, maxExpenses float64) (*models.Group, error) {
	if title == "" || maxExpenses <= 0 {
		return nil, errors.New("CreateGroup: title is empty or max expenses value is not valid")
	}
	group := models.Group{
		Title:       title,
		MaxExpenses: maxExpenses,
	}
	err := db.db.QueryRow("INSERT INTO groups(title, max_expenses) VALUES ($1,$2) RETURNING ID",
		title, maxExpenses).Scan(&group.Id)

	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (db *Database) GetGroups() ([]models.Group, error) {
	var groups []models.Group

	rows, err := db.db.Query("SELECT * from groups")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.Id, &group.Title, &group.MaxExpenses, &group.CurrExpenses); err != nil {
			return nil, err
		}

		group.Expenses, err = db.GetExpensesInGroup(group.Id)
		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}
	return groups, nil
}

func (db *Database) GetExpensesInGroup(id int) ([]models.Expense, error) {
	rows, err := db.db.Query("SELECT id,title,cost from expenses WHERE group_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense

	for rows.Next() {
		var expense models.Expense
		if err := rows.Scan(&expense.Id, &expense.Title, &expense.Cost); err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func (db *Database) ResetDatabase() error {
	if _, err := db.db.Exec("DELETE FROM expenses"); err != nil {
		return err
	}
	if _, err := db.db.Exec("DELETE FROM groups"); err != nil {
		return err
	}
	if _, err := db.db.Exec("ALTER SEQUENCE expenses_id_seq RESTART WITH 1"); err != nil {
		return err
	}
	if _, err := db.db.Exec("ALTER SEQUENCE groups_id_seq RESTART WITH 1"); err != nil {
		return err
	}
	return nil
}

func (db *Database) CloseDatabase() error {
	return db.db.Close()
}
