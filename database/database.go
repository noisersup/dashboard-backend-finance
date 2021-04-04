package database

import (
	"database/sql"
	"fmt"
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