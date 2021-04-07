package handlers

import (
	"net/http"

	"github.com/noisersup/dashboard-backend-finance/database"
)

type Server struct{
	db 	*database.Database
}

func CreateHandlers(db *database.Database) Server{
	return Server{db}
}

func (srv *Server) GetGroups(w http.ResponseWriter, r *http.Request){
	
}