package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/mascot27/cleanApiRepoPattern/member"
	"github.com/mascot27/cleanApiRepoPattern/middleware"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"log"
	"net/url"
	"os"
	"time"

	memberHttpDeliver "github.com/mascot27/cleanApiRepoPattern/member/delivery/http"
	memberRepo "github.com/mascot27/cleanApiRepoPattern/member/repository"
	memberUcase "github.com/mascot27/cleanApiRepoPattern/member/usecase"
)

func init() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

var mr member.MemberRepository

func main() {

	// choose member repository following the database type
	if viper.GetString(`database_type`) == "mysql" {
		var dbConn *sql.DB
		dbHost := viper.GetString(`database_mysql.host`)
		dbPort := viper.GetString(`database_mysql.port`)
		dbUser := viper.GetString(`database_mysql.user`)
		dbPass := viper.GetString(`database_mysql.pass`)
		dbName := viper.GetString(`database_mysql.name`)
		connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
		val := url.Values{}
		val.Add("parseTime", "1")
		val.Add("loc", "Europe/Paris")
		dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
		dbConn, err := sql.Open(`mysql`, dsn)

		if err != nil && viper.GetBool("debug") {
			fmt.Println(err)
		}
		err = dbConn.Ping()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		defer dbConn.Close()
		mr = memberRepo.NewMysqlMemberRepository(dbConn)
	} else if viper.GetString(`database_type`) == "mongodb" {
		session, err := mgo.Dial("localhost")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		c := session.DB("member")
		defer session.Close()
		mr = memberRepo.NewMongoDbMemberRepository(c)
	}

	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.CORS)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	mu := memberUcase.NewMemberUsecase(mr, timeoutContext)
	memberHttpDeliver.NewMemberHttpHandler(e, mu)

	e.Start(viper.GetString("server.address"))
}
