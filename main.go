package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/noisersup/dashboard-backend-finance/api"
	"github.com/noisersup/dashboard-backend-finance/api/pb"
	"github.com/noisersup/dashboard-backend-finance/database"

	u "github.com/noisersup/dashboard-backend-finance/utils"
	"google.golang.org/grpc"
)

type DbConfig struct{
	Address		string `json:"address"`
	Port		int    `json:"port"`
	Username	string `json:"username"`
	Password	string `json:"password"`
	Database	string `json:"database"`
}

func main(){
	addrPtr := flag.String("address","127.0.0.1","Specify ip address of postgres database.")
	portPtr := flag.Int("port",5432,"Specify port of postgres database.")
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
		config.Address = *addrPtr
		config.Port = *portPtr
		config.Database = *dbPtr
		config.Username = *userPtr
		config.Password = *passwdPtr
	}

	

	db,err := database.ConnectToDatabase(config.Address,config.Port,config.Username,config.Password,config.Database)
	if err != nil {log.Fatal(u.Err("Database Error",err))}

	defer func ()  {
		if err := db.CloseDatabase(); err != nil{
			log.Fatal(u.Err("Closing database error",err))
		}
	}()


	lis, err := net.Listen("tcp",":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v",err)
	}
	log.Print(lis.Addr())
	a := api.Api{}

	server := grpc.NewServer()
	pb.RegisterFinanceServiceServer(server,&a)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v",err)
	}
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