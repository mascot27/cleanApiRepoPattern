package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/mascot27/cleanApiRepoPattern/middleware"
	"github.com/spf13/viper"
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

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
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
	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.CORS)
	mr := memberRepo.NewMysqlMemberRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	mu := memberUcase.NewMemberUsecase(mr, timeoutContext)
	memberHttpDeliver.NewMemberHttpHandler(e, mu)

	e.Start(viper.GetString("server.address"))
}
