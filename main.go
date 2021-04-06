package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/noisersup/dashboard-backend-finance/database"
	u "github.com/noisersup/dashboard-backend-finance/utils"
)

type DbConfig struct{
	Username	string `json:"username"`
	Password	string `json:"password"`
	Database	string `json:"database"`
}

func main(){
	userPtr := flag.String("user","","Specify username to postgres database.")
	passwdPtr := flag.String("passwd","","Specify passwd to postgres database.")
	dbPtr := flag.String("database","","Specify name of postgres database.")
	flag.Parse()
	
	var config DbConfig
	var err error
	if *userPtr==""||*passwdPtr==""||*dbPtr=="" { 
		log.Printf("Authentication flags not set, loading passes from config file")
		config, err = loadConfig("config/database.json")
		if err != nil {log.Fatal(u.Err("Config file error",err))}
	}else{
		config.Database = *dbPtr
		config.Username = *userPtr
		config.Password = *passwdPtr
	}
	
	_,err = database.ConnectToDatabase(config.Username,config.Password,config.Database)
	if err != nil {log.Fatal(u.Err("Database Error",err))}
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