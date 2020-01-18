package main

import (
	"log"
	"net/smtp"
	"os"

	"github.com/arthurbailao/suno-kindle/api"
	"github.com/arthurbailao/suno-kindle/database"
	"github.com/arthurbailao/suno-kindle/mail"
	"github.com/arthurbailao/suno-kindle/worker"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func main() {
	db, err := database.Open("my.db")
	if err != nil {
		panic(err)
	}

	mail := mail.New(mail.Config{
		Host:     os.Getenv("SMTP_HOST") + ":" + os.Getenv("SMTP_PORT"),
		Auth:     smtp.PlainAuth("", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST")),
		From:     os.Getenv("MAIL_FROM"),
		FromName: os.Getenv("MAIL_FROMNAME"),
	})

	w := worker.New(db, mail)
	c := api.Controller{DB: db, Worker: w}

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
