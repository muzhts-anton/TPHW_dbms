package domain

type UseCase interface {
	Forum(forum Forum) (Forum, NetError)
	GetForum(forum Forum) (Forum, NetError)
	CreateThreadsForum(thread Thread) (Thread, NetError)
	GetUsersOfForum(forum Forum, limit string, since string, desc string) ([]User, NetError)
	GetThreadsOfForum(forum Forum, limit string, since string, desc string) ([]Thread, NetError)
	GetFullPostInfo(posts PostFull, related []string) (PostFull, NetError)
	UpdatePostInfo(postUpdate PostUpdate) (Post, NetError)
	GetClear() NetError
	GetStatus() Status
	CheckThreadIdOrSlug(slugOrId string) (Thread, NetError)
	CreatePosts(createPosts []Post, thread Thread) ([]Post, NetError)
	UpdateThreadInfo(slugOrId string, upThread Thread) (Thread, NetError)
	GetPostOfThread(limit string, since string, desc string, sort string, ID int) ([]Post, NetError)
	Voted(vote Vote, thread Thread) (Thread, NetError)
	CreateUsers(user User) ([]User, NetError)
	GetUser(user User) (User, NetError)
	ChangeInfoUser(user User) (User, NetError)
}

type Repository interface {
	GetUser(name string) (User, NetError)
	InForum(forum Forum) error
	GetForum(slug string) (Forum, NetError)
	InThread(thread Thread) (Thread, NetError)
	GetThreadSlug(slug string) (Thread, NetError)
	GetUsersOfForum(forum Forum, limit string, since string, desc string) ([]User, NetError)
	GetThreadsOfForum(forum Forum, limit string, since string, desc string) ([]Thread, NetError)
	GetFullPostInfo(posts PostFull, related []string) (PostFull, NetError)
	GetIdThread(id int) (Thread, NetError)
	UpdatePostInfo(post Post, postUpdate PostUpdate) (Post, NetError)
	GetClear() NetError
	GetStatus() Status
	InPosts(posts []Post, thread Thread) ([]Post, error)
	UpdateThreadInfo(upThread Thread) (Thread, NetError)
	GetPostsFlat(limit string, since string, desc string, ID int) ([]Post, NetError)
	GetPostsTree(limit string, since string, desc string, ID int) ([]Post, NetError)
	GetPostsParent(limit string, since string, desc string, ID int) ([]Post, NetError)
	InVoted(vote Vote) error
	UpVote(vote Vote) (Vote, error)
	CreateUsers(user User) (User, NetError)
	ChangeInfoUser(user User) (User, error)
	CheckUserEmailUniq(user []User) ([]User, NetError)
	ForumCheck(forum Forum) (Forum, NetError)
}
