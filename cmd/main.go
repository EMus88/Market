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
		logrus.Fatal(err)
	}
	//db connection
	db, err := repository.NewDB(context.Background())
	if err != nil {
		logrus.Fatal("No database connection ")
	}
	//migration
	if err := repository.AutoMigration(viper.GetBool("db.migration.isAllowed")); err != nil {
		logrus.Fatal(err)
	}

	//init main components
	r := repository.NewRepository(db)
	s := service.NewService(r)
	h := handler.NewHandler(s)

	//init server
	adr := fmt.Sprint(viper.GetString("host"), ":", viper.GetString("port"))
	server := &http.Server{
		Addr:    adr,
		Handler: h.Init(),
	}
	//run server
	go server.ListenAndServe()

	logrus.Printf("Server started by address: %s", adr)

	//shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Println("Server stopped")
}
