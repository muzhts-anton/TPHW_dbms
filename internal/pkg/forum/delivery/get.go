package del

import (
	"dbms/internal/pkg/domain"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
)

// CreateForum /forum/create
func (h *DelHandler) CreateForum(w http.ResponseWriter, r *http.Request) {
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
		out, _ := easyjson.Marshal(domain.ErrorResp{
			Message: nerr.Message,
		})

		w.WriteHeader(nerr.Statuscode)
		w.Write(out)
		return
	}

	out, _ := easyjson.Marshal(*frm)
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// ForumInfo /forum/{slug}/details
func (h *DelHandler) ForumInfo(w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// GetUsersForum /forum/{slug}/users
func (h *DelHandler) GetUsersForum(w http.ResponseWriter, r *http.Request) {
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

	out, _ := easyjson.Marshal(usr) // TODO: fix _easyjson
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// GetThreadsForum /forum/{slug}/threads
func (h *DelHandler) GetThreadsForum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, found := vars["slug"]
	if !found {
		utils.Response(w, domain.NotFound, nil)
		return
	}
	limit := ""
	since := ""
	desc := ""

	query := r.URL.Query()
	if limits := query["limit"]; len(limits) > 0 {
		limit = limits[0]
	}
	if sinces := query["since"]; len(sinces) > 0 {
		since = sinces[0]
	}
	if descs := query["desc"]; len(descs) > 0 {
		desc = descs[0]
	}
	forumS := domain.Forum{Slug: slug}

	users, status := h.dhusc.GetThreadsOfForum(forumS, limit, since, desc)
	utils.Response(w, status, users)
}

// GetPostInfo /post/{id}/details
func (h *DelHandler) GetPostInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idV, found := vars["id"]
	if !found {
		utils.Response(w, domain.NotFound, nil)
		return
	}

	id, _ := strconv.Atoi(idV)
	query := r.URL.Query()

	var related []string
	if relateds := query["related"]; len(relateds) > 0 {
		related = strings.Split(relateds[0], ",")
	}

	postFull := domain.PostFull{}

	postFull.Post.ID = id
	finalPostF, status := h.dhusc.GetFullPostInfo(postFull, related)
	utils.Response(w, status, finalPostF)
}

// UpdatePostInfo /post/{id}/details
func (h *DelHandler) UpdatePostInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids, found := vars["id"]
	if !found {
		utils.Response(w, domain.NotFound, nil)
		return
	}

	postUpdate := domain.PostUpdate{}
	err := easyjson.UnmarshalFromReader(r.Body, &postUpdate)
	if err != nil {
		utils.Response(w, domain.InternalError, nil)
		return
	}
	id, err := strconv.Atoi(ids)

	if err == nil {
		postUpdate.ID = id
	}

	finalPostU, status := h.dhusc.UpdatePostInfo(postUpdate)
	utils.Response(w, status, finalPostU)
}

// GetClear /service/clear
func (h *DelHandler) GetClear(w http.ResponseWriter, _ *http.Request) {
	status := h.dhusc.GetClear()
	utils.Response(w, status, nil)
}

// GetStatus /service/status
func (h *DelHandler) GetStatus(w http.ResponseWriter, _ *http.Request) {
	statusS := h.dhusc.GetStatus()
	utils.Response(w, domain.Okey, statusS)
}

// CreatePosts /thread/{slug_or_id}/create
func (h *DelHandler) CreatePosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId, found := vars["slug_or_id"]
	if !found {
		utils.Response(w, domain.NotFound, nil)
		return
	}

	var posts []domain.Post
	thread, status := h.dhusc.CheckThreadIdOrSlug(slugOrId)
	if status != domain.Okey {
		utils.Response(w, status, nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&posts)
	if err != nil {
		utils.Response(w, domain.InternalError, nil)
		return
	}

	if len(posts) == 0 {
		utils.Response(w, domain.Created, []domain.Post{})
		return
	}

	createPosts, status := h.dhusc.CreatePosts(posts, thread)
	utils.Response(w, status, createPosts)
}

// GetThreadInfo /thread/{slug_or_id}/details
func (h *DelHandler) GetThreadInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId, found := vars["slug_or_id"]
	if !found {
		utils.Response(w, domain.NotFound, nil)
		return
	}
	finalThread, status := h.dhusc.CheckThreadIdOrSlug(slugOrId)
	utils.Response(w, status, finalThread)
}

