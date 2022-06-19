package del

import (
	"dbms/internal/pkg/domain"

	"net/http"

	"github.com/gorilla/mux"
)

type DelHandler struct {
	dhusc domain.UseCase
}

func SetHandlers(router *mux.Router, uc domain.UseCase) {
	handler := &DelHandler{
		dhusc: uc,
	}

	router.HandleFunc("/forum/create", handler.CreateForum).Methods(http.MethodPost)
	router.HandleFunc("/forum/{slug}/create", handler.CreateThreadsForum).Methods(http.MethodPost)
	router.HandleFunc("/post/{id}/details", handler.UpdatePostInfo).Methods(http.MethodPost)
	router.HandleFunc("/service/clear", handler.GetClear).Methods(http.MethodPost)
	router.HandleFunc("/thread/{slug_or_id}/create", handler.CreatePosts).Methods(http.MethodPost)
	router.HandleFunc("/thread/{slug_or_id}/details", handler.UpdateThreadInfo).Methods(http.MethodPost)
	router.HandleFunc("/thread/{slug_or_id}/vote", handler.Voted).Methods(http.MethodPost)
	router.HandleFunc("/user/{nickname}/create", handler.CreateUsers).Methods(http.MethodPost)
	router.HandleFunc("/user/{nickname}/profile", handler.ChangeInfoUser).Methods(http.MethodPost)

	router.HandleFunc("/forum/{slug}/details", handler.ForumInfo).Methods(http.MethodGet)
	router.HandleFunc("/forum/{slug}/users", handler.GetUsersForum).Methods(http.MethodGet)
	router.HandleFunc("/forum/{slug}/threads", handler.GetThreadsForum).Methods(http.MethodGet)
	router.HandleFunc("/post/{id}/details", handler.GetPostInfo).Methods(http.MethodGet)
	router.HandleFunc("/service/status", handler.GetStatus).Methods(http.MethodGet)
	router.HandleFunc("/thread/{slug_or_id}/details", handler.GetThreadInfo).Methods(http.MethodGet)
	router.HandleFunc("/thread/{slug_or_id}/posts", handler.GetPostOfThread).Methods(http.MethodGet)
	router.HandleFunc("/user/{nickname}/profile", handler.GetUser).Methods(http.MethodGet)
}
