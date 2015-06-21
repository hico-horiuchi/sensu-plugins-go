package main

import (
	"database/sql"
	"fmt"
	"regexp"

	"../lib/check"
	_ "github.com/lib/pq"
)

func main() {
	var (
		host     string
		port     int
		database string
		user     string
		password string
	)

	c := check.New("CheckPostgres")
	c.Option.StringVarP(&host, "host", "h", "localhost", "HOST")
	c.Option.IntVarP(&port, "port", "P", 5432, "PORT")
	c.Option.StringVarP(&user, "user", "u", "", "USER")
	c.Option.StringVarP(&password, "password", "p", "", "PASSWORD")
	c.Option.StringVarP(&database, "database", "d", "test", "DATABASE")
	c.Init()

	version, err := selectVersion(host, port, user, password, database)
	if err != nil {
		c.Error(err)
	}

	c.Ok(fmt.Sprint("Server version ", version))
}

func selectVersion(host string, port int, user string, password string, database string) (string, error) {
	var info string

	source := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, database)
	db, err := sql.Open("postgres", source)
	if err != nil {
		return "", err
	}
	defer db.Close()

	err = db.QueryRow("select version()").Scan(&info)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile("PostgreSQL ([0-9\\.]+)")
	return re.FindStringSubmatch(info)[1], nil
}
