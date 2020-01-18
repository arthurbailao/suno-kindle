package worker

import (
	"log"
	"os"
	"time"

	"github.com/arthurbailao/suno-kindle/database"
	"github.com/arthurbailao/suno-kindle/domain"
	"github.com/arthurbailao/suno-kindle/mail"
	"github.com/arthurbailao/suno-kindle/suno"
	"github.com/pkg/errors"
)

// Worker ...
type Worker struct {
	db       *database.DB
	requests chan chan Response
	mail     *mail.Mail
}

// Response ...
type Response struct {
	Reports []suno.Report
	Err     error
}

// New ...
func New(d *database.DB, m *mail.Mail) *Worker {
	return &Worker{
		db:       d,
		requests: make(chan chan Response),
		mail:     m,
	}
}

// Start ...
func (w Worker) Start() error {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case req := <-w.requests:
			client, err := suno.New(suno.Credentials{Username: os.Getenv("SUNO_USERNAME"), Password: os.Getenv("SUNO_PASSWORD")})
			if err != nil {
				req <- Response{nil, errors.Wrap(err, "failed to create suno client")}
				continue
			}

			reports, err := client.Scrape()
			if err != nil {
				req <- Response{nil, errors.Wrap(err, "failed to scrape")}
				continue
			}

			d := domain.Device{Name: "bla", Email: "arthurbailao@gmail.com"}

			if err := client.Download(reports[0]); err != nil {
				req <- Response{nil, errors.Wrap(err, "failed download report")}
				continue
			}

			if err := w.mail.SendMail(d, reports[0]); err != nil {
				req <- Response{nil, errors.Wrap(err, "failed to send mail")}
				continue
			}

			req <- Response{reports, nil}
		case t := <-ticker.C:
			log.Printf("Tick at %s", t)
		}
	}
}

// Call ...
func (w Worker) Call(c chan Response) {
	w.requests <- c
}
