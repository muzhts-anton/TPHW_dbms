package rep

import (
	"context"
	"dbms/internal/pkg/database"
	"dbms/internal/pkg/domain"
	_ "dbms/internal/pkg/utils/cast"

	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
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
	resp := r.dbm.Pool.QueryRow(context.Background(), SelectUserByNickname, name)
	// if err != nil {
	// 	return domain.User{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServerError,
	// 	}
	// }

	// if len(resp) == 0 {
	// 	return domain.User{}, domain.NetError{
	// 		Err:        errors.New(domain.ErrorNotFound),
	// 		Statuscode: http.StatusNotFound,
	// 		Message:    domain.ErrorNotFound,
	// 	}
	// }

	var tmp domain.User
	err := resp.Scan(&tmp.Nickname, &tmp.Fullname, &tmp.About, &tmp.Email)
	if err != nil {
		return domain.User{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	return tmp, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) ForumCheck(forum domain.Forum) (domain.Forum, domain.NetError) {
	resp := r.dbm.Pool.QueryRow(context.Background(), SelectSlugFromForum, forum.Slug)
	// if err != nil {
	// 	return domain.Forum{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServerError,
	// 	}
	// }

	// if len(resp) == 0 {
	// 	return domain.Forum{}, domain.NetError{
	// 		Err:        errors.New(domain.ErrorNotFound),
	// 		Statuscode: http.StatusNotFound,
	// 		Message:    domain.ErrorNotFound,
	// 	}
	// }

	// forum.Slug = cast.ToString(resp[0][0])

	err := resp.Scan(&forum.Slug)
	if err != nil {
		return domain.Forum{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	return forum, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) CheckSlug(thread domain.Thread) (domain.Thread, domain.NetError) {
	// resp, err := r.dbm.Query(SelectThreadShort, thread.Slug)
	// if err != nil {
	// 	return domain.Thread{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServerError,
	// 	}
	// }

	// if len(resp) == 0 {
	// 	return domain.Thread{}, domain.NetError{
	// 		Err:        errors.New(domain.ErrorNotFound),
	// 		Statuscode: http.StatusNotFound,
	// 		Message:    domain.ErrorNotFound,
	// 	}
	// }

	// thread.Slug = cast.ToString(resp[0][0])
	// thread.Author = cast.ToString(resp[0][1])

	row := r.dbm.Pool.QueryRow(context.Background(), SelectThreadShort, thread.Slug)
	err := row.Scan(&thread.Slug, &thread.Author)
	if err != nil {
		return domain.Thread{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	return thread, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) GetThreadBySlug(check string, thread domain.Thread) (domain.Thread, domain.NetError) {
	resp := r.dbm.Pool.QueryRow(context.Background(), SelectThread, check)
	// if err != nil {
	// 	return domain.Thread{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServerError,
	// 	}
	// }

	// if len(resp) == 0 {
	// 	return domain.Thread{}, domain.NetError{
	// 		Err:        errors.New(domain.ErrorNotFound),
	// 		Statuscode: http.StatusNotFound,
	// 		Message:    domain.ErrorNotFound,
	// 	}
	// }

	err := resp.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		return domain.Thread{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	return thread, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) InForum(forum domain.Forum) error {
	_, err := r.dbm.Pool.Exec(context.Background(), InsertInForum, forum.Slug, forum.User, forum.Title)
	return err
}

func (r *repHandler) GetForum(slug string) (domain.Forum, domain.NetError) {
	resp := r.dbm.Pool.QueryRow(context.Background(), SelectForumBySlug, slug)
	// if err != nil {
	// 	return domain.Forum{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServerError,
	// 	}
	// }

	// if len(resp) == 0 {
	// 	return domain.Forum{}, domain.NetError{
	// 		Err:        errors.New(domain.ErrorNotFound),
	// 		Statuscode: http.StatusNotFound,
	// 		Message:    domain.ErrorNotFound,
	// 	}
	// }

	var tmp domain.Forum
	err := resp.Scan(&tmp.Title, &tmp.User, &tmp.Slug, &tmp.Posts, &tmp.Threads)
	if err != nil {
		return domain.Forum{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	return tmp, domain.NetError{
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
		if nerr.Err == nil {
			tmp, _ := r.GetThreadBySlug(thread.Slug, trd)
			return tmp, domain.NetError{
				Err:        nerr.Err,
				Statuscode: http.StatusConflict,
				Message:    domain.ErrorConflict,
			}
		}
	}

	row := r.dbm.Pool.QueryRow(context.Background(), InsertThread, thread.Author, thread.Message, thread.Title, thread.Created, thread.Forum, thread.Slug, 0)
	err := row.Scan(&trd.Id)
	if err != nil {
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

	return trd, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusCreated,
		Message:    "",
	}
}

func (r *repHandler) GetThreadSlug(slug string) (domain.Thread, domain.NetError) {
	resp := r.dbm.Pool.QueryRow(context.Background(), SelectThreadSlug, slug)
	// if err != nil {
	// 	return domain.Thread{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServerError,
	// 	}
	// }
	// if len(resp) == 0 {
	// 	return domain.Thread{}, domain.NetError{
	// 		Err:        errors.New(domain.ErrorNotFound),
	// 		Statuscode: http.StatusNotFound,
	// 		Message:    domain.ErrorNotFound,
	// 	}
	// }

	var tmp domain.Thread
	err := resp.Scan(&tmp.Id, &tmp.Title, &tmp.Author, &tmp.Forum, &tmp.Message, &tmp.Votes, &tmp.Slug, &tmp.Created)
	if err != nil {
		return domain.Thread{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	return tmp, domain.NetError{
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

	resp, err := r.dbm.Pool.Query(context.Background(), query, forum.Slug, limit)
	if err != nil {
		return nil, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	// for i := range resp {
	// 	usr = append(usr, domain.User{
	// 		Id:       0,
	// 		Nickname: cast.ToString(resp[i][0]),
	// 		Fullname: cast.ToString(resp[i][1]),
	// 		About:    cast.ToString(resp[i][2]),
	// 		Email:    cast.ToString(resp[i][3]),
	// 	})
	// }

	defer resp.Close() // FIXME: little trick btw
	for resp.Next() {
		user := domain.User{}
		resp.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		usr = append(usr, user)
	}

	return usr, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) GetThreadsOfForum(forum domain.Forum, limit string, since string, desc string) ([]domain.Thread, domain.NetError) {
	trd := make([]domain.Thread, 0)

	if since != "" {
		if desc == "true" {
			resp, err := r.dbm.Pool.Query(context.Background(), GetThreadsSinceDescNotNil, forum.Slug, since, limit)
			if err != nil {
				return nil, domain.NetError{
					Err:        err,
					Statuscode: http.StatusNotFound,
					Message:    domain.ErrorNotFound,
				}
			}

			// for i := range resp {
			// 	trd = append(trd, domain.Thread{
			// 		Id:      cast.ToInt(resp[i][0]),
			// 		Title:   cast.ToString(resp[i][1]),
			// 		Author:  cast.ToString(resp[i][2]),
			// 		Forum:   cast.ToString(resp[i][3]),
			// 		Message: cast.ToString(resp[i][4]),
			// 		Votes:   cast.ToInt(resp[i][5]),
			// 		Slug:    cast.ToString(resp[i][6]),
			// 		Created: cast.ToTime(resp[i][7]),
			// 	})
			// }
			defer resp.Close()
			for resp.Next() {
				threadS := domain.Thread{}
				err := resp.Scan(&threadS.Id, &threadS.Title, &threadS.Author, &threadS.Forum, &threadS.Message, &threadS.Votes, &threadS.Slug, &threadS.Created)
				if err != nil {
					continue
				}

				trd = append(trd, threadS)
			}
		} else {
			resp, err := r.dbm.Pool.Query(context.Background(), GetThreadsSinceDescNil, forum.Slug, since, limit)
			if err != nil {
				return nil, domain.NetError{
					Err:        err,
					Statuscode: http.StatusNotFound,
					Message:    domain.ErrorNotFound,
				}
			}

			// for i := range resp {
			// 	trd = append(trd, domain.Thread{
			// 		Id:      cast.ToInt(resp[i][0]),
			// 		Title:   cast.ToString(resp[i][1]),
			// 		Author:  cast.ToString(resp[i][2]),
			// 		Forum:   cast.ToString(resp[i][3]),
			// 		Message: cast.ToString(resp[i][4]),
			// 		Votes:   cast.ToInt(resp[i][5]),
			// 		Slug:    cast.ToString(resp[i][6]),
			// 		Created: cast.ToTime(resp[i][7]),
			// 	})
			// }
			defer resp.Close()
			for resp.Next() {
				threadS := domain.Thread{}
				err := resp.Scan(&threadS.Id, &threadS.Title, &threadS.Author, &threadS.Forum, &threadS.Message, &threadS.Votes, &threadS.Slug, &threadS.Created)
				if err != nil {
					continue
				}

				trd = append(trd, threadS)
			}
		}
	} else {
		if desc == "true" {
			resp, err := r.dbm.Pool.Query(context.Background(), GetThreadsDescNotNil, forum.Slug, limit)
			if err != nil {
				return nil, domain.NetError{
					Err:        err,
					Statuscode: http.StatusNotFound,
					Message:    domain.ErrorNotFound,
				}
			}

			// for i := range resp {
			// 	trd = append(trd, domain.Thread{
			// 		Id:      cast.ToInt(resp[i][0]),
			// 		Title:   cast.ToString(resp[i][1]),
			// 		Author:  cast.ToString(resp[i][2]),
			// 		Forum:   cast.ToString(resp[i][3]),
			// 		Message: cast.ToString(resp[i][4]),
			// 		Votes:   cast.ToInt(resp[i][5]),
			// 		Slug:    cast.ToString(resp[i][6]),
			// 		Created: cast.ToTime(resp[i][7]),
			// 	})
			// }
			defer resp.Close()
			for resp.Next() {
				threadS := domain.Thread{}
				err := resp.Scan(&threadS.Id, &threadS.Title, &threadS.Author, &threadS.Forum, &threadS.Message, &threadS.Votes, &threadS.Slug, &threadS.Created)
				if err != nil {
					continue
				}

				trd = append(trd, threadS)
			}
		} else {
			resp, err := r.dbm.Pool.Query(context.Background(), GetThreadsDescNil, forum.Slug, limit)
			if err != nil {
				return nil, domain.NetError{
					Err:        err,
					Statuscode: http.StatusNotFound,
					Message:    domain.ErrorNotFound,
				}
			}

			// for i := range resp {
			// 	trd = append(trd, domain.Thread{
			// 		Id:      cast.ToInt(resp[i][0]),
			// 		Title:   cast.ToString(resp[i][1]),
			// 		Author:  cast.ToString(resp[i][2]),
			// 		Forum:   cast.ToString(resp[i][3]),
			// 		Message: cast.ToString(resp[i][4]),
			// 		Votes:   cast.ToInt(resp[i][5]),
			// 		Slug:    cast.ToString(resp[i][6]),
			// 		Created: cast.ToTime(resp[i][7]),
			// 	})
			// }
			defer resp.Close()
			for resp.Next() {
				threadS := domain.Thread{}
				err := resp.Scan(&threadS.Id, &threadS.Title, &threadS.Author, &threadS.Forum, &threadS.Message, &threadS.Votes, &threadS.Slug, &threadS.Created)
				if err != nil {
					continue
				}

				trd = append(trd, threadS)
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
	resp := r.dbm.Pool.QueryRow(context.Background(), SelectThreadId, id)
	// if err != nil {
	// 	return domain.Thread{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServerError,
	// 	}
	// }

	// if len(resp) == 0 {
	// 	return domain.Thread{}, domain.NetError{
	// 		Err:        errors.New(domain.ErrorNotFound),
	// 		Statuscode: http.StatusNotFound,
	// 		Message:    domain.ErrorNotFound,
	// 	}
	// }

	var tmp domain.Thread
	err := resp.Scan(&tmp.Id, &tmp.Title, &tmp.Author, &tmp.Forum, &tmp.Message, &tmp.Votes, &tmp.Slug, &tmp.Created)
	if err != nil {
		return domain.Thread{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	return tmp, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) GetFullPostInfo(posts domain.PostFull, related []string) (domain.PostFull, domain.NetError) {
	resp := r.dbm.Pool.QueryRow(context.Background(), SelectPostById, posts.Post.Id)
	// if err != nil {
	// 	return domain.PostFull{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServerError,
	// 	}
	// }

	// if len(resp) == 0 {
	// 	return domain.PostFull{}, domain.NetError{
	// 		Err:        errors.New(domain.ErrorNotFound),
	// 		Statuscode: http.StatusNotFound,
	// 		Message:    domain.ErrorNotFound,
	// 	}
	// }

	var pstf domain.PostFull
	// pstf.Post = domain.Post{
	// 	Id:       posts.Post.Id,
	// 	Parent:   cast.ToInt(resp[0][0]),
	// 	Author:   cast.ToString(resp[0][1]),
	// 	Message:  cast.ToString(resp[0][2]),
	// 	IsEdited: cast.ToBool(resp[0][3]),
	// 	Forum:    cast.ToString(resp[0][4]),
	// 	Thread:   cast.ToInt(resp[0][5]),
	// 	Created:  cast.ToTime(resp[0][6]),
	// }

	pstf.Post.Id = posts.Post.Id
	err := resp.Scan(&pstf.Post.Parent, &pstf.Post.Author, &pstf.Post.Message, &pstf.Post.IsEdited, &pstf.Post.Forum, &pstf.Post.Thread, &pstf.Post.Created)
	if err != nil {
		return domain.PostFull{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
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
	resp := r.dbm.Pool.QueryRow(context.Background(), UpdatePostMessage, postUpdate.Message, post.Id)
	// if err != nil {
	// 	return domain.Post{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServerError,
	// 	}
	// }

	// if len(resp) == 0 {
	// 	return domain.Post{}, domain.NetError{
	// 		Err:        errors.New(domain.ErrorNotFound),
	// 		Statuscode: http.StatusNotFound,
	// 		Message:    domain.ErrorNotFound,
	// 	}
	// }

	// post = domain.Post{
	// 	Id:       cast.ToInt(resp[0][0]),
	// 	Parent:   cast.ToInt(resp[0][1]),
	// 	Author:   cast.ToString(resp[0][2]),
	// 	Message:  cast.ToString(resp[0][3]),
	// 	IsEdited: cast.ToBool(resp[0][4]),
	// 	Forum:    cast.ToString(resp[0][5]),
	// 	Thread:   cast.ToInt(resp[0][6]),
	// 	Created:  cast.ToTime(resp[0][7]),
	// 	Path:     cast.ToInt8Arr(resp[0][8]),
	// }

	err := resp.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created, &post.Path)
	if err != nil {
		return domain.Post{}, domain.NetError{
			Err:        errors.New(domain.ErrorNotFound),
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	return post, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) GetClear() domain.NetError {
	_, err := r.dbm.Pool.Exec(context.Background(), ClearAll)
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
	resp := r.dbm.Pool.QueryRow(context.Background(), SelectCountUsers)
	err := resp.Scan(&sts.User)
	if err != nil {
		sts.User = 0
	}

	resp = r.dbm.Pool.QueryRow(context.Background(), SelectCountForum)
	err = resp.Scan(&sts.Forum)
	if err != nil {
		sts.Forum = 0
	}

	resp = r.dbm.Pool.QueryRow(context.Background(), SelectCountThreads)
	err = resp.Scan(&sts.Thread)
	if err != nil {
		sts.Thread = 0
	}

	resp = r.dbm.Pool.QueryRow(context.Background(), SelectCountPosts)
	err = resp.Scan(&sts.Post)
	if err != nil {
		sts.Post = 0
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

	// resp, err := r.dbm.Query(query, values...)
	// if err != nil {
	// 	return nil, err
	// }

	// if len(resp) == 0 {
	// 	return nil, errors.New(domain.ErrorInternalServerError)
	// }

	// for i, pst := range posts {
	// 	pst.Id = cast.ToInt(resp[i][0])
	// 	pst.Created = cast.ToTime(resp[i][1])
	// 	pst.Forum = cast.ToString(resp[i][2])
	// 	pst.IsEdited = cast.ToBool(resp[i][3])
	// 	pst.Thread = cast.ToInt(resp[i][4])
	// }

	rows, err := r.dbm.Pool.Query(context.Background(), query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for i := range posts {
		if rows.Next() {
			err := rows.Scan(&posts[i].Id, &posts[i].Created, &posts[i].Forum, &posts[i].IsEdited, &posts[i].Thread)
			if err != nil {
				return nil, err
			}
		}
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return posts, nil
}

func (r *repHandler) UpdateThreadInfo(upThread domain.Thread) (domain.Thread, domain.NetError) {
	var resp pgx.Row
	if upThread.Slug == "" {
		resp = r.dbm.Pool.QueryRow(context.Background(), fmt.Sprintf(UpdateThread, `id=$3`), upThread.Title, upThread.Message, upThread.Id)
	} else {
		resp = r.dbm.Pool.QueryRow(context.Background(), fmt.Sprintf(UpdateThread, `slug=$3`), upThread.Title, upThread.Message, upThread.Slug)
	}

	var tmp domain.Thread
	err := resp.Scan(&tmp.Id, &tmp.Title, &tmp.Author, &tmp.Forum, &tmp.Message, &tmp.Votes, &tmp.Slug, &tmp.Created)
	if err != nil {
		return domain.Thread{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	// if len(resp) == 0 {
	// 	return domain.Thread{}, domain.NetError{
	// 		Err:        errors.New(domain.ErrorNotFound),
	// 		Statuscode: http.StatusNotFound,
	// 		Message:    domain.ErrorNotFound,
	// 	}
	// }

	return tmp, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) GetPostsFlat(limit string, since string, desc string, id int) ([]domain.Post, domain.NetError) {
	pst := make([]domain.Post, 0)
	if since == "" {
		if desc == "true" {
			resp, err := r.dbm.Pool.Query(context.Background(), SelectPostSinceDescNotNil, id, limit)
			if err != nil {
				return nil, domain.NetError{
					Err:        err,
					Statuscode: http.StatusInternalServerError,
					Message:    domain.ErrorInternalServerError,
				}
			}

			// for i := range resp {
			// 	pst = append(pst, domain.Post{
			// 		Id:       cast.ToInt(resp[i][0]),
			// 		Parent:   cast.ToInt(resp[i][1]),
			// 		Author:   cast.ToString(resp[i][2]),
			// 		Message:  cast.ToString(resp[i][3]),
			// 		IsEdited: cast.ToBool(resp[i][4]),
			// 		Forum:    cast.ToString(resp[i][5]),
			// 		Thread:   cast.ToInt(resp[i][6]),
			// 		Created:  cast.ToTime(resp[i][7]),
			// 	})
			// }
			defer resp.Close()
			for resp.Next() {
				onePost := domain.Post{}
				resp.Scan(&onePost.Id, &onePost.Parent, &onePost.Author, &onePost.Message, &onePost.IsEdited, &onePost.Forum, &onePost.Thread, &onePost.Created)
				pst = append(pst, onePost)
			}
		} else {
			resp, err := r.dbm.Pool.Query(context.Background(), SelectPostSinceDescNil, id, limit)
			if err != nil {
				return nil, domain.NetError{
					Err:        err,
					Statuscode: http.StatusInternalServerError,
					Message:    domain.ErrorInternalServerError,
				}
			}

			// for i := range resp {
			// 	pst = append(pst, domain.Post{
			// 		Id:       cast.ToInt(resp[i][0]),
			// 		Parent:   cast.ToInt(resp[i][1]),
			// 		Author:   cast.ToString(resp[i][2]),
			// 		Message:  cast.ToString(resp[i][3]),
			// 		IsEdited: cast.ToBool(resp[i][4]),
			// 		Forum:    cast.ToString(resp[i][5]),
			// 		Thread:   cast.ToInt(resp[i][6]),
			// 		Created:  cast.ToTime(resp[i][7]),
			// 	})
			// }
			defer resp.Close()
			for resp.Next() {
				onePost := domain.Post{}
				resp.Scan(&onePost.Id, &onePost.Parent, &onePost.Author, &onePost.Message, &onePost.IsEdited, &onePost.Forum, &onePost.Thread, &onePost.Created)
				pst = append(pst, onePost)
			}
		}
	} else {
		if desc == "true" {
			resp, err := r.dbm.Pool.Query(context.Background(), SelectPostDescNotNil, id, since, limit)
			if err != nil {
				return nil, domain.NetError{
					Err:        err,
					Statuscode: http.StatusInternalServerError,
					Message:    domain.ErrorInternalServerError,
				}
			}

			// for i := range resp {
			// 	pst = append(pst, domain.Post{
			// 		Id:       cast.ToInt(resp[i][0]),
			// 		Parent:   cast.ToInt(resp[i][1]),
			// 		Author:   cast.ToString(resp[i][2]),
			// 		Message:  cast.ToString(resp[i][3]),
			// 		IsEdited: cast.ToBool(resp[i][4]),
			// 		Forum:    cast.ToString(resp[i][5]),
			// 		Thread:   cast.ToInt(resp[i][6]),
			// 		Created:  cast.ToTime(resp[i][7]),
			// 	})
			// }
			defer resp.Close()
			for resp.Next() {
				onePost := domain.Post{}
				resp.Scan(&onePost.Id, &onePost.Parent, &onePost.Author, &onePost.Message, &onePost.IsEdited, &onePost.Forum, &onePost.Thread, &onePost.Created)
				pst = append(pst, onePost)
			}
		} else {
			resp, err := r.dbm.Pool.Query(context.Background(), SelectPostDescNil, id, since, limit)
			if err != nil {
				return nil, domain.NetError{
					Err:        err,
					Statuscode: http.StatusInternalServerError,
					Message:    domain.ErrorInternalServerError,
				}
			}

			// for i := range resp {
			// 	pst = append(pst, domain.Post{
			// 		Id:       cast.ToInt(resp[i][0]),
			// 		Parent:   cast.ToInt(resp[i][1]),
			// 		Author:   cast.ToString(resp[i][2]),
			// 		Message:  cast.ToString(resp[i][3]),
			// 		IsEdited: cast.ToBool(resp[i][4]),
			// 		Forum:    cast.ToString(resp[i][5]),
			// 		Thread:   cast.ToInt(resp[i][6]),
			// 		Created:  cast.ToTime(resp[i][7]),
			// 	})
			// }
			defer resp.Close()
			for resp.Next() {
				onePost := domain.Post{}
				resp.Scan(&onePost.Id, &onePost.Parent, &onePost.Author, &onePost.Message, &onePost.IsEdited, &onePost.Forum, &onePost.Thread, &onePost.Created)
				pst = append(pst, onePost)
			}
		}
	}

	return pst, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) getTree(id int, since, limit, desc string) (resp pgx.Rows, err error) {
	queryRow := ""

	if limit == "" && since == "" {
		if desc == "true" {
			queryRow += SelectTreeLimitSinceNil
		} else {
			queryRow += SelectTreeLimitSinceDescNil
		}
		resp, err = r.dbm.Pool.Query(context.Background(), queryRow, id)
	} else {
		if limit != "" && since == "" {
			if desc == "true" {
				queryRow += SelectTreeSinceNil
			} else {
				queryRow += SelectTreeSinceDescNil
			}
			resp, err = r.dbm.Pool.Query(context.Background(), queryRow, id, limit)
		}
		if limit != "" && since != "" {
			if desc == "true" {
				queryRow = SelectTreeNotNil
			} else {
				queryRow = SelectTree
			}
			resp, err = r.dbm.Pool.Query(context.Background(), queryRow, id, since, limit)
		}
		if limit == "" && since != "" {
			if desc == "true" {
				queryRow = SelectTreeSinceNilDesc
			} else {
				queryRow = SelectTreeSinceNilDescNil
			}
			resp, err = r.dbm.Pool.Query(context.Background(), queryRow, id, since)
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
	for resp.Next() {
		var onePost domain.Post
		err = resp.Scan(&onePost.Id, &onePost.Parent, &onePost.Author, &onePost.Message, &onePost.IsEdited, &onePost.Forum, &onePost.Thread, &onePost.Created)
		if err != nil {
			return pst, domain.NetError{
				Err:        err,
				Statuscode: http.StatusInternalServerError,
				Message:    domain.ErrorInternalServerError,
			}
		}

		pst = append(pst, onePost)
	}
	// for i := range resp {
	// 	pst = append(pst, domain.Post{
	// 		Id:       cast.ToInt(resp[i][0]),
	// 		Parent:   cast.ToInt(resp[i][1]),
	// 		Author:   cast.ToString(resp[i][2]),
	// 		Message:  cast.ToString(resp[i][3]),
	// 		IsEdited: cast.ToBool(resp[i][4]),
	// 		Forum:    cast.ToString(resp[i][5]),
	// 		Thread:   cast.ToInt(resp[i][6]),
	// 		Created:  cast.ToTime(resp[i][7]),
	// 	})
	// }

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

	resp, err := r.dbm.Pool.Query(context.Background(), queryRow)
	if err != nil {
		return nil, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	pst := make([]domain.Post, 0)
	// for i := range resp {
	// 	pst = append(pst, domain.Post{
	// 		Id:       cast.ToInt(resp[i][0]),
	// 		Parent:   cast.ToInt(resp[i][1]),
	// 		Author:   cast.ToString(resp[i][2]),
	// 		Message:  cast.ToString(resp[i][3]),
	// 		IsEdited: cast.ToBool(resp[i][4]),
	// 		Forum:    cast.ToString(resp[i][5]),
	// 		Thread:   cast.ToInt(resp[i][6]),
	// 		Created:  cast.ToTime(resp[i][7]),
	// 	})
	// }
	for resp.Next() {
		var post domain.Post
		err = resp.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		pst = append(pst, post)
	}

	return pst, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) InVoted(vote domain.Vote) error {
	_, err := r.dbm.Pool.Exec(context.Background(), InsertVote, vote.Nickname, vote.Voice, vote.Thread)
	return err
}

func (r *repHandler) UpVote(vote domain.Vote) (domain.Vote, error) {
	_, err := r.dbm.Pool.Query(context.Background(), UpdateVote, vote.Voice, vote.Nickname, vote.Thread)
	if err != nil {
		return domain.Vote{}, err
	}

	return vote, nil
}

func (r *repHandler) CheckUserEmailUniq(usersS []domain.User) ([]domain.User, domain.NetError) {
	resp, err := r.dbm.Pool.Query(context.Background(), SelectUserByEmailOrNickname, usersS[0].Nickname, usersS[0].Email)
	defer resp.Close()
	if err != nil {
		return []domain.User{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServerError,
		}
	}

	usr := make([]domain.User, 0)
	// for i := range resp {
	// 	usr = append(usr, domain.User{
	// 		Nickname: cast.ToString(resp[i][0]),
	// 		Fullname: cast.ToString(resp[i][1]),
	// 		About:    cast.ToString(resp[i][2]),
	// 		Email:    cast.ToString(resp[i][3]),
	// 	})
	// }
	for resp.Next() {
		userOne := domain.User{}
		resp.Scan(&userOne.Nickname, &userOne.Fullname, &userOne.About, &userOne.Email)
		usr = append(usr, userOne)
	}

	return usr, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) CreateUsers(user domain.User) (domain.User, domain.NetError) {
	_, err := r.dbm.Pool.Exec(context.Background(), `Insert INTO users(Nickname, FullName, About, Email) VALUES ($1, $2, $3, $4);`, user.Nickname, user.Fullname, user.About, user.Email)
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
	// resp, err := r.dbm.Query(UpdateUser, user.Fullname, user.About, user.Email, user.Nickname)
	// if err != nil {
	// 	fmt.Println("huh?")
	// 	return domain.User{}, err
	// }
	// if len(resp) == 0 {
	// 	return domain.User{}, errors.New(domain.ErrorNotFound)
	// }

	// return domain.User{
	// 		Id:       0,
	// 		Nickname: cast.ToString(resp[0][0]),
	// 		Fullname: cast.ToString(resp[0][1]),
	// 		About:    cast.ToString(resp[0][2]),
	// 		Email:    cast.ToString(resp[0][3]),
	// 	},
	// 	nil
	upUser := domain.User{}
	row := r.dbm.Pool.QueryRow(context.Background(), UpdateUser, user.Fullname, user.About, user.Email, user.Nickname)
	err := row.Scan(&upUser.Nickname, &upUser.Fullname, &upUser.About, &upUser.Email)
	if err != nil {
		return domain.User{}, err
	}
	return upUser, nil
}
