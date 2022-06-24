package rep

import (
	"context"
	"dbms/internal/pkg/database"
	"dbms/internal/pkg/domain"
	_ "dbms/internal/pkg/utils/cast"

	"errors"
	"fmt"
	"net/http"
	_ "strings"
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
	resp := r.dbm.Pool.QueryRow(context.Background(), queryGetUserByNickname, name)
	// if err != nil {
	// 	return domain.User{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServer,
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
	resp := r.dbm.Pool.QueryRow(context.Background(), querySelectSlugFromForum, forum.Slug)
	// if err != nil {
	// 	return domain.Forum{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServer,
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
	// resp, err := r.dbm.Query(querySelectThreadShort, thread.Slug)
	// if err != nil {
	// 	return domain.Thread{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServer,
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

	row := r.dbm.Pool.QueryRow(context.Background(), querySelectThreadShort, thread.Slug)
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
	resp := r.dbm.Pool.QueryRow(context.Background(), queryGetThread, check)
	// if err != nil {
	// 	return domain.Thread{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServer,
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
	_, err := r.dbm.Pool.Exec(context.Background(), queryInsertInForum, forum.Slug, forum.User, forum.Title)
	return err
}

func (r *repHandler) GetForum(slug string) (domain.Forum, domain.NetError) {
	resp := r.dbm.Pool.QueryRow(context.Background(), querySelectForumBySlug, slug)
	// if err != nil {
	// 	return domain.Forum{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServer,
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

	row := r.dbm.Pool.QueryRow(context.Background(), queryInsertThread, thread.Author, thread.Message, thread.Title, thread.Created, thread.Forum, thread.Slug, 0)
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
	resp := r.dbm.Pool.QueryRow(context.Background(), querySelectThreadSlug, slug)
	// if err != nil {
	// 	return domain.Thread{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServer,
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
	var query string
	if desc == "true" {
		if since != "" {
			query = fmt.Sprintf(queryGetUsersOfForumDescNotNilSince, since)
		} else {
			query = queryGetUsersOfForumDescSinceNil
		}
	} else {
		query = fmt.Sprintf(queryGetUsersOfForumDescNil, since)
	}

	usr := make([]domain.User, 0)

	resp, err := r.dbm.Pool.Query(context.Background(), query, forum.Slug, limit)
	if err != nil {
		return nil, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServer,
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
		var tmp domain.User
		resp.Scan(&tmp.Nickname, &tmp.Fullname, &tmp.About, &tmp.Email)
		usr = append(usr, tmp)
	}

	return usr, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) GetThreadsOfForum(forum domain.Forum, limit string, since string, desc string) ([]domain.Thread, domain.NetError) {
	trd := make([]domain.Thread, 0)

	var resp pgx.Rows
	var err error
	if since != "" && desc == "true" {
		resp, err = r.dbm.Pool.Query(context.Background(), queryGetThreadsSinceDescNotNil, forum.Slug, since, limit)
	} else if since != "" && desc != "true" {
		resp, err = r.dbm.Pool.Query(context.Background(), queryGetThreadsSinceDescNil, forum.Slug, since, limit)
	} else if since == "" && desc == "true" {
		resp, err = r.dbm.Pool.Query(context.Background(), queryGetThreadsDescNotNil, forum.Slug, limit)
	} else if since == "" && desc != "true" {
		resp, err = r.dbm.Pool.Query(context.Background(), queryGetThreadsDescNil, forum.Slug, limit)
	}

	if err != nil {
		return nil, domain.NetError{
			Err:        err,
			Statuscode: http.StatusNotFound,
			Message:    domain.ErrorNotFound,
		}
	}

	defer resp.Close()
	for resp.Next() {
		var tmp domain.Thread
		resp.Scan(&tmp.Id, &tmp.Title, &tmp.Author, &tmp.Forum, &tmp.Message, &tmp.Votes, &tmp.Slug, &tmp.Created)
		trd = append(trd, tmp)
	}

	return trd, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) GetIdThread(id int) (domain.Thread, domain.NetError) {
	resp := r.dbm.Pool.QueryRow(context.Background(), querySelectThreadId, id)
	// if err != nil {
	// 	return domain.Thread{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServer,
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
	resp := r.dbm.Pool.QueryRow(context.Background(), querySelectPostById, posts.Post.Id)
	// if err != nil {
	// 	return domain.PostFull{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServer,
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

	for _, item := range related {
		if "thread" == item {
			tmp, _ := r.GetIdThread(pstf.Post.Thread)
			pstf.Thread = &tmp
		}
		if "user" == item {
			tmp, _ := r.GetUser(pstf.Post.Author)
			pstf.Author = &tmp
		}
		if "forum" == item {
			tmp, _ := r.GetForum(pstf.Post.Forum)
			pstf.Forum = &tmp
		}
	}

	return pstf, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) UpdatePostInfo(post domain.Post, postUpdate domain.PostUpdate) (domain.Post, domain.NetError) {
	resp := r.dbm.Pool.QueryRow(context.Background(), queryUpdatePostMessage, postUpdate.Message, post.Id)
	// if err != nil {
	// 	return domain.Post{}, domain.NetError{
	// 		Err:        err,
	// 		Statuscode: http.StatusInternalServerError,
	// 		Message:    domain.ErrorInternalServer,
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
	_, err := r.dbm.Pool.Exec(context.Background(), queryClearAll)
	if err != nil {
		return domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServer,
		}
	}

	return domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) GetStatus() (sts domain.Status) {
	resp := r.dbm.Pool.QueryRow(context.Background(), querySelectCountUsers)
	if err := resp.Scan(&sts.User); err != nil {
		sts.User = 0
	}

	resp = r.dbm.Pool.QueryRow(context.Background(), querySelectCountForum)
	if err := resp.Scan(&sts.Forum); err != nil {
		sts.Forum = 0
	}

	resp = r.dbm.Pool.QueryRow(context.Background(), querySelectCountThreads)
	if err := resp.Scan(&sts.Thread); err != nil {
		sts.Thread = 0
	}

	resp = r.dbm.Pool.QueryRow(context.Background(), querySelectCountPosts)
	if err := resp.Scan(&sts.Post); err != nil {
		sts.Post = 0
	}

	return
}

