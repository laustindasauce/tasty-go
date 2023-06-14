package tasty

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	var err error
	client, err = NewClient(http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}
	sessionToken := "default-token-1234"
	client.Session = Session{
		SessionToken: &sessionToken,
	}
	client.baseURL = server.URL
	// Required for customGet method
	client.baseHost = strings.Split(server.URL, "/")[2]
}

func teardown() {
	server.Close()
}
