package rep

// users
const (
	SelectUserByNickname = `
	SELECT nickname, fullname, about, email
	FROM users
	WHERE nickname = $1
	LIMIT 1;
	`

	SelectUserByEmailOrNickname = `
	SELECT nickname, fullname, about, email
	FROM users
	WHERE nickname = $1 or email = $2
	LIMIT 2;
	`

	GetUsersOfForumDescNotNilSince = `
	SELECT nickname, fullname, about, email
	FROM users_forum
	WHERE slug = $1 and nickname < '%s'
	ORDER BY nickname DESC
	LIMIT nullif($2, 0);
	`

	GetUsersOfForumDescSinceNil = `
	SELECT nickname, fullname, about, email
	FROM users_forum
	WHERE slug = $1
	ORDER BY nickname DESC
	LIMIT nullif($2, 0);
	`

	GetUsersOfForumDescNil = `
	SELECT nickname, fullname, about, email
	FROM users_forum
	WHERE slug = $1 and nickname > '%s'
	ORDER BY nickname
	LIMIT nullif($2, 0);
	`

	UpdateUser = `
	UPDATE users
	SET
		fullname = coalesce(nullif($1, ''), fullname),
		about = coalesce(nullif($2, ''), about),
		email = coalesce(nullif($3, ''), email)
	WHERE nickname = $4
	RETURNING nickname, fullname, about, email;
	`

	InsertUser = `
	INSERT INTO users (nickname, fullname, about,email)
	VALUES ($1, $2, $3, $4);
	`
)

// thread
const (
	SelectThread = `
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE slug = $1
	LIMIT 1;
	`

	SelectThreadSlug = `
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE slug = $1
	LIMIT 1;
	`
	GetThreadsSinceDescNotNil = `
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE forum = $1 and created <= $2
	ORDER BY created DESC
	LIMIT $3;
	`

	GetThreadsSinceDescNil = `
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE forum = $1 AND created >= $2
	ORDER BY created ASC
	LIMIT $3;
	`

	GetThreadsDescNotNil = `
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE forum = $1
	ORDER BY created DESC
	LIMIT $2;
	`

	GetThreadsDescNil = `
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE forum = $1
	ORDER BY created ASC
	LIMIT $2;
	`

	SelectThreadId = `
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE id = $1
	LIMIT 1;
	`

	InsertThread = `
	INSERT INTO threads (author, message, title, created, forum, slug, votes)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id;
	`

	UpdateThread = `
	UPDATE threads
	SET
		title = coalesce(nullif($1, ''), title),
		message = coalesce(nullif($2, ''), message)
	WHERE %s = $3
	RETURNING id, title, author, forum, message, votes, slug, created;
	`

	SelectThreadShort = `
	SELECT slug, author
	FROM threads
	WHERE slug = $1;
	`
)

// votes
const (
	UpdateVote = `
	UPDATE votes
	SET voice = $1
	WHERE author = $2 AND thread = $3;
	`

	InsertVote = `
	INSERT INTO votes (author, voice, thread)
	VALUES ($1, $2, $3);
	`
)

// forum
const (
	SelectForumBySlug = `
	SELECT title, "user", slug, posts, threads
	FROM forum
	WHERE slug = $1
	LIMIT 1;
	`

	InsertInForum = `
	INSERT INTO forum (slug, "user", title)
	VALUES ($1, $2, $3);
	`

	SelectSlugFromForum = `
	SELECT slug
	FROM forum
	WHERE slug = $1;
	`
)

// other
const (
	ClearAll = `
	TRUNCATE table users, forum, threads, posts, votes, users_forum CASCADE;
	`

	SelectCountUsers = `
	SELECT COUNT(*) FROM users;
	`

	SelectCountForum = `
	SELECT COUNT(*) FROM forum;
	`

	SelectCountThreads = `
	SELECT COUNT(*) FROM threads;
	`

	SelectCountPosts = `
	SELECT COUNT(*) FROM posts;
	`
)

// posts
const (
	SelectPostById = `
	SELECT parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE id = $1;
	`

	UpdatePostMessage = `
	UPDATE posts
	SET
		message = coalesce(nullif($1, ''), message),
		isedited = case WHEN $1 = '' OR message = $1 THEN isedited ELSE true END
	WHERE id = $2
	RETURNING id, parent, author, message, isedited, forum, thread, created, path;
	`

	SelectPostSinceDescNotNil = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE thread = $1
	ORDER BY id DESC
	LIMIT $2;
	`

	SelectPostSinceDescNil = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE thread = $1
	ORDER BY id
	LIMIT $2;
	`

	SelectPostDescNotNil = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE thread = $1 and id < $2
	ORDER BY id DESC
	LIMIT $3;
	`

	SelectPostDescNil = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE thread = $1 and id > $2
	ORDER BY id
	LIMIT $3;
	`

	InsertIntoPosts = `
	INSERT INTO posts (author, created, forum, message, parent, thread)
	VALUES %s
	RETURNING id, created, forum, isedited, thread;
	`

	SelectTreeLimitSinceNil = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE thread = $1
	ORDER BY path, id DESC;
	`

	SelectTreeLimitSinceDescNil = `
	SELECT id, parent, author, message, isedited, forum, thread, created FROM posts
	WHERE thread = $1
	ORDER BY path, id ASC;
	`

	SelectTreeSinceNil = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE thread = $1
	ORDER BY path DESC, id DESC
	LIMIT $2;
	`

	SelectTreeSinceDescNil = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE thread = $1
	ORDER BY path, id ASC
	LIMIT $2;
	`

	SelectTreeNotNil = `
	SELECT posts.id, posts.parent, posts.author, posts.message, posts.isedited, posts.forum, posts.thread, posts.created
	FROM posts
	JOIN posts parent ON parent.id = $2
	WHERE posts.path < parent.path AND posts.thread = $1
	ORDER BY posts.path DESC, posts.id DESC
	LIMIT $3;
	`

	SelectTreeSinceNilDesc = `
	SELECT posts.id, posts.parent, posts.author, posts.message, posts.isedited, posts.forum, posts.thread, posts.created
	FROM posts
	JOIN posts parent ON parent.id = $2
	WHERE posts.path < parent.path AND posts.thread = $1
	ORDER BY posts.path DESC, posts.id DESC;
	`

	SelectTreeSinceNilDescNil = `
	SELECT posts.id, posts.parent, posts.author, posts.message, posts.isedited, posts.forum, posts.thread, posts.created
	FROM posts
	JOIN posts parent ON parent.id = $2
	WHERE posts.path > parent.path AND posts.thread = $1
	ORDER BY posts.path ASC, posts.id ASC;
	`

	SelectTree = `
	SELECT posts.id, posts.parent, posts.author, posts.message, posts.isedited, posts.forum, posts.thread, posts.created
	FROM posts
	JOIN posts parent ON parent.id = $2
	WHERE posts.path > parent.path AND posts.thread = $1
	ORDER BY posts.path ASC, posts.id ASC
	LIMIT $3;
	`

	SelectOnPostsParentDesc = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE path[1] = ANY (%s)
	ORDER BY path[1] DESC, path, id;
	`

	SelectOnPostsParentAsc = `
	SELECT id, parent, author, message, isedited, forum, thread, created
	FROM posts
	WHERE path[1] = ANY (%s)
	ORDER BY path[1] ASC, path, id;
	`
)
