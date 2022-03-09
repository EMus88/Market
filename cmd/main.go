package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EMus88/Market/configs"
	"github.com/EMus88/Market/internal/handler"
	"github.com/EMus88/Market/internal/repository"
	"github.com/EMus88/Market/internal/service"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Internet-shop API
// @version 1.0
// @description API Server for catalog of internet-shop

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	//init logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})
	logger.SetLevel(logrus.DebugLevel)

	//init configs
	if err := configs.InitConfig(); err != nil {
		logger.Fatal(err)
	}
	//db connection
	db, err := repository.NewDB(context.Background())
	if err != nil {
		logger.Fatal("No database connection ")
	}
	logger.Info("DB connection success")
	//migration
	if err := repository.AutoMigration(viper.GetBool("db.migration.isAllowed")); err != nil {
		logger.Fatal(err)
	}

	//init main components
	r := repository.NewRepository(db, logger)
	s := service.NewService(r, logger)
	h := handler.NewHandler(s, logger)

	//init server
	adr := fmt.Sprint(viper.GetString("host"), ":", viper.GetString("port"))
	server := &http.Server{
		Addr:    adr,
		Handler: h.Init(),
	}
	//run server
	go server.ListenAndServe()

	logger.Infof("Server started by address: %s", adr)

	//shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	<-time.After(time.Second * 5)
	logrus.Println("Server stopped")
}
