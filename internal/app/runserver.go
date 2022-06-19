package serv

import (
	"dbms/internal/pkg/database"
	"dbms/internal/pkg/middlewares"
	"dbms/internal/pkg/utils/setter"
	"dbms/internal/pkg/utils/log"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Runserver(port string) {
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	api.Use(middlewares.Logger)
	api.Use(middlewares.PanicRecovery)

	db := database.InitDatabase()
	db.Connect()
	defer db.Disconnect()

	setter.SetHandlers(setter.Services{
		Forum: setter.Data{Db: db, Api: api},
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	log.Info("Connecting to port " + port)

	if err := server.ListenAndServe(); err != nil {
		log.Error(err)
	}
}
