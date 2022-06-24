package del

import (
	"dbms/internal/pkg/domain"

	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
)

// CreateForum /forum/create
func (h *DelHandler) CreateForum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorInternalServerError,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	frm := new(domain.Forum)
	err = easyjson.Unmarshal(b, frm)
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorInternalServerError,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	var nerr domain.NetError
	*frm, nerr = h.dhusc.Forum(*frm)
	if nerr.Err != nil {
		var out []byte
		if nerr.Statuscode == http.StatusConflict {
			out, _ = easyjson.Marshal(frm)
		} else {
			out, _ = easyjson.Marshal(domain.ErrorResp{
				Message: nerr.Message,
			})
		}

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	out, _ := easyjson.Marshal(*frm)
	w.WriteHeader(nerr.Statuscode)
	w.Write(out)
}

// ForumInfo /forum/{slug}/details
func (h *DelHandler) ForumInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	param, ok := mux.Vars(r)["slug"]
	if !ok {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorNotFound,
		})

		w.WriteHeader(http.StatusNotFound)
		w.Write(out)
		return
	}

	frm, nerr := h.dhusc.GetForum(domain.Forum{Slug: param})
	if nerr.Err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	out, _ := easyjson.Marshal(frm)
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// CreateThreadsForum /forum/{slug}/create
func (h *DelHandler) CreateThreadsForum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorInternalServerError,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	trd := new(domain.Thread)
	err = easyjson.Unmarshal(b, trd)
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorInternalServerError,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	var ok bool
	trd.Forum, ok = mux.Vars(r)["slug"]
	if !ok {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorNotFound,
		})

		w.WriteHeader(http.StatusNotFound)
		w.Write(out)
		return
	}

	var nerr domain.NetError
	*trd, nerr = h.dhusc.CreateThreadsForum(*trd)
	if nerr.Err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	out, _ := easyjson.Marshal(trd)
	w.WriteHeader(nerr.Statuscode)
	w.Write(out)
}

// GetUsersForum /forum/{slug}/users
func (h *DelHandler) GetUsersForum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	param, ok := mux.Vars(r)["slug"]
	if !ok {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorNotFound,
		})

		w.WriteHeader(http.StatusNotFound)
		w.Write(out)
		return
	}

	var limit, since, desc string // TODO: pack to struct{...}
	query := r.URL.Query()

	var tmp []string
	if tmp = query["limit"]; len(tmp) > 0 {
		limit = tmp[0]
	}
	if tmp := query["since"]; len(tmp) > 0 {
		since = tmp[0]
	}
	if tmp := query["desc"]; len(tmp) > 0 {
		desc = tmp[0]
	}
	if limit == "" { // TODO: move to usecase
		limit = "100"
	}

	usr, nerr := h.dhusc.GetUsersOfForum(domain.Forum{Slug: param}, limit, since, desc)
	if nerr.Err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	out, _ := easyjson.Marshal(domain.Users(usr))
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// GetThreadsForum /forum/{slug}/threads
func (h *DelHandler) GetThreadsForum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	param, ok := mux.Vars(r)["slug"]
	if !ok {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorNotFound,
		})

		w.WriteHeader(http.StatusNotFound)
		w.Write(out)
		return
	}

	var limit, since, desc string // TODO: pack to struct{...}
	query := r.URL.Query()

	var tmp []string
	if tmp = query["limit"]; len(tmp) > 0 {
		limit = tmp[0]
	}
	if tmp := query["since"]; len(tmp) > 0 {
		since = tmp[0]
	}
	if tmp := query["desc"]; len(tmp) > 0 {
		desc = tmp[0]
	}

	trd, nerr := h.dhusc.GetThreadsOfForum(domain.Forum{Slug: param}, limit, since, desc)
	if nerr.Err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	out, _ := easyjson.Marshal(domain.Threads(trd))
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// GetPostInfo /post/{id}/details
func (h *DelHandler) GetPostInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 0) // TODO: mb here error of map will be incorrect
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorNotFound,
		})

		w.WriteHeader(http.StatusNotFound)
		w.Write(out)
		return
	}

	var related []string
	if tmp := r.URL.Query()["related"]; len(tmp) > 0 {
		related = strings.Split(tmp[0], ",")
	}

	var pf domain.PostFull
	pf.Post.Id = int(id)

	pf, nerr := h.dhusc.GetFullPostInfo(pf, related)
	if nerr.Err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	out, _ := easyjson.Marshal(pf)
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// UpdatePostInfo /post/{id}/details
func (h *DelHandler) UpdatePostInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 0) // TODO: mb here error of map will be incorrect
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorNotFound,
		})

		w.WriteHeader(http.StatusNotFound)
		w.Write(out)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorInternalServerError,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	pu := new(domain.PostUpdate)
	err = easyjson.Unmarshal(b, pu)
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorInternalServerError,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	pu.Id = int(id)

	pst, nerr := h.dhusc.UpdatePostInfo(*pu)
	if nerr.Err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	out, _ := easyjson.Marshal(pst)
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// GetClear /service/clear
func (h *DelHandler) GetClear(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nerr := h.dhusc.GetClear()
	if nerr.Err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetStatus /service/status
func (h *DelHandler) GetStatus(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sts := h.dhusc.GetStatus()

	out, _ := easyjson.Marshal(sts)
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// CreatePosts /thread/{slug_or_id}/create
func (h *DelHandler) CreatePosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slid, ok := mux.Vars(r)["slug_or_id"]
	if !ok {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorNotFound,
		})

		w.WriteHeader(http.StatusNotFound)
		w.Write(out)
		return
	}

	trd, nerr := h.dhusc.CheckThreadIdOrSlug(slid)
	if nerr.Err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorInternalServerError,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	var pst domain.Posts
	err = easyjson.Unmarshal(b, &pst)
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorInternalServerError,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	if len(pst) == 0 {
		out, _ := easyjson.Marshal(pst)
		w.WriteHeader(http.StatusCreated)
		w.Write(out)
		return
	}

	pst, nerr = h.dhusc.CreatePosts(pst, trd)
	if nerr.Err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	out, _ := easyjson.Marshal(pst)
	w.WriteHeader(nerr.Statuscode)
	w.Write(out)
}