func (r *repHandler) InPosts(posts []domain.Post, thread domain.Thread) ([]domain.Post, error) {
	now := time.Now()
	var qr string
	var values []interface{}
	for i, item := range posts {
		qr += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
		values = append(values,
			item.Author,
			now,
			thread.Forum,
			item.Message,
			item.Parent,
			thread.Id,
		)
	}

	resp, err := r.dbm.Pool.Query(context.Background(), fmt.Sprintf(queryInsertIntoPosts, qr[:len(qr)-1]), values...)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	for i := range posts {
		if resp.Next() {
			if err := resp.Scan(&posts[i].Id, &posts[i].Created, &posts[i].Forum, &posts[i].IsEdited, &posts[i].Thread); err != nil {
				return nil, err
			}
		}
	}

	if resp.Err() != nil {
		return nil, resp.Err()
	}

	return posts, nil
}

func (r *repHandler) UpdateThreadInfo(upThread domain.Thread) (domain.Thread, domain.NetError) {
	var resp pgx.Row
	if upThread.Slug == "" {
		resp = r.dbm.Pool.QueryRow(context.Background(), fmt.Sprintf(queryUpdateThread, `id`), upThread.Title, upThread.Message, upThread.Id)
	} else {
		resp = r.dbm.Pool.QueryRow(context.Background(), fmt.Sprintf(queryUpdateThread, `slug`), upThread.Title, upThread.Message, upThread.Slug)
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
	var resp pgx.Rows
	var err error
	if since == "" && desc == "true" {
		resp, err = r.dbm.Pool.Query(context.Background(), querySelectPostSinceDescNotNil, id, limit)
	} else if since == "" && desc != "true" {
		resp, err = r.dbm.Pool.Query(context.Background(), querySelectPostSinceDescNil, id, limit)
	} else if since != "" && desc == "true" {
		resp, err = r.dbm.Pool.Query(context.Background(), querySelectPostDescNotNil, id, since, limit)
	} else if since != "" && desc != "true" {
		resp, err = r.dbm.Pool.Query(context.Background(), querySelectPostDescNil, id, since, limit)
	}

	if err != nil {
		return nil, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServer,
		}
	}

	defer resp.Close()
	for resp.Next() {
		var tmp domain.Post
		resp.Scan(&tmp.Id, &tmp.Parent, &tmp.Author, &tmp.Message, &tmp.IsEdited, &tmp.Forum, &tmp.Thread, &tmp.Created)
		pst = append(pst, tmp)
	}

	return pst, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) getTree(id int, since, limit, desc string) (pgx.Rows, error) {
	var qr string
	var params []interface{}

	if limit == "" && since == "" && desc == "true" {
		qr = querySelectTreeLimitSinceNil
		params = []interface{}{id}
	} else if limit == "" && since == "" && desc != "true" {
		qr = querySelectTreeLimitSinceDescNil
		params = []interface{}{id}
	} else if limit == "" && since != "" && desc == "true" {
		qr = querySelectTreeSinceNilDesc
		params = []interface{}{id, since}
	} else if limit == "" && since != "" && desc != "true" {
		qr = querySelectTreeSinceNilDescNil
		params = []interface{}{id, since}
	} else if limit != "" && since == "" && desc == "true" {
		qr = querySelectTreeSinceNil
		params = []interface{}{id, limit}
	} else if limit != "" && since == "" && desc != "true" {
		qr = querySelectTreeSinceDescNil
		params = []interface{}{id, limit}
	} else if limit != "" && since != "" && desc == "true" {
		qr = querySelectTreeNotNil
		params = []interface{}{id, since, limit}
	} else if limit != "" && since != "" && desc != "true" {
		qr = querySelectTree
		params = []interface{}{id, since, limit}
	}

	return r.dbm.Pool.Query(context.Background(), qr, params...)
}

