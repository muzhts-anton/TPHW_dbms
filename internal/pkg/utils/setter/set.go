package setter

import (
	"dbms/internal/pkg/database"

	_ "dbms/internal/pkg/forum/delivery"
	_ "dbms/internal/pkg/forum/repository"
	_ "dbms/internal/pkg/forum/usecase"

	"github.com/gorilla/mux"
)

type Data struct {
	Db  *database.DBManager
	Api *mux.Router
}

type Services struct {
	Forum Data
}

func SetHandlers(svs Services) {
	// rep := rep.InitRep(svs.Forum.Db)
	// usc := usc.InitUsc(rep)
	// del.SetHandlers(svs.Forum.Api, usc)
}
