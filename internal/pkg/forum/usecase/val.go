package usc

import (
	"dbms/internal/pkg/domain"

	"errors"
	_ "fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgconn"
)

type UscHandler struct {
	uhrep domain.Repository
}

func InitUsc(uhrep domain.Repository) domain.UseCase {
	return &UscHandler{
		uhrep: uhrep,
	}
}

func (u *UscHandler) Forum(forum domain.Forum) (domain.Forum, domain.NetError) {
	usr, nerr := u.uhrep.GetUser(forum.User)
	if nerr.Err != nil {
		return domain.Forum{}, nerr
	}

	forum.User = usr.Nickname

	err := u.uhrep.InForum(forum)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23503" {
			return domain.Forum{}, domain.NetError{
				Err:        err,
				Statuscode: http.StatusNotFound,
				Message:    domain.ErrorNotFound,
			}
		}
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			tmp, _ := u.uhrep.GetForum(forum.Slug)
			return tmp, domain.NetError{
				Err:        err,
				Statuscode: http.StatusConflict,
				Message:    domain.ErrorConflict,
			}
		}
		return domain.Forum{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	forum.Posts = 0
	forum.Threads = 0

	return forum, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusCreated,
		Message:    "",
	}
}

func (u *UscHandler) GetForum(forum domain.Forum) (domain.Forum, domain.NetError) {
	return u.uhrep.GetForum(forum.Slug)
}

func (u *UscHandler) CreateThreadsForum(thread domain.Thread) (domain.Thread, domain.NetError) {
	return u.uhrep.InThread(thread)
}

func (u *UscHandler) GetUsersOfForum(forum domain.Forum, limit string, since string, desc string) ([]domain.User, domain.NetError) {
	_, nerr := u.uhrep.GetForum(forum.Slug)
	if nerr.Err != nil {
		return nil, nerr
	}

	return u.uhrep.GetUsersOfForum(forum, limit, since, desc)
}

func (u *UscHandler) GetThreadsOfForum(forum domain.Forum, limit string, since string, desc string) ([]domain.Thread, domain.NetError) {
	_, nerr := u.uhrep.GetForum(forum.Slug)
	if nerr.Err != nil {
		return nil, nerr
	}

	return u.uhrep.GetThreadsOfForum(forum, limit, since, desc)
}

func (u *UscHandler) GetFullPostInfo(posts domain.PostFull, related []string) (domain.PostFull, domain.NetError) {
	return u.uhrep.GetFullPostInfo(posts, related)
}

func (u *UscHandler) UpdatePostInfo(postUpdate domain.PostUpdate) (domain.Post, domain.NetError) {
	pst, nerr := u.uhrep.UpdatePostInfo(domain.Post{Id: postUpdate.Id}, postUpdate)
	if nerr.Err != nil {
		return domain.Post{}, domain.NetError{
			Err:        nerr.Err,
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	return pst, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (u *UscHandler) GetClear() domain.NetError {
	return u.uhrep.GetClear()
}

func (u *UscHandler) GetStatus() domain.Status {
	return u.uhrep.GetStatus()
}

func (u *UscHandler) CheckThreadIdOrSlug(slugOrId string) (domain.Thread, domain.NetError) {
	id, err := strconv.ParseInt(slugOrId, 10, 0)
	if err != nil {
		return u.uhrep.GetThreadSlug(slugOrId)
	}
	return u.uhrep.GetIdThread(int(id))
}

func (u *UscHandler) CreatePosts(posts []domain.Post, thread domain.Thread) ([]domain.Post, domain.NetError) {
	pst, err := u.uhrep.InPosts(posts, thread)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23503" {
			return nil, domain.NetError{
				Err:        err,
				Statuscode: http.StatusNotFound,
				Message:    domain.ErrorNotFound,
			}
		}
		return nil, domain.NetError{
			Err:        err,
			Statuscode: http.StatusConflict,
			Message:    domain.ErrorConflict,
		}
	}

	return pst, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusCreated,
		Message:    "",
	}
}

func (u *UscHandler) UpdateThreadInfo(slugOrId string, thread domain.Thread) (domain.Thread, domain.NetError) {
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		thread.Slug = slugOrId
	} else {
		thread.Id = id
	}

	return u.uhrep.UpdateThreadInfo(thread)
}

func (u *UscHandler) GetPostOfThread(limit string, since string, desc string, sort string, id int) ([]domain.Post, domain.NetError) {
	switch sort {
	case "flat":
		return u.uhrep.GetPostsFlat(limit, since, desc, id)
	case "tree":
		return u.uhrep.GetPostsTree(limit, since, desc, id)
	case "parent_tree":
		return u.uhrep.GetPostsParent(limit, since, desc, id)
	default:
		return u.uhrep.GetPostsFlat(limit, since, desc, id)
	}
}

func (u *UscHandler) Voted(vote domain.Vote, thread domain.Thread) (domain.Thread, domain.NetError) {
	err := u.uhrep.InVoted(vote)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			_, err := u.uhrep.UpVote(vote)
			if err != nil {
				return domain.Thread{}, domain.NetError{
					Err:        err,
					Statuscode: http.StatusInternalServerError,
					Message:    domain.ErrorInternalServerError,
				}
			}

			return thread, domain.NetError{
				Err:        nil,
				Statuscode: http.StatusOK,
				Message:    "",
			}
		}

		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23503" {
			return domain.Thread{}, domain.NetError{
				Err:        err,
				Statuscode: http.StatusNotFound,
				Message:    domain.ErrorNotFound,
			}
		}

		return domain.Thread{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	return thread, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (u *UscHandler) CreateUsers(user domain.User) ([]domain.User, domain.NetError) {
	usr := make([]domain.User, 0)
	usr = append(usr, user)

	usr, _ = u.uhrep.CheckUserEmailUniq(usr)
	if len(usr) > 0 {
		return usr, domain.NetError{
			Err:        errors.New(domain.ErrorConflict),
			Statuscode: http.StatusConflict,
			Message:    domain.ErrorConflict,
		}
	}

	usr = make([]domain.User, 0)
	newusr, _ := u.uhrep.CreateUsers(user)
	usr = append(usr, newusr)

	return usr, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusCreated,
		Message:    "",
	}
}

func (u *UscHandler) GetUser(user domain.User) (domain.User, domain.NetError) {
	return u.uhrep.GetUser(user.Nickname)
}

func (u *UscHandler) ChangeInfoUser(user domain.User) (domain.User, domain.NetError) {
	usr, err := u.uhrep.ChangeInfoUser(user)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return domain.User{}, domain.NetError{
				Err:        err,
				Statuscode: http.StatusConflict,
				Message:    domain.ErrorConflict,
			}
		}

		return domain.User{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	return usr, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}