func (r *repHandler) GetPostsTree(limit string, since string, desc string, id int) ([]domain.Post, domain.NetError) {
	resp, err := r.getTree(id, since, limit, desc)
	if err != nil {
		return nil, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServer,
		}
	}
	defer resp.Close()

	pst := make([]domain.Post, 0)
	for resp.Next() {
		var tmp domain.Post
		err = resp.Scan(&tmp.Id, &tmp.Parent, &tmp.Author, &tmp.Message, &tmp.IsEdited, &tmp.Forum, &tmp.Thread, &tmp.Created)
		if err != nil {
			return pst, domain.NetError{
				Err:        err,
				Statuscode: http.StatusInternalServerError,
				Message:    domain.ErrorInternalServer,
			}
		}

		pst = append(pst, tmp)
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

	var qr string
	if desc == "true" {
		qr = fmt.Sprintf(querySelectOnPostsParentDesc, par)
	} else {
		qr = fmt.Sprintf(querySelectOnPostsParentAsc, par)
	}

	resp, err := r.dbm.Pool.Query(context.Background(), qr)
	if err != nil {
		return nil, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServer,
		}
	}

	defer resp.Close()

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
		var tmp domain.Post
		err = resp.Scan(&tmp.Id, &tmp.Parent, &tmp.Author, &tmp.Message, &tmp.IsEdited, &tmp.Forum, &tmp.Thread, &tmp.Created)
		pst = append(pst, tmp)
	}

	return pst, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) InVoted(vote domain.Vote) error {
	_, err := r.dbm.Pool.Exec(context.Background(), queryInsertVote, vote.Nickname, vote.Voice, vote.Thread)
	return err
}

func (r *repHandler) UpVote(vote domain.Vote) (domain.Vote, error) {
	_, err := r.dbm.Pool.Exec(context.Background(), queryUpdateVote, vote.Voice, vote.Nickname, vote.Thread)
	if err != nil {
		return domain.Vote{}, err
	}

	return vote, nil
}

func (r *repHandler) CheckUserEmailUniq(usersS []domain.User) ([]domain.User, domain.NetError) {
	resp, err := r.dbm.Pool.Query(context.Background(), queryGetUserByNicknameEmail, usersS[0].Nickname, usersS[0].Email)
	if err != nil {
		return []domain.User{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServer,
		}
	}

	defer resp.Close()

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
		var tmp domain.User
		resp.Scan(&tmp.Nickname, &tmp.Fullname, &tmp.About, &tmp.Email)
		usr = append(usr, tmp)
	}

	return usr, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusOK,
		Message:    "",
	}
}

func (r *repHandler) CreateUsers(user domain.User) (domain.User, domain.NetError) {
	_, err := r.dbm.Pool.Exec(context.Background(), queryInsertUser, user.Nickname, user.Fullname, user.About, user.Email)
	if err != nil {
		return domain.User{}, domain.NetError{
			Err:        err,
			Statuscode: http.StatusInternalServerError,
			Message:    domain.ErrorInternalServer,
		}
	}

	return user, domain.NetError{
		Err:        nil,
		Statuscode: http.StatusCreated,
		Message:    "",
	}
}

func (r *repHandler) ChangeInfoUser(user domain.User) (domain.User, error) {
	// resp, err := r.dbm.Query(queryUpdateUser, user.Fullname, user.About, user.Email, user.Nickname)
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
	var tmp domain.User
	row := r.dbm.Pool.QueryRow(context.Background(), queryUpdateUser, user.Fullname, user.About, user.Email, user.Nickname)
	err := row.Scan(&tmp.Nickname, &tmp.Fullname, &tmp.About, &tmp.Email)
	if err != nil {
		return domain.User{}, err
	}
	return tmp, nil
}
