package main

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"

	_ "github.com/lib/pq"
	"github.com/spf13/pflag"
)

func main() {
	var (
		host     string
		port     int
		database string
		user     string
		password string
	)

	pflag.StringVarP(&host, "host", "h", "localhost", "HOST")
	pflag.IntVarP(&port, "port", "p", 5432, "PORT")
	pflag.StringVarP(&user, "user", "u", "", "USER")
	pflag.StringVarP(&password, "password", "w", "", "PASSWORD")
	pflag.StringVarP(&database, "database", "d", "test", "DATABASE")
	pflag.Parse()

	version := selectVersion(host, port, user, password, database)
	fmt.Println("CheckPostgres OK: Server version", version)
	os.Exit(0)
}

func selectVersion(host string, port int, user string, password string, database string) string {
	var info string

	source := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, database)
	db, err := sql.Open("postgres", source)
	if err != nil {
		fmt.Println("CheckPostgres CRITICAL:", err)
		os.Exit(2)
	}
	defer db.Close()

	err = db.QueryRow("select version()").Scan(&info)
	if err != nil {
		fmt.Println("CheckPostgres CRITICAL:", err)
		os.Exit(2)
	}

	re := regexp.MustCompile("PostgreSQL ([0-9\\.]+)")
	return re.FindStringSubmatch(info)[1]
}
