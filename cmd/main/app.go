package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	config "rest-api-tutorial/internal/config"
	"rest-api-tutorial/internal/user"
	"rest-api-tutorial/internal/user/db"
	"rest-api-tutorial/pkg/client/mongodb"
	"rest-api-tutorial/pkg/logging"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Create router")
	router := httprouter.New()

	cfg := config.GetConfig()
	cfgMongo := cfg.MongoDB
	mongoDBClient, err := mongodb.NewClient(context.Background(), cfgMongo.Host, cfgMongo.Port, cfgMongo.Username, cfg.MongoDB.Password, cfg.MongoDB.Database, cfg.MongoDB.AuthDB)
	if err != nil {
		panic(err)
	}
	storage := db.NewStorage(mongoDBClient, cfg.MongoDB.Collection, logger)

	user1 := user.User{
		ID:           "",
		Email:        "raissov1@gmail.com",
		Username:     "raissov",
		PasswordHash: "12345",
	}
	user1ID, err := storage.Create(context.Background(), user1)
	if err != nil {
		panic(err)
	}
	logger.Info(user1ID)

	logger.Info("register user handler")
	handler := user.NewHandler(*logger)
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("starting of application")
	var listener net.Listener
	var listeErr error
	logger.Println(cfg.Listen.Type)
	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("listen unix socket")
		listener, listeErr = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket: %s", socketPath)

	} else {
		logger.Info("listen tcp socket")
		listener, listeErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening in port 1234 %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)

	}
	if listeErr != nil {
		logger.Fatal(listeErr)
	}
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Fatal(server.Serve(listener))
}
