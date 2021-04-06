package main

import (
	"log"
	"testing"

	"github.com/noisersup/dashboard-backend-finance/database"
)

type TestVars struct{
	Db	*database.Database
}

var test TestVars
func TestMain(m *testing.M) {
	var config DbConfig
	
	config,err := loadConfig("config/database.json") 
	if err != nil {log.Fatal(err)}
 

	test.Db,err = database.ConnectToDatabase(config.Username,config.Password,config.Database)
	if err != nil {log.Fatal(err)}

	err = test.Db.ResetDatabase()
	if err != nil {log.Fatal(err)}
	
	m.Run()
}

func TestLoadConfig(t *testing.T){
	var config DbConfig
	config,err := loadConfig("config/tests/database.json")
	if err != nil { t.Fatal(err) }

	if config.Database != "databaseTest" || 
	   config.Password != "passwordTest" || 
	   config.Username != "usernameTest" {
		   t.Fatalf(`Data provided by test file invalid: ("%s", "%s", "%s")`,
		   			config.Username,config.Database,config.Password)
	}
}