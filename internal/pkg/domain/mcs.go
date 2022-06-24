package domain

type UseCase interface {
	UscForum(forum Forum) (Forum, NetError)
	UscGetForum(forum Forum) (Forum, NetError)
	UscCreateThreadsForum(thread Thread) (Thread, NetError)
	UscGetUsersOfForum(forum Forum, limit string, since string, desc string) ([]User, NetError)
	UscGetThreadsOfForum(forum Forum, limit string, since string, desc string) ([]Thread, NetError)
	UscGetFullPostInfo(posts PostFull, related []string) (PostFull, NetError)
	UscUpdatePostInfo(postUpdate PostUpdate) (Post, NetError)
	UscGetClear() NetError
	UscGetStatus() Status
	UscCheckThreadIdOrSlug(slugOrId string) (Thread, NetError)
	UscCreatePosts(createPosts []Post, thread Thread) ([]Post, NetError)
	UscUpdateThreadInfo(slugOrId string, upThread Thread) (Thread, NetError)
	UscGetPostOfThread(limit string, since string, desc string, sort string, ID int) ([]Post, NetError)
	UscVoted(vote Vote, thread Thread) (Thread, NetError)
	UscCreateUsers(user User) ([]User, NetError)
	UscGetUser(user User) (User, NetError)
	UscChangeInfoUser(user User) (User, NetError)
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
