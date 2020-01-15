package suno

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

const sunoURL = "https://membros.sunoresearch.com.br/"

type Client struct {
	httpClient http.Client
}

type Credentials struct {
	Username, Password string
}

type Report struct {
	ID, Title, URL, Description string
}

type authenticatedRoundTripper struct {
	transport http.RoundTripper
	cookie    *http.Cookie
}

func New(credentials Credentials) (*Client, error) {
	if credentials.Username == "" || credentials.Password == "" {
		return nil, errors.New("failed with empty credentials")
	}

	client := http.Client{
		Timeout: time.Second * 3,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	form := url.Values{}
	form.Set("log", credentials.Username)
	form.Set("pwd", credentials.Password)
	form.Set("wp-submit", "Acessar")

	response, err := client.PostForm(sunoURL+"wp-login.php", form)
	if err != nil {
		return nil, errors.Wrap(err, "login request failed")
	}

	var loggedCookie *http.Cookie

	for _, cookie := range response.Cookies() {
		if strings.HasPrefix(cookie.Name, "wordpress_logged_in_") {
			loggedCookie = cookie
		}
	}

	if loggedCookie == nil {
		return nil, errors.New("authentication failed: login cookie not found")
	}

	return &Client{http.Client{
		Timeout:   time.Second * 3,
		Transport: authenticatedRoundTripper{transport: http.DefaultTransport, cookie: loggedCookie},
	}}, nil
}

func (c Client) Scrape() ([]Report, error) {
	res, err := c.httpClient.Get(sunoURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch home page")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.Wrapf(err, "http request failed with status %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse html")
	}

	var reports []Report
	cards := doc.Find("article.card").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return s.Find("a.download").Length() > 0
	})

	cards.Each(func(_ int, s *goquery.Selection) {
		URL, exists := s.Find("a.download").Attr("href")
		if !exists {
			return
		}

		re := regexp.MustCompile(`\/download\/([a-zA-Z0-9\-]+)`)
		matches := re.FindStringSubmatch(URL)
		if len(matches) < 2 {
			return
		}

		title := strings.TrimSpace(s.Find(".title h2").Text())
		description := strings.TrimSpace(s.Find(".description").Text())

		reports = append(reports, Report{
			ID:          matches[1],
			URL:         URL,
			Title:       title,
			Description: description,
		})
	})
	return reports, nil
}

func (c Client) Download(r Report) error {
	// Get the data
	resp, err := c.httpClient.Get(r.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(r.Title + ".pdf")
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err

}

func (art authenticatedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.AddCookie(art.cookie)
	return art.transport.RoundTrip(req)
}
