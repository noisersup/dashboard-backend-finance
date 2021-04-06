package main

import (
	"log"
	"reflect"
	"testing"

	"github.com/noisersup/dashboard-backend-finance/database"
	"github.com/noisersup/dashboard-backend-finance/models"
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

func TestCreateReadGroups(t *testing.T){
	expectedGroup := models.Group{
		Id: 1,
		Title: "testGroup",
		MaxExpenses: 100.01,
		CurrExpenses: 0.00,
	}

	group, err := test.Db.CreateGroup(expectedGroup.Title,expectedGroup.MaxExpenses)
	if err != nil { t.Fatal(err) }

	exp := getMapFromStruct(expectedGroup)
	score := getMapFromStruct(*group)

	if len(exp) != len(score) {
		t.Errorf("Database record has a different number of columns. \nExpected: %v \nGot: %v",exp,score)
	}

	if expectedGroup.Id != group.Id {t.Errorf("Group id does not match expected \nExpected: %d Got: %d",
	expectedGroup.Id,group.Id)}

	if expectedGroup.Title != group.Title {t.Errorf("Group title does not match expected \nExpected: %s Got: %s",
	expectedGroup.Title,group.Title)}

	if expectedGroup.MaxExpenses != group.MaxExpenses {t.Errorf("Group MaxExpenses does not match expected \nExpected: %f Got: %f",
	expectedGroup.MaxExpenses,group.MaxExpenses)}

	if expectedGroup.CurrExpenses != group.CurrExpenses {t.Errorf("Group CurrExpenses does not match expected \nExpected: %f Got: %f",
	expectedGroup.CurrExpenses,group.CurrExpenses)}
}

func TestCreateReadExpenses(t *testing.T){
	expectedExpense := models.Expense{
		Id: 1,
		Title: "testExpense",
		Cost: 100.01,
	}

	expense, err := test.Db.CreateExpense(expectedExpense.Title,expectedExpense.Cost,1)
	if err != nil { t.Fatal(err) }

	exp := getMapFromStruct(expectedExpense)
	score := getMapFromStruct(*expense)

	if len(exp) != len(score) {
		t.Errorf("Database record has a different number of columns. \nExpected: %v \nGot: %v",exp,score)
	}

	if expectedExpense.Id != expense.Id {t.Errorf("Expense id does not match expected \nExpected: %d Got: %d",
	expectedExpense.Id,expense.Id)}

	if expectedExpense.Title != expense.Title {t.Errorf("Expense title does not match expected \nExpected: %s Got: %s",
	expectedExpense.Title,expense.Title)}

	if expectedExpense.Cost != expense.Cost {t.Errorf("Expense Cost does not match expected \nExpected: %f Got: %f",
	expectedExpense.Cost,expense.Cost)}
}


func getMapFromStruct(inpStruct interface{}) []interface{} {
	val := reflect.ValueOf(inpStruct)
	expValues := make([]interface{}, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		expValues[i] = val.Field(i).Interface()
	}
	return expValues
}
