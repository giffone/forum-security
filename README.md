## forum: version secure

**forum based on: "basic forum" + "autentication" + "image upload" + "security"**

![screenshot](/git/GFaiz/forum-security/src/branch/master/screenshot.png)

In this version:

- [x] should implement a Hypertext Transfer Protocol Secure (HTTPS) protocol :
- encrypted connection : 
- - for this you will have to generate an SSL certificate, you can think of this like a identity card for your website. You can create your certificates or use "Certificate Authorities"(CA's)
- [x] the implementation of Rate Limiting must be present on this project
- [x] should encrypt at least the clients passwords. As a Bonus you can also encrypt the database, for this you will have to create a password for your database.
- [x] clients session cookies should be unique. For instance, the session state is stored on the server and the session should present an unique identifier. This way the client has no direct access to it. Therefore, there is no way for attackers to read or tamper with session state.

### Objectives

This project consists in creating a web forum that allows:

- communication between users.
- associating categories to posts.
- liking and disliking posts and comments.
- filtering posts.

### SQLite

In this project we use sqlite db.

### Run

To run project please type in command line 
```console
$ go run ./cmd/forumsqlite/
``` 
or 
```console
$ make run
```

### Randomizer

Delete current db `forum/db/database-sqlite3.db`

To use random `user`, `categories` and `post` for testing, need to uncomment this lines in directory : `forum/internal/app/app.go`:

```go
_, _, schema := repo.ExportSettings()
repository.NewLoremIpsum().Run(db, schema)
```

You can then return comment.

For example random user to login:

> login: `blackbeard`
>
> password: `12345Aa`

### Docker

Run command in command line
```console
$ make docker
```
