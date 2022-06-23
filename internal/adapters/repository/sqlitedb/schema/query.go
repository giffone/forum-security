package schema

import "github.com/giffone/forum-security/internal/config"

func Query() map[string]string {
	return map[string]string{
		config.QueAttach: `ATTACH DATABASE ? AS ?;`,

		config.QueDetach: `DETACH DATABASE ?;`,

		config.QueRestore: `INSERT INTO %s SELECT * FROM %s.%s;`,

		config.QueInsert2: `INSERT INTO %s (%s, %s)  
		VALUES (?,?);`,

		config.QueInsert3: `INSERT INTO %s (%s, %s, %s)
		VALUES (?,?,?);`,

		config.QueInsert4: `INSERT INTO %s (%s, %s, %s, %s)
		VALUES (?,?,?,?);`,

		config.QueInsert5: `INSERT INTO %s (%s, %s, %s, %s, %s)
		VALUES (?,?,?,?,?);`,

		config.QueInsert6: `INSERT INTO %s (%s, %s, %s, %s, %s, %s)
		VALUES (?,?,?,?,?,?);`,

		config.TabUsers: `CREATE TABLE IF NOT EXISTS %s (
		"id"		INTEGER NOT NULL,
		"login"		TEXT NOT NULL UNIQUE,
		"name"		TEXT NOT NULL,
		"password"	TEXT NOT NULL,
		"email"		TEXT NOT NULL UNIQUE,
		"root"		INTEGER NOT NULL DEFAULT 0,
		"created"	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY("id" AUTOINCREMENT));`,

		config.TabCategories: `CREATE TABLE IF NOT EXISTS %s (
		"id"	INTEGER NOT NULL,
		"body"	TEXT NOT NULL UNIQUE,
		PRIMARY KEY("id"));`,

		config.TabLikes: `CREATE TABLE IF NOT EXISTS %s (
		"id"	INTEGER NOT NULL,
		"body"	TEXT NOT NULL UNIQUE,
		PRIMARY KEY("id"));`,

		config.TabPosts: `CREATE TABLE IF NOT EXISTS %s (
		"id"		INTEGER NOT NULL,
		"user"		INTEGER NOT NULL,
		"title"		TEXT NOT NULL,
		"body"		TEXT NOT NULL,
		"created"	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY("user") REFERENCES "src_users"("id"),
		PRIMARY KEY("id" AUTOINCREMENT));`,

		config.TabFiles: `CREATE TABLE IF NOT EXISTS %s (
		"id"			INTEGER NOT NULL,
		"id_variety"	INTEGER NOT NULL,
		"variety"		TEXT NOT NULL,
		"path"			TEXT NOT NULL,
		"mime"			TEXT NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT));`,

		config.TabPostsLikes: `CREATE TABLE IF NOT EXISTS %s (
		"id"		INTEGER NOT NULL,
		"user"		INTEGER NOT NULL,
		"post"		INTEGER NOT NULL,
		"like"		INTEGER NOT NULL,
		"created"	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY("user") REFERENCES "src_users"("id"),
		FOREIGN KEY("post") REFERENCES "posts"("id"),
		FOREIGN KEY("like") REFERENCES "src_likes"("id"),
		PRIMARY KEY("id" AUTOINCREMENT));`,

		config.TabPostsCategories: `CREATE TABLE IF NOT EXISTS %s (
		"id"		INTEGER NOT NULL UNIQUE,
		"post"		INTEGER NOT NULL,
		"category"	INTEGER NOT NULL,
		FOREIGN KEY("post") REFERENCES "posts"("id"),
		FOREIGN KEY("category") REFERENCES "src_categories"("id"),
		PRIMARY KEY("id" AUTOINCREMENT));`,

		config.TabComments: `CREATE TABLE  IF NOT EXISTS %s (
		"id"		INTEGER NOT NULL,
		"user"		INTEGER NOT NULL,
		"post"		INTEGER NOT NULL,
		"body"		TEXT NOT NULL,
		"created"	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY("user") REFERENCES "src_users"("id"),
		FOREIGN KEY("post") REFERENCES "posts"("id"),
		PRIMARY KEY("id" AUTOINCREMENT));`,

		config.TabCommentsLikes: `CREATE TABLE IF NOT EXISTS %s (
		"id"		INTEGER NOT NULL,
		"user"		INTEGER NOT NULL,
		"comment"	INTEGER NOT NULL,
		"like"		INTEGER NOT NULL,
		"created"	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY("user") REFERENCES "src_users"("id"),
		FOREIGN KEY("comment") REFERENCES "comments"("id"),
		FOREIGN KEY("like") REFERENCES "src_likes"("id"),
		PRIMARY KEY("id" AUTOINCREMENT));`,

		config.TabSessions: `CREATE TABLE IF NOT EXISTS %s (
		"user"		TEXT NOT NULL,
		"uuid"		TEXT NOT NULL UNIQUE,
		"expire"	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP);`,

		config.TabTokens: `CREATE TABLE IF NOT EXISTS %s (
		"id"		TEXT NOT NULL,
		"token"		TEXT NOT NULL UNIQUE,
		"expire"	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP);`,

		config.QueSelect: `SELECT id
		FROM		%s
		WHERE		%s.%s = ?;`,

		config.QueSelectPosts: `SELECT posts.id, posts.title, posts.body, posts.user, src_users.name, posts.created,
			(SELECT	files.path
				FROM	files
				WHERE	posts.id = files.id_variety AND
						files.variety = 'post' AND
						files.mime = 'image'
			) AS path
		FROM		posts
		LEFT JOIN	src_users ON posts.user = src_users.id
		ORDER BY	posts.id DESC;`,

		config.QueSelectUsers: `SELECT src_users.id, src_users.login, src_users.name, src_users.password, src_users.email, src_users.root, src_users.created
		FROM	src_users;`,

		config.QueSelectCategories: `SELECT src_categories.id, src_categories.body
		FROM		src_categories
		ORDER BY	src_categories.id ASC;`,

		config.QueSelectUserBy: `SELECT src_users.id, src_users.login, src_users.name, src_users.password, src_users.email, src_users.root, src_users.created
		FROM	src_users
		WHERE	src_users.%s = ?;`,

		config.QueSelectPostsBy: `SELECT posts.id, posts.title, posts.body, posts.user, src_users.name, posts.created,
			(SELECT	files.path
				FROM	files
				WHERE	posts.id = files.id_variety AND
						files.variety = 'post' AND
						files.mime = 'image'
			) AS path
		FROM		posts
		LEFT JOIN	src_users ON src_users.id = posts.user
		LEFT JOIN	files ON posts.id = files.id_variety
		WHERE		%s.%s = ?
		ORDER BY	posts.id DESC;`,

		config.QueSelectCommentsBy: `SELECT comments.id, src_users.name, comments.body, comments.created, comments.post
		FROM		comments
		LEFT JOIN	src_users ON comments.user = src_users.id
		WHERE		%s.%s = ?
		ORDER BY	comments.id DESC;`,

		config.QueSelectSessionBy: `SELECT sessions.user, sessions.uuid, sessions.expire
		FROM	sessions
		WHERE	sessions.%s = ?
		AND 	sessions.%s = ?;`,

		config.QueSelectCategoryBy: `SELECT src_categories.id, src_categories.body
		FROM		posts_categories
		LEFT JOIN	src_categories ON src_categories.id = posts_categories.category
		WHERE 		posts_categories.%s = ?
		ORDER BY	src_categories.body;`,

		config.QueSelectLikeBy: `SELECT %s.id, %s.like, src_likes.body, %s.created
		FROM		%s
		LEFT JOIN	src_likes ON src_likes.id = %s.like
		WHERE		%s.%s = ?
		AND %s.%s = ?;`,

		config.QueSelectPostsRatedBy: `SELECT posts.id, posts.title, posts.body, posts.user, src_users.name, posts.created, files.path, src_likes.body
		FROM		posts
		LEFT JOIN	src_users ON src_users.id = posts.user
		INNER JOIN	posts_likes ON posts_likes.post = posts.id
		LEFT JOIN	src_likes ON src_likes.id = posts_likes.like
		WHERE		%s.%s = ?
		ORDER BY	posts.id DESC;`,

		config.QueSelectCommentsRatedBy: `SELECT comments.id, src_users.name, comments.body, comments.created, comments.post, src_likes.body
		FROM		comments
		LEFT JOIN	src_users ON src_users.id = comments.user
		INNER JOIN	comments_likes ON comments_likes.comment = comments.id
		LEFT JOIN	src_likes ON src_likes.id = comments_likes.like
		WHERE		%s.%s = ?
		ORDER BY	comments.id DESC`,

		config.QueSelectPostsAndCategoryBy: `SELECT posts.id, posts.title, posts.body, posts.user, src_users.name, posts.created
		FROM		posts
		LEFT JOIN	src_users ON src_users.id = posts.user
		LEFT JOIN	posts_categories ON posts.id = posts_categories.post
		WHERE		%s.%s = ?
		ORDER BY	posts.id DESC;`,

		config.QueSelectLikeCountBy: `SELECT src_likes.id, src_likes.body, COUNT(*) as N
		FROM		%s
		LEFT JOIN	src_likes ON src_likes.id = %s.like
		WHERE 		%s.%s = ?
		GROUP BY	src_likes.body;`,

		config.QueSelectLikedOrNot: `SELECT src_likes.body
		FROM		%s
		LEFT JOIN	src_likes ON src_likes.id = %s.like
		WHERE 		%s.%s = ?
		AND 		%s.%s = ?;`,

		config.QueDeleteBy: `DELETE
		FROM	%s
		WHERE	%s.%s = ?;`,

		config.QueSelectCount: `SELECT COUNT(*) as N
		FROM	%s;`,
	}
}
