package controllers

import (
	"opay/api/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	//Postgres lib dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//Server Struct that contains DB GORM and Mux Router
type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

//Initialize method to connecting to spesific DB
func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	server.DB, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", Dbdriver)
	}

	server.DB.Debug().AutoMigrate(&models.Wallet{}) //database migration

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

//Run Method to run Service
func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
