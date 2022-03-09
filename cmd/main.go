package main

import (
	"JWT_auth/configs"
	"JWT_auth/internal/handler"
	"JWT_auth/internal/repository"
	"JWT_auth/internal/service"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

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
	<-time.After(time.Second * 2)
	logrus.Println("Server stopped")
}
