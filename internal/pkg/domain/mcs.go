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
	RepGetUser(name string) (User, NetError)
	RepInForum(forum Forum) error
	RepGetForum(slug string) (Forum, NetError)
	RepInThread(thread Thread) (Thread, NetError)
	RepGetThreadSlug(slug string) (Thread, NetError)
	RepGetUsersOfForum(forum Forum, limit string, since string, desc string) ([]User, NetError)
	RepGetThreadsOfForum(forum Forum, limit string, since string, desc string) ([]Thread, NetError)
	RepGetFullPostInfo(posts PostFull, related []string) (PostFull, NetError)
	RepGetIdThread(id int) (Thread, NetError)
	RepUpdatePostInfo(post Post, postUpdate PostUpdate) (Post, NetError)
	RepGetClear() NetError
	RepGetStatus() Status
	RepInPosts(posts []Post, thread Thread) ([]Post, error)
	RepUpdateThreadInfo(upThread Thread) (Thread, NetError)
	RepGetPostsFlat(limit string, since string, desc string, ID int) ([]Post, NetError)
	RepGetPostsTree(limit string, since string, desc string, ID int) ([]Post, NetError)
	RepGetPostsParent(limit string, since string, desc string, ID int) ([]Post, NetError)
	RepInVoted(vote Vote) error
	RepUpVote(vote Vote) (Vote, error)
	RepCreateUsers(user User) (User, NetError)
	RepChangeInfoUser(user User) (User, error)
	RepCheckUserEmailUniq(user []User) ([]User, NetError)
	ForumCheck(forum Forum) (Forum, NetError)
}
