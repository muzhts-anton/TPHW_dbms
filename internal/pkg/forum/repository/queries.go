package rep

// users
const (
	queryGetUserByNickname = `
	SELECT nickname, fullname, about, email
	FROM users
	WHERE nickname = $1
	LIMIT 1;
	`

	queryGetUserByNicknameEmail = `
	SELECT nickname, fullname, about, email
	FROM users
	WHERE nickname = $1 or email = $2
	LIMIT 2;
	`

	queryGetUsersOfForumDescNotNilSince = `
	SELECT nickname, fullname, about, email
	FROM users_forum
	WHERE slug = $1 and nickname < '%s'
	ORDER BY nickname DESC
	LIMIT nullif($2, 0);
	`

	queryGetUsersOfForumDescSinceNil = `
	SELECT nickname, fullname, about, email
	FROM users_forum
	WHERE slug = $1
	ORDER BY nickname DESC
	LIMIT nullif($2, 0);
	`

	queryGetUsersOfForumDescNil = `
	SELECT nickname, fullname, about, email
	FROM users_forum
	WHERE slug = $1 and nickname > '%s'
	ORDER BY nickname
	LIMIT nullif($2, 0);
	`

	queryUpdateUser = `
	UPDATE users
	SET
		fullname = coalesce(nullif($1, ''), fullname),
		about = coalesce(nullif($2, ''), about),
		email = coalesce(nullif($3, ''), email)
	WHERE nickname = $4
	RETURNING nickname, fullname, about, email;
	`

	queryInsertUser = `
	INSERT INTO users (nickname, fullname, about,email)
	VALUES ($1, $2, $3, $4);
	`
)

// thread
const (
	queryGetThread = `
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE slug = $1
	LIMIT 1;
	`

	querySelectThreadSlug = `
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE slug = $1
	LIMIT 1;
	`
	queryGetThreadsSinceDescNotNil = `
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE forum = $1 and created <= $2
	ORDER BY created DESC
	LIMIT $3;
	`

	queryGetThreadsSinceDescNil = `
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE forum = $1 AND created >= $2
	ORDER BY created ASC
	LIMIT $3;
	`

	queryGetThreadsDescNotNil = `
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE forum = $1
	ORDER BY created DESC
	LIMIT $2;
	`

	queryGetThreadsDescNil = `
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE forum = $1
	ORDER BY created ASC
	LIMIT $2;
	`

	querySelectThreadId = `
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE id = $1
	LIMIT 1;
	`

	queryInsertThread = `
	INSERT INTO threads (author, message, title, created, forum, slug, votes)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id;
	`

	queryUpdateThread = `
	UPDATE threads
	SET
		title = coalesce(nullif($1, ''), title),
		message = coalesce(nullif($2, ''), message)
	WHERE %s = $3
	RETURNING id, title, author, forum, message, votes, slug, created;
	`

	querySelectThreadShort = `
	SELECT slug, author
	FROM threads
	WHERE slug = $1;
	`
)

// votes
const (
	queryUpdateVote = `
	UPDATE votes
	SET voice = $1
	WHERE author = $2 AND thread = $3;
	`

	queryInsertVote = `
	INSERT INTO votes (author, voice, thread)
	VALUES ($1, $2, $3);
	`
)

// forum
const (
	querySelectForumBySlug = `
	SELECT title, "user", slug, posts, threads
	FROM forum
	WHERE slug = $1
	LIMIT 1;
	`

	queryInsertInForum = `
	INSERT INTO forum (slug, "user", title)
	VALUES ($1, $2, $3);
	`

	querySelectSlugFromForum = `
	SELECT slug
	FROM forum
	WHERE slug = $1;
	`
)

// other
const (
	queryClearAll = `
	TRUNCATE table users, forum, threads, posts, votes, users_forum CASCADE;
	`

	querySelectCountUsers = `
	SELECT COUNT(*) FROM users;
	`

	querySelectCountForum = `
	SELECT COUNT(*) FROM forum;
	`

	querySelectCountThreads = `
	SELECT COUNT(*) FROM threads;
	`

	querySelectCountPosts = `
	SELECT COUNT(*) FROM posts;
	`
)

// posts
const (
	querySelectPostById = `
	SELECT parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE id = $1;
	`

	queryUpdatePostMessage = `
	UPDATE posts
	SET
		message = coalesce(nullif($1, ''), message),
		isedited = case WHEN $1 = '' OR message = $1 THEN isedited ELSE true END
	WHERE id = $2
	RETURNING id, parent, author, message, isedited, forum, thread, created, path;
	`

	querySelectPostSinceDescNotNil = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE thread = $1
	ORDER BY id DESC
	LIMIT $2;
	`

	querySelectPostSinceDescNil = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE thread = $1
	ORDER BY id
	LIMIT $2;
	`

	querySelectPostDescNotNil = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE thread = $1 and id < $2
	ORDER BY id DESC
	LIMIT $3;
	`

	querySelectPostDescNil = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE thread = $1 and id > $2
	ORDER BY id
	LIMIT $3;
	`

	queryInsertIntoPosts = `
	INSERT INTO posts (author, created, forum, message, parent, thread)
	VALUES %s
	RETURNING id, created, forum, isedited, thread;
	`

	querySelectTreeLimitSinceNil = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE thread = $1
	ORDER BY path, id DESC;
	`

	querySelectTreeLimitSinceDescNil = `
	SELECT id, parent, author, message, isedited, forum, thread, created FROM posts
	WHERE thread = $1
	ORDER BY path, id ASC;
	`

	querySelectTreeSinceNil = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE thread = $1
	ORDER BY path DESC, id DESC
	LIMIT $2;
	`

	querySelectTreeSinceDescNil = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE thread = $1
	ORDER BY path, id ASC
	LIMIT $2;
	`

	querySelectTreeNotNil = `
	SELECT posts.id, posts.parent, posts.author, posts.message, posts.isedited, posts.forum, posts.thread, posts.created
	FROM posts
	JOIN posts parent ON parent.id = $2
	WHERE posts.path < parent.path AND posts.thread = $1
	ORDER BY posts.path DESC, posts.id DESC
	LIMIT $3;
	`

	querySelectTreeSinceNilDesc = `
	SELECT posts.id, posts.parent, posts.author, posts.message, posts.isedited, posts.forum, posts.thread, posts.created
	FROM posts
	JOIN posts parent ON parent.id = $2
	WHERE posts.path < parent.path AND posts.thread = $1
	ORDER BY posts.path DESC, posts.id DESC;
	`

	querySelectTreeSinceNilDescNil = `
	SELECT posts.id, posts.parent, posts.author, posts.message, posts.isedited, posts.forum, posts.thread, posts.created
	FROM posts
	JOIN posts parent ON parent.id = $2
	WHERE posts.path > parent.path AND posts.thread = $1
	ORDER BY posts.path ASC, posts.id ASC;
	`

	querySelectTree = `
	SELECT posts.id, posts.parent, posts.author, posts.message, posts.isedited, posts.forum, posts.thread, posts.created
	FROM posts
	JOIN posts parent ON parent.id = $2
	WHERE posts.path > parent.path AND posts.thread = $1
	ORDER BY posts.path ASC, posts.id ASC
	LIMIT $3;
	`

	querySelectOnPostsParentDesc = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE path[1] = ANY (%s)
	ORDER BY path[1] DESC, path, id;
	`

	querySelectOnPostsParentAsc = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE path[1] = ANY (%s)
	ORDER BY path[1] ASC, path, id;
	`
)
