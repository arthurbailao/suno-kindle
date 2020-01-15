package worker

import (
	"log"
	"os"
	"time"

	"github.com/arthurbailao/suno-kindle/database"
	"github.com/arthurbailao/suno-kindle/suno"
	"github.com/pkg/errors"
)

type Worker struct {
	db       *database.DB
	requests chan chan Response
}

type Response struct {
	Reports []suno.Report
	Err     error
}

func New(d *database.DB) *Worker {
	return &Worker{
		db:       d,
		requests: make(chan chan Response),
	}
}

func (w Worker) Start() error {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case req := <-w.requests:
			client, err := suno.New(suno.Credentials{os.Getenv("SUNO_USERNAME"), os.Getenv("SUNO_PASSWORD")})
			if err != nil {
				req <- Response{nil, errors.Wrap(err, "failed to create suno client")}
				continue
			}

			reports, err := client.Scrape()
			if err != nil {
				req <- Response{nil, errors.Wrap(err, "failed to scrape")}
				continue
			}

			req <- Response{reports, nil}
		case t := <-ticker.C:
			log.Printf("Tick at %s", t)
		}
	}
}

func (w Worker) Call(c chan Response) {
	w.requests <- c
}
