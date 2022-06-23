FROM golang

LABEL project="Forum"

WORKDIR /web

COPY . .

RUN go build --tags sqlite_userauth cmd/forumsqlite/main.go 

CMD ["./main"]