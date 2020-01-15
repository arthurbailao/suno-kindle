# suno-kindle

The main goal here is leaning. And that's it. But at the end, SUNO reports will be sent automaticaly to my Kindle. 

## How to use
For now, we are able to scrape our SUNO account and list available reports to download. A simple devices API with create/update and list operations are already in place too.

First of all, setup environment variables and start the application:

```bash
export SUNO_USERNAME=<username>
export SUNO_PASSWORD=<password>
go run cmd/suno-kindle/main.go
```

Then, we can make some basic requests:

```bash
# Scraping SUNO homepage looking for downloadable reports
curl -s -H"Content-Type: application/json" -XPOST http://localhost:8080/process

# Listing registed devices 
curl -s -H"Content-Type: application/json" -XGET http://localhost:8080/devices

# Creating/updating devices
curl -s -H"Content-Type: application/json" -XPUT http://localhost:8080/devices \
-d '{"name":"My Kindle","email":"mykindle@kindle.com"}'
```
