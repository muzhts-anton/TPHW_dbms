package usc

import (
	"dbms/internal/pkg/domain"

	"errors"
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

func (u *UscHandler) UscForum(forum domain.Forum) (domain.Forum, domain.NetError) {
	usr, nerr := u.uhrep.RepGetUser(forum.User)
	if nerr.Err != nil {
		return domain.Forum{}, nerr
	}

	forum.User = usr.Nickname

	err := u.uhrep.RepInForum(forum)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == domain.ErrorPsqlNotFound {
			return domain.Forum{}, domain.NetError{
				Err:        err,
				Statuscode: http.StatusNotFound,
				Message:    domain.ErrorNotFound,
			}
		}
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == domain.ErrorPsqlConflict {
			tmp, _ := u.uhrep.RepGetForum(forum.Slug)
			return tmp, domain.NetError{
				Err:        err,
				Statuscode: http.StatusConflict,
				Message:    domain.ErrorConflict,
			}
		}
		return domain.Forum{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServer,
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

func (u *UscHandler) UscGetForum(forum domain.Forum) (domain.Forum, domain.NetError) {
	return u.uhrep.RepGetForum(forum.Slug)
}

func (u *UscHandler) UscCreateThreadsForum(thread domain.Thread) (domain.Thread, domain.NetError) {
	return u.uhrep.RepInThread(thread)
}

func (u *UscHandler) UscGetUsersOfForum(forum domain.Forum, limit string, since string, desc string) ([]domain.User, domain.NetError) {
	_, nerr := u.uhrep.RepGetForum(forum.Slug)
	if nerr.Err != nil {
		return nil, nerr
	}

	return u.uhrep.RepGetUsersOfForum(forum, limit, since, desc)
}

func (u *UscHandler) UscGetThreadsOfForum(forum domain.Forum, limit string, since string, desc string) ([]domain.Thread, domain.NetError) {
	_, nerr := u.uhrep.RepGetForum(forum.Slug)
	if nerr.Err != nil {
		return nil, nerr
	}

	return u.uhrep.RepGetThreadsOfForum(forum, limit, since, desc)
}

func (u *UscHandler) UscGetFullPostInfo(posts domain.PostFull, related []string) (domain.PostFull, domain.NetError) {
	return u.uhrep.RepGetFullPostInfo(posts, related)
}

func (u *UscHandler) UscUpdatePostInfo(postUpdate domain.PostUpdate) (domain.Post, domain.NetError) {
	pst, nerr := u.uhrep.RepUpdatePostInfo(domain.Post{Id: postUpdate.Id}, postUpdate)
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

func (u *UscHandler) UscGetClear() domain.NetError {
	return u.uhrep.RepGetClear()
}

func (u *UscHandler) UscGetStatus() domain.Status {
	return u.uhrep.RepGetStatus()
}

func (u *UscHandler) UscCheckThreadIdOrSlug(slugOrId string) (domain.Thread, domain.NetError) {
	id, err := strconv.ParseInt(slugOrId, 10, 0)
	if err != nil {
		return u.uhrep.RepGetThreadSlug(slugOrId)
	}
	return u.uhrep.RepGetIdThread(int(id))
}

func (u *UscHandler) UscCreatePosts(posts []domain.Post, thread domain.Thread) ([]domain.Post, domain.NetError) {
	pst, err := u.uhrep.RepInPosts(posts, thread)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == domain.ErrorPsqlNotFound {
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

func (u *UscHandler) UscUpdateThreadInfo(slugOrId string, thread domain.Thread) (domain.Thread, domain.NetError) {
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		thread.Slug = slugOrId
	} else {
		thread.Id = id
	}

	return u.uhrep.RepUpdateThreadInfo(thread)
}

func (u *UscHandler) UscGetPostOfThread(limit string, since string, desc string, sort string, id int) ([]domain.Post, domain.NetError) {
	switch sort {
	case "flat":
		return u.uhrep.RepGetPostsFlat(limit, since, desc, id)
	case "tree":
		return u.uhrep.RepGetPostsTree(limit, since, desc, id)
	case "parent_tree":
		return u.uhrep.RepGetPostsParent(limit, since, desc, id)
	default:
		return u.uhrep.RepGetPostsFlat(limit, since, desc, id)
	}
}

func (u *UscHandler) UscVoted(vote domain.Vote, thread domain.Thread) (domain.Thread, domain.NetError) {
	err := u.uhrep.RepInVoted(vote)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == domain.ErrorPsqlConflict {
			_, err := u.uhrep.RepUpVote(vote)
			if err != nil {
				return domain.Thread{}, domain.NetError{
					Err:        err,
					Statuscode: http.StatusInternalServerError,
					Message:    domain.ErrorInternalServer,
				}
			}

			return thread, domain.NetError{
				Err:        nil,
				Statuscode: http.StatusOK,
				Message:    "",
			}
		}

		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == domain.ErrorPsqlNotFound {
			return domain.Thread{}, domain.NetError{
				Err:        err,
				Statuscode: http.StatusNotFound,
				Message:    domain.ErrorNotFound,
			}
		}

		return domain.Thread{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServer,
		}
	}

	return thread, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (u *UscHandler) UscCreateUsers(user domain.User) ([]domain.User, domain.NetError) {
	usr := make([]domain.User, 0)
	usr = append(usr, user)

	usr, _ = u.uhrep.RepCheckUserEmailUniq(usr)
	if len(usr) > 0 {
		return usr, domain.NetError{
			Err:        errors.New(domain.ErrorConflict),
			Statuscode: http.StatusConflict,
			Message:    domain.ErrorConflict,
		}
	}

	usr = make([]domain.User, 0)
	newusr, _ := u.uhrep.RepCreateUsers(user)
	usr = append(usr, newusr)

	return usr, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusCreated,
		Message:    "",
	}
}

func (u *UscHandler) UscGetUser(user domain.User) (domain.User, domain.NetError) {
	return u.uhrep.RepGetUser(user.Nickname)
}

func (u *UscHandler) UscChangeInfoUser(user domain.User) (domain.User, domain.NetError) {
	usr, err := u.uhrep.RepChangeInfoUser(user)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == domain.ErrorPsqlConflict {
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
