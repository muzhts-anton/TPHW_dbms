package rep

import (
	"dbms/internal/pkg/database"
	"dbms/internal/pkg/domain"
	"dbms/internal/pkg/utils/cast"
	"errors"

	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgconn"
)

type repHandler struct {
	dbm *database.DBManager
}

func InitRep(dbm *database.DBManager) domain.Repository {
	return &repHandler{
		dbm: dbm,
	}
}

func (r *repHandler) GetUser(name string) (domain.User, domain.NetError) {
	resp, err := r.dbm.Query(SelectUserByNickname, name)
	if err != nil {
		return domain.User{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	if len(resp) == 0 {
		return domain.User{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	return domain.User{
			Nickname: cast.ToString(resp[0][0]),
			Fullname: cast.ToString(resp[0][1]),
			About:    cast.ToString(resp[0][2]),
			Email:    cast.ToString(resp[0][3]),
		},
		domain.NetError{
			Err:        nil,
			Statuscode: http.StatusOK,
			Message:    "",
		}
}

func (r *repHandler) ForumCheck(forum domain.Forum) (domain.Forum, domain.NetError) {
	resp, err := r.dbm.Query(SelectSlugFromForum, forum.Slug)
	if err != nil {
		return domain.Forum{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	if len(resp) == 0 {
		return domain.Forum{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	forum.Slug = cast.ToString(resp[0][0])

	return forum, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) CheckSlug(thread domain.Thread) (domain.Thread, domain.NetError) {
	resp, err := r.dbm.Query(SelectThreadShort, thread.Slug)
	if err != nil {
		return domain.Thread{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	if len(resp) == 0 {
		return domain.Thread{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	thread.Slug = cast.ToString(resp[0][0])
	thread.Author = cast.ToString(resp[0][0])

	return thread, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) GetThreadBySlug(check string, thread domain.Thread) (domain.Thread, domain.NetError) {
	resp, err := r.dbm.Query(SelectThread, check)
	if err != nil {
		return domain.Thread{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	if len(resp) == 0 {
		return domain.Thread{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	return domain.Thread{
			Id:      cast.ToInt(resp[0][0]),
			Title:   cast.ToString(resp[0][1]),
			Author:  cast.ToString(resp[0][2]),
			Forum:   cast.ToString(resp[0][3]),
			Message: cast.ToString(resp[0][4]),
			Votes:   cast.ToInt(resp[0][5]),
			Slug:    cast.ToString(resp[0][6]),
			Created: cast.ToTime(resp[0][7]),
		},
		domain.NetError{
			Err:        nil,
			Statuscode: http.StatusOK,
			Message:    "",
		}
}

func (r *repHandler) InForum(forum domain.Forum) error {
	_, err := r.dbm.Query(InsertInForum, forum.Slug, forum.User, forum.Title)
	return err
}

func (r *repHandler) GetForum(slug string) (domain.Forum, domain.NetError) {
	resp, err := r.dbm.Query(SelectForumBySlug, slug)
	if err != nil {
		return domain.Forum{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	if len(resp) == 0 {
		return domain.Forum{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	return domain.Forum{
			Title:   cast.ToString(resp[0][0]),
			User:    cast.ToString(resp[0][1]),
			Slug:    cast.ToString(resp[0][2]),
			Posts:   cast.ToInt(resp[0][3]),
			Threads: cast.ToInt(resp[0][4]),
		},
		domain.NetError{
			Err:        nil,
			Statuscode: http.StatusOK,
			Message:    "",
		}
}

func (r *repHandler) InThread(thread domain.Thread) (domain.Thread, domain.NetError) {
	usr, nerr := r.GetUser(thread.Author)
	if nerr.Err != nil {
		return domain.Thread{}, nerr
	}

	frm, nerr := r.ForumCheck(domain.Forum{Slug: thread.Forum})
	if nerr.Err != nil {
		return domain.Thread{}, nerr
	}

	thread.Author = usr.Nickname
	thread.Forum = frm.Slug

	trd := thread

	if thread.Slug != "" {
		thread, nerr := r.CheckSlug(thread)
		if nerr.Err != nil {
			tmp, _ := r.GetThreadBySlug(thread.Slug, trd)
			return tmp, domain.NetError{
				Err:        nerr.Err,
				Statuscode: http.StatusConflict,
				Message:    domain.ErrorConflict,
			}
		}
	}

	resp, err := r.dbm.Query(InsertThread, thread.Author, thread.Message, thread.Title, thread.Created, thread.Forum, thread.Slug, 0)
	if err != nil {
		return domain.Thread{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	if len(resp) == 0 { // FIXME: wtf with err
		if pgerr, ok := err.(*pgconn.PgError); ok {
			switch pgerr.Code {
			case "23505":
				return trd, domain.NetError{
					Err:        err,
					Statuscode: http.StatusConflict,
					Message:    domain.ErrorConflict,
				}
			default:
				return domain.Thread{}, domain.NetError{
					Err:        err,
					Statuscode: http.StatusNotFound,
					Message:    domain.ErrorNotFound,
				}
			}
		}
	}

	trd.Id = cast.ToInt(resp[0][0])

	return trd, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusCreated,
		Message:    "",
	}
}

func (r *repHandler) GetThreadSlug(slug string) (domain.Thread, domain.NetError) {
	resp, err := r.dbm.Query(SelectThreadSlug, slug)
	if err != nil {
		return domain.Thread{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	if len(resp) == 0 {
		return domain.Thread{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	return domain.Thread{
			Id:      cast.ToInt(resp[0][0]),
			Title:   cast.ToString(resp[0][1]),
			Author:  cast.ToString(resp[0][2]),
			Forum:   cast.ToString(resp[0][3]),
			Message: cast.ToString(resp[0][4]),
			Votes:   cast.ToInt(resp[0][5]),
			Slug:    cast.ToString(resp[0][6]),
			Created: cast.ToTime(resp[0][7]),
		},
		domain.NetError{
			Err:        nil,
			Statuscode: http.StatusOK,
			Message:    "",
		}
}

func (r *repHandler) GetUsersOfForum(forum domain.Forum, limit string, since string, desc string) ([]domain.User, domain.NetError) {
	var query string // TODO: rewrite
	if desc == "true" {
		if since != "" {
			query = fmt.Sprintf(GetUsersOfForumDescNotNilSince, since)
		} else {
			query = GetUsersOfForumDescSinceNil
		}
	} else {
		query = fmt.Sprintf(GetUsersOfForumDescNil, since)
		if since == "" {
			query = fmt.Sprintf(GetUsersOfForumDescNil, 0) // FIXME: why %s <- 0?
		} else {
			query = fmt.Sprintf(GetUsersOfForumDescNil, since)
		}
	}

	usr := make([]domain.User, 0)

	resp, err := r.dbm.Query(query, forum.Slug, limit)
	if err != nil {
		return nil, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	if len(resp) == 0 {
		return nil, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	for i := range resp {
		usr = append(usr, domain.User{
			Id:       0,
			Nickname: cast.ToString(resp[i][0]),
			Fullname: cast.ToString(resp[i][1]),
			About:    cast.ToString(resp[i][2]),
			Email:    cast.ToString(resp[i][3]),
		})
	}

	return usr, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

// FIXME: no len==0 checking
func (r *repHandler) GetThreadsOfForum(forum domain.Forum, limit string, since string, desc string) ([]domain.Thread, domain.NetError) {
	trd := make([]domain.Thread, 0)

	if since != "" {
		if desc == "true" {
			resp, err := r.dbm.Query(GetThreadsSinceDescNotNil, forum.Slug, since, limit)
			if err != nil {
				return nil, domain.NetError{
					Err:        err,
					Statuscode: http.StatusNotFound,
					Message:    domain.ErrorNotFound,
				}
			}

			for i := range resp {
				trd = append(trd, domain.Thread{
					Id:      cast.ToInt(resp[i][0]),
					Title:   cast.ToString(resp[i][1]),
					Author:  cast.ToString(resp[i][2]),
					Forum:   cast.ToString(resp[i][3]),
					Message: cast.ToString(resp[i][4]),
					Votes:   cast.ToInt(resp[i][5]),
					Slug:    cast.ToString(resp[i][6]),
					Created: cast.ToTime(resp[i][7]),
				})
			}
		} else {
			resp, err := r.dbm.Query(GetThreadsSinceDescNil, forum.Slug, since, limit)
			if err != nil {
				return nil, domain.NetError{
					Err:        err,
					Statuscode: http.StatusNotFound,
					Message:    domain.ErrorNotFound,
				}
			}

			for i := range resp {
				trd = append(trd, domain.Thread{
					Id:      cast.ToInt(resp[i][0]),
					Title:   cast.ToString(resp[i][1]),
					Author:  cast.ToString(resp[i][2]),
					Forum:   cast.ToString(resp[i][3]),
					Message: cast.ToString(resp[i][4]),
					Votes:   cast.ToInt(resp[i][5]),
					Slug:    cast.ToString(resp[i][6]),
					Created: cast.ToTime(resp[i][7]),
				})
			}
		}
	} else {
		if desc == "true" {
			resp, err := r.dbm.Query(GetThreadsDescNotNil, forum.Slug, limit)
			if err != nil {
				return nil, domain.NetError{
					Err:        err,
					Statuscode: http.StatusNotFound,
					Message:    domain.ErrorNotFound,
				}
			}

			for i := range resp {
				trd = append(trd, domain.Thread{
					Id:      cast.ToInt(resp[i][0]),
					Title:   cast.ToString(resp[i][1]),
					Author:  cast.ToString(resp[i][2]),
					Forum:   cast.ToString(resp[i][3]),
					Message: cast.ToString(resp[i][4]),
					Votes:   cast.ToInt(resp[i][5]),
					Slug:    cast.ToString(resp[i][6]),
					Created: cast.ToTime(resp[i][7]),
				})
			}
		} else {
			resp, err := r.dbm.Query(GetThreadsDescNil, forum.Slug, limit)
			if err != nil {
				return nil, domain.NetError{
					Err:        err,
					Statuscode: http.StatusNotFound,
					Message:    domain.ErrorNotFound,
				}
			}

			for i := range resp {
				trd = append(trd, domain.Thread{
					Id:      cast.ToInt(resp[i][0]),
					Title:   cast.ToString(resp[i][1]),
					Author:  cast.ToString(resp[i][2]),
					Forum:   cast.ToString(resp[i][3]),
					Message: cast.ToString(resp[i][4]),
					Votes:   cast.ToInt(resp[i][5]),
					Slug:    cast.ToString(resp[i][6]),
					Created: cast.ToTime(resp[i][7]),
				})
			}
		}
	}

	return trd, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) GetIdThread(id int) (domain.Thread, domain.NetError) {
	resp, err := r.dbm.Query(SelectThreadId, id)
	if err != nil {
		return domain.Thread{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	if len(resp) == 0 {
		return domain.Thread{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	return domain.Thread{
			Id:      cast.ToInt(resp[0][0]),
			Title:   cast.ToString(resp[0][1]),
			Author:  cast.ToString(resp[0][2]),
			Forum:   cast.ToString(resp[0][3]),
			Message: cast.ToString(resp[0][4]),
			Votes:   cast.ToInt(resp[0][5]),
			Slug:    cast.ToString(resp[0][6]),
			Created: cast.ToTime(resp[0][7]),
		},
		domain.NetError{
			Err:        nil,
			Statuscode: http.StatusOK,
			Message:    "",
		}
}

func (r *repHandler) GetFullPostInfo(posts domain.PostFull, related []string) (domain.PostFull, domain.NetError) {
	resp, err := r.dbm.Query(SelectPostById, posts.Post.Id)
	if err != nil {
		return domain.PostFull{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	if len(resp) == 0 {
		return domain.PostFull{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	var pstf domain.PostFull
	pstf.Post = domain.Post{
		Id:       posts.Post.Id,
		Parent:   cast.ToInt(resp[0][0]),
		Author:   cast.ToString(resp[0][1]),
		Message:  cast.ToString(resp[0][2]),
		IsEdited: cast.ToBool(resp[0][3]),
		Forum:    cast.ToString(resp[0][4]),
		Thread:   cast.ToInt(resp[0][5]),
		Created:  cast.ToTime(resp[0][6]),
	}

	for i := 0; i < len(related); i++ {
		if "user" == related[i] {
			tmp, _ := r.GetUser(pstf.Post.Author)
			pstf.Author = &tmp
		}
		if "forum" == related[i] {
			tmp, _ := r.GetForum(pstf.Post.Forum)
			pstf.Forum = &tmp
		}
		if "thread" == related[i] {
			tmp, _ := r.GetIdThread(pstf.Post.Thread)
			pstf.Thread = &tmp

		}
	}

	return pstf, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) UpdatePostInfo(post domain.Post, postUpdate domain.PostUpdate) (domain.Post, domain.NetError) {
	resp, err := r.dbm.Query(UpdatePostMessage, postUpdate.Message, post.Id)
	if err != nil {
		return domain.Post{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	if len(resp) == 0 {
		return domain.Post{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	// TODO: check if it is really necessary to change post (postOne). If not just return Post{...}
	post = domain.Post{
		Id:       cast.ToInt(resp[0][0]),
		Parent:   cast.ToInt(resp[0][1]),
		Author:   cast.ToString(resp[0][2]),
		Message:  cast.ToString(resp[0][3]),
		IsEdited: cast.ToBool(resp[0][4]),
		Forum:    cast.ToString(resp[0][5]),
		Thread:   cast.ToInt(resp[0][6]),
		Created:  cast.ToTime(resp[0][7]),
		Path:     cast.ToInt8Arr(resp[0][8]),
	}

	return post, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) GetClear() domain.NetError {
	_, err := r.dbm.Query(ClearAll)
	if err != nil {
		return domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	return domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) GetStatus() (sts domain.Status) {
	resp, err := r.dbm.Query(SelectCountUsers)
	if err != nil {
		sts.User = 0
	} else {
		sts.User = int64(cast.ToUint64(resp[0][0]))
	}

	resp, err = r.dbm.Query(SelectCountForum)
	if err != nil {
		sts.Forum = 0
	} else {
		sts.Forum = int64(cast.ToUint64(resp[0][0]))
	}

	resp, err = r.dbm.Query(SelectCountThreads)
	if err != nil {
		sts.Thread = 0
	} else {
		sts.Thread = int64(cast.ToUint64(resp[0][0]))
	}

	resp, err = r.dbm.Query(SelectCountPosts)
	if err != nil {
		sts.Post = 0
	} else {
		sts.Post = int64(cast.ToUint64(resp[0][0]))
	}

	return
}

// TODO: test & rewrite if necessary
func (r *repHandler) InPosts(posts []domain.Post, thread domain.Thread) ([]domain.Post, error) {
	query := InsertIntoPosts

	var values []interface{}
	created := time.Now()
	for i, post := range posts {
		value := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
		query += value
		values = append(values, post.Author)
		values = append(values, created)
		values = append(values, thread.Forum)
		values = append(values, post.Message)
		values = append(values, post.Parent)
		values = append(values, thread.Id)
	}

	query = strings.TrimSuffix(query, ",")
	query += ` RETURNING id, created, forum, isEdited, thread;`

	resp, err := r.dbm.Query(query, values...)
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New(domain.ErrorInternalServerError)
	}

	for i, pst := range posts {
		pst.Id = cast.ToInt(resp[i][0])
		pst.Created = cast.ToTime(resp[i][1])
		pst.Forum = cast.ToString(resp[i][2])
		pst.IsEdited = cast.ToBool(resp[i][3])
		pst.Thread = cast.ToInt(resp[i][4])
	}

	return posts, nil
}

func (r *repHandler) UpdateThreadInfo(upThread domain.Thread) (domain.Thread, domain.NetError) {
	var err error
	var resp []database.DBbyterow
	if upThread.Slug == "" {
		resp, err = r.dbm.Query(fmt.Sprintf(UpdateThread, `id=$3`), upThread.Title, upThread.Message, upThread.Id)
	} else {
		resp, err = r.dbm.Query(fmt.Sprintf(UpdateThread, `slug=$3`), upThread.Title, upThread.Message, upThread.Slug)
	}

	if err != nil {
		return domain.Thread{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	return domain.Thread{
			Id:      cast.ToInt(resp[0][0]),
			Title:   cast.ToString(resp[0][1]),
			Author:  cast.ToString(resp[0][2]),
			Forum:   cast.ToString(resp[0][3]),
			Message: cast.ToString(resp[0][4]),
			Votes:   cast.ToInt(resp[0][5]),
			Slug:    cast.ToString(resp[0][6]),
			Created: cast.ToTime(resp[0][7]),
		},
		domain.NetError{
			Err:        nil,
			Statuscode: http.StatusOK,
			Message:    "",
		}
}

// TODO: don't foorget 'bout me lmao
func (r *repHandler) GetPostsFlat(limit string, since string, desc string, id int) ([]domain.Post, domain.NetError) {
	pst := make([]domain.Post, 0)
	if since == "" {
		if desc == "true" {
			resp, err := r.dbm.Query(SelectPostSinceDescNotNil, id, limit)
			if err != nil {
				return nil, domain.NetError{
					Err:        err,
					Statuscode: http.StatusInternalServerError,
					Message:    domain.ErrorInternalServerError,
				}
			}

			for i := range resp {
				pst = append(pst, domain.Post{
					Id:       cast.ToInt(resp[i][0]),
					Parent:   cast.ToInt(resp[i][1]),
					Author:   cast.ToString(resp[i][2]),
					Message:  cast.ToString(resp[i][3]),
					IsEdited: cast.ToBool(resp[i][4]),
					Forum:    cast.ToString(resp[i][5]),
					Thread:   cast.ToInt(resp[i][6]),
					Created:  cast.ToTime(resp[i][7]),
				})
			}
		} else {
			resp, err := r.dbm.Query(SelectPostSinceDescNil, id, limit)
			if err != nil {
				return nil, domain.NetError{
					Err:        err,
					Statuscode: http.StatusInternalServerError,
					Message:    domain.ErrorInternalServerError,
				}
			}

			for i := range resp {
				pst = append(pst, domain.Post{
					Id:       cast.ToInt(resp[i][0]),
					Parent:   cast.ToInt(resp[i][1]),
					Author:   cast.ToString(resp[i][2]),
					Message:  cast.ToString(resp[i][3]),
					IsEdited: cast.ToBool(resp[i][4]),
					Forum:    cast.ToString(resp[i][5]),
					Thread:   cast.ToInt(resp[i][6]),
					Created:  cast.ToTime(resp[i][7]),
				})
			}
		}
	} else {
		if desc == "true" {
			resp, err := r.dbm.Query(SelectPostDescNotNil, id, since, limit)
			if err != nil {
				return nil, domain.NetError{
					Err:        err,
					Statuscode: http.StatusInternalServerError,
					Message:    domain.ErrorInternalServerError,
				}
			}

			for i := range resp {
				pst = append(pst, domain.Post{
					Id:       cast.ToInt(resp[i][0]),
					Parent:   cast.ToInt(resp[i][1]),
					Author:   cast.ToString(resp[i][2]),
					Message:  cast.ToString(resp[i][3]),
					IsEdited: cast.ToBool(resp[i][4]),
					Forum:    cast.ToString(resp[i][5]),
					Thread:   cast.ToInt(resp[i][6]),
					Created:  cast.ToTime(resp[i][7]),
				})
			}
		} else {
			resp, err := r.dbm.Query(SelectPostDescNil, id, since, limit)
			if err != nil {
				return nil, domain.NetError{
					Err:        err,
					Statuscode: http.StatusInternalServerError,
					Message:    domain.ErrorInternalServerError,
				}
			}

			for i := range resp {
				pst = append(pst, domain.Post{
					Id:       cast.ToInt(resp[i][0]),
					Parent:   cast.ToInt(resp[i][1]),
					Author:   cast.ToString(resp[i][2]),
					Message:  cast.ToString(resp[i][3]),
					IsEdited: cast.ToBool(resp[i][4]),
					Forum:    cast.ToString(resp[i][5]),
					Thread:   cast.ToInt(resp[i][6]),
					Created:  cast.ToTime(resp[i][7]),
				})
			}
		}
	}

	return pst, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) getTree(id int, since, limit, desc string) (resp []database.DBbyterow, err error) {
	queryRow := ""

	if limit == "" && since == "" {
		if desc == "true" {
			queryRow += SelectTreeLimitSinceNil
		} else {
			queryRow += SelectTreeLimitSinceDescNil
		}
		resp, err = r.dbm.Query(queryRow, id)
	} else {
		if limit != "" && since == "" {
			if desc == "true" {
				queryRow += SelectTreeSinceNil
			} else {
				queryRow += SelectTreeSinceDescNil
			}
			resp, err = r.dbm.Query(queryRow, id, limit)
		}
		if limit != "" && since != "" {
			if desc == "true" {
				queryRow = SelectTreeNotNil
			} else {
				queryRow = SelectTree
			}
			resp, err = r.dbm.Query(queryRow, id, since, limit)
		}
		if limit == "" && since != "" {
			if desc == "true" {
				queryRow = SelectTreeSinceNilDesc
			} else {
				queryRow = SelectTreeSinceNilDescNil
			}
			resp, err = r.dbm.Query(queryRow, id, since)
		}
	}

	return
}

func (r *repHandler) GetPostsTree(limit string, since string, desc string, id int) ([]domain.Post, domain.NetError) {
	resp, err := r.getTree(id, since, limit, desc)
	if err != nil {
		return nil, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	pst := make([]domain.Post, 0)
	for i := range resp {
		pst = append(pst, domain.Post{
			Id:       cast.ToInt(resp[i][0]),
			Parent:   cast.ToInt(resp[i][1]),
			Author:   cast.ToString(resp[i][2]),
			Message:  cast.ToString(resp[i][3]),
			IsEdited: cast.ToBool(resp[i][4]),
			Forum:    cast.ToString(resp[i][5]),
			Thread:   cast.ToInt(resp[i][6]),
			Created:  cast.ToTime(resp[i][7]),
		})
	}

	return pst, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

// TODO: rewrite queries handling
func (r *repHandler) GetPostsParent(limit string, since string, desc string, id int) ([]domain.Post, domain.NetError) {
	par := fmt.Sprintf(`SELECT id FROM posts WHERE thread = %d AND parent = 0 `, id)
	if since != "" {
		if desc == "true" {
			par += ` AND path[1] < ` + fmt.Sprintf(`(SELECT path[1] FROM posts WHERE id = %s) `, since)
		} else {
			par += ` AND path[1] > ` + fmt.Sprintf(`(SELECT path[1] FROM posts WHERE id = %s) `, since)
		}
	}
	if desc == "true" {
		par += ` ORDER BY id DESC `
	} else {
		par += ` ORDER BY id ASC `
	}
	if limit != "" {
		par += " LIMIT " + limit
	}
	queryRow := fmt.Sprintf(`SELECT id, parent, author, message, isedited, forum, thread, created FROM posts WHERE path[1] = ANY (%s) `, par)
	if desc == "true" {
		queryRow += ` ORDER BY path[1] DESC, path,  id `
	} else {
		queryRow += ` ORDER BY path[1] ASC, path,  id `
	}

	resp, err := r.dbm.Query(queryRow)
	if err != nil {
		return nil, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	pst := make([]domain.Post, 0)
	for i := range resp {
		pst = append(pst, domain.Post{
			Id:       cast.ToInt(resp[i][0]),
			Parent:   cast.ToInt(resp[i][1]),
			Author:   cast.ToString(resp[i][2]),
			Message:  cast.ToString(resp[i][3]),
			IsEdited: cast.ToBool(resp[i][4]),
			Forum:    cast.ToString(resp[i][5]),
			Thread:   cast.ToInt(resp[i][6]),
			Created:  cast.ToTime(resp[i][7]),
		})
	}

	return pst, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) InVoted(vote domain.Vote) error {
	_, err := r.dbm.Query(InsertVote, vote.Nickname, vote.Voice, vote.Thread)
	return err
}

func (r *repHandler) UpVote(vote domain.Vote) (domain.Vote, error) {
	_, err := r.dbm.Query(UpdateVote, vote.Voice, vote.Nickname, vote.Thread)
	if err != nil {
		return domain.Vote{}, err
	}

	return vote, nil
}

func (r *repHandler) CheckUserEmailUniq(usersS []domain.User) ([]domain.User, domain.NetError) {
	resp, err := r.dbm.Query(SelectUserByEmailOrNickname, usersS[0].Nickname, usersS[0].Email)
	if err != nil {
		return []domain.User{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	usr := make([]domain.User, 0)
	for i := range resp {
		usr = append(usr, domain.User{
			Nickname: cast.ToString(resp[i][0]),
			Fullname: cast.ToString(resp[i][1]),
			About:    cast.ToString(resp[i][2]),
			Email:    cast.ToString(resp[i][3]),
		})
	}

	return usr, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) CreateUsers(user domain.User) (domain.User, domain.NetError) {
	_, err := r.dbm.Query(`Insert INTO users(Nickname, FullName, About, Email) VALUES ($1, $2, $3, $4);`, user.Nickname, user.Fullname, user.About, user.Email)
	if err != nil {
		return domain.User{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	return user, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusCreated,
		Message:    "",
	}
}

func (r *repHandler) ChangeInfoUser(user domain.User) (domain.User, error) {
	resp, err := r.dbm.Query(UpdateUser, user.Fullname, user.About, user.Email, user.Nickname)
	if err != nil {
		return domain.User{}, err
	}
	if len(resp) == 0 {
		return domain.User{}, errors.New(domain.ErrorNotFound)
	}

	return domain.User{
			Id:       0,
			Nickname: cast.ToString(resp[0][0]),
			Fullname: cast.ToString(resp[0][1]),
			About:    cast.ToString(resp[0][2]),
			Email:    cast.ToString(resp[0][3]),
		},
		nil
}
