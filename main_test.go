package main

import "testing"

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