package slack

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/firba1/util/assert"
)

func jsonTestServer(code int, body string) (*httptest.Server, *http.Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	httpClient := &http.Client{Transport: transport}

	return server, httpClient
}

func TestAPIRequest(t *testing.T) {
	testcases := []struct {
		in  string
		out BasicResponse
		err error
	}{
		{
			`{"ok":true}`,
			BasicResponse{OK: true},
			nil,
		},
		{
			`{"ok":false,"error":"error_key"}`,
			BasicResponse{false, "error_key"},
			SlackError{"error_key"},
		},
	}

	for _, test := range testcases {
		server, client := jsonTestServer(200, test.in)
		defer server.Close()

		api := API{"deadbeef", server.URL, client}
		r := BasicResponse{}
		err := api.request("fart.butt", make(map[string]string), &r)

		assert.Equal(t, r, test.out)
		assert.Equal(t, err, test.err)
	}
}