// UpdateThreadInfo /thread/{slug_or_id}/details
func (h *DelHandler) UpdateThreadInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId, found := vars["slug_or_id"]
	if !found {
		utils.Response(w, domain.NotFound, nil)
		return
	}
	threadS := domain.Thread{}
	err := easyjson.UnmarshalFromReader(r.Body, &threadS)
	if err != nil {
		utils.Response(w, domain.InternalError, nil)
		return
	}
	finalThread, status := h.dhusc.UpdateThreadInfo(slugOrId, threadS)
	utils.Response(w, status, finalThread)
}

// GetPostOfThread /thread/{slug_or_id}/posts
func (h *DelHandler) GetPostOfThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId, found := vars["slug_or_id"]
	if !found {
		utils.Response(w, domain.NotFound, nil)
		return
	}
	limit := ""
	since := ""
	desc := ""
	sort := ""

	query := r.URL.Query()
	if limits := query["limit"]; len(limits) > 0 {
		limit = limits[0]
	}
	if sinces := query["since"]; len(sinces) > 0 {
		since = sinces[0]
	}
	if descs := query["desc"]; len(descs) > 0 {
		desc = descs[0]
	}
	if sorts := query["sort"]; len(sorts) > 0 {
		sort = sorts[0]
	}

	thread, status := h.dhusc.CheckThreadIdOrSlug(slugOrId)
	if status != domain.Okey {
		utils.Response(w, status, nil) // return not found
		return
	}

	finalPosts, status := h.dhusc.GetPostOfThread(limit, since, desc, sort, thread.ID)
	utils.Response(w, status, finalPosts)
}

// Voted /thread/{slug_or_id}/vote
func (h DelHandler) Voted(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId, found := vars["slug_or_id"]
	if !found {
		utils.Response(w, domain.NotFound, nil)
		return
	}

	thread, status := h.dhusc.CheckThreadIdOrSlug(slugOrId)
	if status != domain.Okey {
		utils.Response(w, status, nil) // return not found
		return
	}

	voteS := domain.Vote{}
	err := easyjson.UnmarshalFromReader(r.Body, &voteS)
	if err != nil {
		utils.Response(w, domain.InternalError, nil)
		return
	}

	if thread.ID != 0 {
		voteS.Thread = thread.ID
	}

	_, statusV := h.dhusc.Voted(voteS, thread)
	if statusV != domain.Okey {
		utils.Response(w, statusV, nil)
		return
	}

	finalThread, statusT := h.dhusc.CheckThreadIdOrSlug(slugOrId)
	utils.Response(w, statusT, finalThread)
}

// CreateUsers /user/{nickname}/create
func (h *DelHandler) CreateUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname, found := vars["nickname"]
	if !found {
		utils.Response(w, domain.NotFound, nil)
		return
	}

	userS := domain.User{}
	err := easyjson.UnmarshalFromReader(r.Body, &userS)
	if err != nil {
		utils.Response(w, domain.InternalError, nil)
		return
	}
	userS.NickName = nickname

	finalUser, status := h.dhusc.CreateUsers(userS)
	if status == domain.Created {
		newU := finalUser[0]
		utils.Response(w, status, newU)
		return
	}
	utils.Response(w, status, finalUser)
}

// GetUser /user/{nickname}/profile
func (h *DelHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname, found := vars["nickname"]
	if !found {
		utils.Response(w, domain.NotFound, nil)
		return
	}

	userS := domain.User{}
	userS.NickName = nickname

	finalUser, status := h.dhusc.GetUser(userS)
	utils.Response(w, status, finalUser)
}

// ChangeInfoUser /user/{nickname}/profile
func (h *DelHandler) ChangeInfoUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname, found := vars["nickname"]
	if !found {
		utils.Response(w, domain.NotFound, nil)
		return
	}

	userS := domain.User{}
	err := easyjson.UnmarshalFromReader(r.Body, &userS)
	if err != nil {
		utils.Response(w, domain.InternalError, nil)
		return
	}
	userS.NickName = nickname

	finalUser, status := h.dhusc.ChangeInfoUser(userS)
	utils.Response(w, status, finalUser)
}
