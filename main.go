package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/noisersup/dashboard-backend-finance/database"
	u "github.com/noisersup/dashboard-backend-finance/utils"
)

type DbConfig struct{
	Address		string `json:"address"`
	Port		int    `json:"port"`
	Username	string `json:"username"`
	Password	string `json:"password"`
	Database	string `json:"database"`
}

func main(){
	config := getVars()
	
	db,err := database.ConnectToDatabase(config.Address,config.Port,config.Username,config.Password,config.Database)
	if err != nil {log.Fatal(u.Err("Database Error",err))}
	
	defer func ()  {
		if err := db.CloseDatabase(); err != nil{
			log.Fatal(u.Err("Closing database error",err))
		}
	}()
}

func getVars() *DbConfig {	
	var config DbConfig

	config.Address = os.Getenv("DB_ADDRESS")
	config.Database = os.Getenv("DB_NAME")
	config.Username = os.Getenv("DB_USERNAME")
	config.Password = os.Getenv("DB_PASSWORD")

	config.Port,_ = strconv.Atoi(os.Getenv("DB_PORT")) //default: 5432

	if config.Address=="" || config.Database==""||config.Username==""||config.Password==""{ 
		log.Printf("ENV variables not set, loading config from config file")
		
		var err error
		config, err = loadConfig("config/database.json")
		if err != nil {log.Fatal(u.Err("Config file error",err))}
	}

	if config.Port == 0 {
		log.Print("Port invalid or not provided. Setting default (5432)")
		config.Port=5432
	}
	return &config
}

func loadConfig(jsonFile string) (DbConfig, error){
	file, err := os.Open(jsonFile)
	if err != nil {return DbConfig{},err}
	
	defer file.Close()
	
	bytes, err := ioutil.ReadAll(file)
	if err != nil {return DbConfig{},err}
	
	var config DbConfig

	json.Unmarshal([]byte(bytes),&config)

	return config, nil
}