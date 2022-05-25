package main

import (
	"fmt"
	"github.com/LittleLe6owski/LinkShortener/link/delivery/http"
	"github.com/LittleLe6owski/LinkShortener/link/repository/postgres"
	"github.com/LittleLe6owski/LinkShortener/link/repository/ssdb"
	"github.com/LittleLe6owski/LinkShortener/link/usecase"
	"github.com/gocraft/dbr/v2"
	"github.com/labstack/echo/v4"
	"github.com/seefan/gossdb/v2/conf"
	"github.com/seefan/gossdb/v2/pool"
	"github.com/spf13/viper"
	"log"
)

func InitPGConnection() dbr.Session {
	host := viper.GetString("postgres.host")
	port := viper.GetString("postgres.port")
	user := viper.GetString("postgres.user")
	pass := viper.GetString("postgres.pass")
	dbname := viper.GetString("postgres.name")
	sslmode := viper.GetString("postgres.sslmode")

	conn, _ := dbr.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbname, sslmode), nil)
	conn.SetMaxOpenConns(10)
	sess := conn.NewSession(nil)
	_, err := sess.Begin()
	if err != nil {
		log.Fatal("pg connection not begin", err)
	}
	return *sess
}

func InitSSDBConnection() *pool.Client {
	conn := pool.NewConnectors(&conf.Config{
		Host: viper.GetString("ssdb.host"),
		Port: viper.GetInt("ssdb.port"),
	})
	err := conn.Start()
	if err != nil {
		log.Fatal("ssdb connection not started ", err)
	}
	return conn.GetClient()
}

func main() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}
	var (
		pgSess = InitPGConnection()
		client = InitSSDBConnection()
	)

	var (
		pgRepo   = postgres.NewLinkRepository(&pgSess)
		ssdbRepo = ssdb.NewLinkRepository(client)
	)

	var (
		linkUseCase = usecase.NewLinkUseCase(pgRepo, ssdbRepo)
	)

	linkUseCase.InitRestoreLinks()

	e := echo.New()
	http.NewLinkHandler(e, linkUseCase)
	e.Logger.Fatal(e.Start(":8080"))
}
