package main

import (
	"log"

	"github.com/arthurbailao/suno-kindle/api"
	"github.com/arthurbailao/suno-kindle/database"
	"github.com/arthurbailao/suno-kindle/worker"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func main() {
	db, err := database.Open("my.db")
	if err != nil {
		panic(err)
	}

	w := worker.New(db)
	c := api.Controller{db, w}

	router := gin.Default()
	router.GET("/devices", c.ListDevices)
	router.PUT("/devices", c.CreateOrUpdateDevice)
	router.POST("/process", c.Process)

	var g errgroup.Group
	g.Go(func() error {
		return router.Run()
	})
	g.Go(w.Start)

	if err := g.Wait(); err == nil {
		log.Println("Success!")
	}
}
