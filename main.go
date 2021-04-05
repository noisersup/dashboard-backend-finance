package main

import (
	"flag"
	"log"

	"github.com/noisersup/dashboard-backend-finance/database"
	u "github.com/noisersup/dashboard-backend-finance/utils"
)

func main(){
	userPtr := flag.String("user","","Specify username to postgres database.")
	passwdPtr := flag.String("passwd","","Specify passwd to postgres database.")
	dbPtr := flag.String("database","","Specify name of postgres database.")
	flag.Parse()
	
	if *userPtr==""||*passwdPtr==""||*dbPtr=="" { log.Fatalf("You must specify username, password and name of database!")}


	db,err := database.ConnectToDatabase(*userPtr,*passwdPtr,*dbPtr)
	if err != nil {u.Err("Databse Error",err)}

	
}