// GetThreadInfo /thread/{slug_or_id}/details
func (h *DelHandler) GetThreadInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slid, ok := mux.Vars(r)["slug_or_id"]
	if !ok {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorNotFound,
		})

		w.WriteHeader(http.StatusNotFound)
		w.Write(out)
		return
	}

	trd, nerr := h.dhusc.CheckThreadIdOrSlug(slid)
	if nerr.Err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	out, _ := easyjson.Marshal(trd)
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// UpdateThreadInfo /thread/{slug_or_id}/details
func (h *DelHandler) UpdateThreadInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slid, ok := mux.Vars(r)["slug_or_id"]
	if !ok {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorNotFound,
		})

		w.WriteHeader(http.StatusNotFound)
		w.Write(out)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorInternalServerError,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	trd := new(domain.Thread)
	err = easyjson.Unmarshal(b, trd)
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorInternalServerError,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	var nerr domain.NetError
	*trd, nerr = h.dhusc.UpdateThreadInfo(slid, *trd)
	if nerr.Err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	out, _ := easyjson.Marshal(trd)
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// GetPostOfThread /thread/{slug_or_id}/posts
func (h *DelHandler) GetPostOfThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slid, ok := mux.Vars(r)["slug_or_id"]
	if !ok {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorNotFound,
		})

		w.WriteHeader(http.StatusNotFound)
		w.Write(out)
		return
	}

	var limit, since, desc, sort string
	query := r.URL.Query()

	var tmp []string
	if tmp = query["limit"]; len(tmp) > 0 {
		limit = tmp[0]
	}
	if tmp = query["since"]; len(tmp) > 0 {
		since = tmp[0]
	}
	if tmp = query["desc"]; len(tmp) > 0 {
		desc = tmp[0]
	}
	if tmp = query["sort"]; len(tmp) > 0 {
		sort = tmp[0]
	}

	trd, nerr := h.dhusc.CheckThreadIdOrSlug(slid)
	if nerr.Err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	pst, nerr := h.dhusc.GetPostOfThread(limit, since, desc, sort, trd.Id)
	if nerr.Err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	out, _ := easyjson.Marshal(domain.Posts(pst))
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// Voted /thread/{slug_or_id}/vote
func (h DelHandler) Voted(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slid, ok := mux.Vars(r)["slug_or_id"]
	if !ok {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorNotFound,
		})

		w.WriteHeader(http.StatusNotFound)
		w.Write(out)
		return
	}

	trd, nerr := h.dhusc.CheckThreadIdOrSlug(slid)
	if nerr.Statuscode != http.StatusOK {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorInternalServerError,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	vt := new(domain.Vote)
	err = easyjson.Unmarshal(b, vt)
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorInternalServerError,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	if trd.Id != 0 {
		vt.Thread = trd.Id
	}

	_, nerr = h.dhusc.Voted(*vt, trd)
	if nerr.Statuscode != http.StatusOK {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	trd, nerr = h.dhusc.CheckThreadIdOrSlug(slid)
	if nerr.Err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	out, _ := easyjson.Marshal(trd)
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// CreateUsers /user/{nickname}/create
func (h *DelHandler) CreateUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorInternalServerError,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	usr := new(domain.User)
	err = easyjson.Unmarshal(b, usr)
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorInternalServerError,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	var ok bool
	usr.Nickname, ok = mux.Vars(r)["nickname"]
	if !ok {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorNotFound,
		})

		w.WriteHeader(http.StatusNotFound)
		w.Write(out)
		return
	}

	usrs, nerr := h.dhusc.CreateUsers(*usr)
	if nerr.Err != nil {
		var out []byte
		if nerr.Statuscode == http.StatusConflict {
			out, _ = easyjson.Marshal(domain.Users(usrs))
		} else {
			out, _ = easyjson.Marshal(domain.ErrorResp{
				Message: nerr.Message,
			})
		}

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	var out []byte
	if nerr.Statuscode == http.StatusCreated {
		out, _ = easyjson.Marshal(usrs[0])
		w.WriteHeader(http.StatusCreated)
	} else {
		out, _ = easyjson.Marshal(domain.Users(usrs))
		w.WriteHeader(http.StatusOK)
	}
	w.Write(out)
}

// GetUser /user/{nickname}/profile
func (h *DelHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ok bool
	var usr domain.User
	usr.Nickname, ok = mux.Vars(r)["nickname"]
	if !ok {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorNotFound,
		})

		w.WriteHeader(http.StatusNotFound)
		w.Write(out)
		return
	}

	usr, nerr := h.dhusc.GetUser(usr)
	if nerr.Err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	out, _ := easyjson.Marshal(usr)
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// ChangeInfoUser /user/{nickname}/profile
func (h *DelHandler) ChangeInfoUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorInternalServerError,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	usr := new(domain.User)
	err = easyjson.Unmarshal(b, usr)
	if err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorInternalServerError,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	var ok bool
	usr.Nickname, ok = mux.Vars(r)["nickname"]
	if !ok {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: domain.ErrorNotFound,
		})

		w.WriteHeader(http.StatusNotFound)
		w.Write(out)
		return
	}

	var nerr domain.NetError
	*usr, nerr = h.dhusc.ChangeInfoUser(*usr)
	if nerr.Err != nil {
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	out, _ := easyjson.Marshal(usr)
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
