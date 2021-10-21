package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}

func main() {
	router := gin.Default()

	// add the ping-pong function
	router.GET("/ping")

	// add the echo function
	router.POST("/echo", EchoFunction)

	srv := http.Server{
		Addr:         "0.0.0.0:5000",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	logger.Info("the server is now running at port 5000")
	go func() {
		srv.ListenAndServe()
	}()

	logger.Info("if you want to stop it, please enter the Ctrl+C")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("shutdown http server failed", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("the http server is now stopped")
}
func EchoFunction(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}
	c.String(http.StatusOK, string(data))
}
