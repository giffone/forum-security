package config

import (
	"fmt"
	"log"
	"os"

	"github.com/giffone/forum-security/pkg/read_env"
)

const (
	AdminPath = "secure/db/.env"
	// AdminPath = ".env"

	/*------------------------------------------------------*/

	SqliteDriver = "sqlite3"
	sqlitePort   = ":3306"

	/*------------------------------------------------------*/

	MysqlDriver = "mysql"
	mysqlPort   = ":3306"
)

type DriverConf struct {
	Path, PathB,
	Driver, Port, Connection string
}

func NewSqlite() *DriverConf {
	err := read_env.ReadEnv(AdminPath)
	if err != nil {
		log.Fatalf("read env admin pass: %s", err.Error())
	}
	name := fmt.Sprintf("database-%s.db", SqliteDriver)
	path := fmt.Sprintf("%s/%s", PathDBs, name)
	admin := os.Getenv("ADMIN")
	adminPw := os.Getenv("PASSWORD")
	os.Clearenv()
	return &DriverConf{
		Path:       path,
		PathB:      fmt.Sprintf("%s/%s/%s", PathDBs, PathDBsBackup, name),
		Driver:     SqliteDriver,
		Port:       sqlitePort,
		Connection: fmt.Sprintf("file:%s?_auth&_auth_user=%s&_auth_pass=%s&_auth_crypt=sha256", path, admin, adminPw),
		// file:test.s3db?_auth&_auth_user=admin&_auth_pass=admin&_auth_crypt=sha1
	}
}

func NewMysql() *DriverConf {
	name := fmt.Sprintf("database-%s.db", SqliteDriver)
	return &DriverConf{
		Path:       fmt.Sprintf("%s/%s", PathDBs, name),
		PathB:      fmt.Sprintf("%s/%s/%s", PathDBs, PathDBsBackup, name),
		Driver:     MysqlDriver,
		Port:       mysqlPort,
		Connection: "admin:admin@tcp(localhost:3306)/forum_db", //<username>:<pw>@tcp(<HOST>:<port>)/<dbname>
	}
}
