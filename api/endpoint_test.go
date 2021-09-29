package api_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"minesweeper_api"
	"minesweeper_api/api"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	validJson = `{
    "num_line": 4,
    "num_column": 4,
    "num_mine": 6
}`
	invalidJson = `{`
	contentType = "application/json"
)

type testRequest struct {
	client      *http.Client
	method      string
	url         string
	contentType string
	token       string
	body        io.Reader
}

type clientServer struct {
	service minesweeper_api.MineSweepService
	server  *httptest.Server
}

func (tr testRequest) make() (*http.Response, error) {
	req, err := http.NewRequest(tr.method, tr.url, tr.body)
	if err != nil {
		return nil, err
	}
	if tr.token != "" {
		req.Header.Set("Authorization", tr.token)
	}
	if tr.contentType != "" {
		req.Header.Set("Content-Type", tr.contentType)
	}
	return tr.client.Do(req)
}

func newService() minesweeper_api.MineSweepService {
	logger, _ := zap.NewDevelopment()
	return minesweeper_api.NewService(logger)
}

func newServer(svc minesweeper_api.MineSweepService) *httptest.Server {
	mux := api.MakeHandler("minesweep", svc)
	return httptest.NewServer(mux)
}

func newClientServer(t *testing.T) clientServer {
	t.Helper()
	mineService := newService()
	mineServer := newServer(mineService)

	return clientServer{
		service: mineService,
		server:  mineServer,
	}
}

func TestCreateAgentGroup(t *testing.T) {
	cli := newClientServer(t)
	defer cli.server.Close()

	cases := map[string]struct {
		req         string
		contentType string
		status      int
		location    string
	}{
		"add a valida game": {
			req:         validJson,
			contentType: contentType,
			status:      http.StatusCreated,
			location:    "/agent_groups",
		},
		"add a valid game with a invalid json": {
			req:         invalidJson,
			contentType: contentType,
			status:      http.StatusBadRequest,
			location:    "/agent_groups",
		},
		"add a valid game without a content type": {
			req:         validJson,
			contentType: "",
			status:      http.StatusUnsupportedMediaType,
			location:    "/agent_groups",
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      cli.server.Client(),
				method:      http.MethodPost,
				url:         fmt.Sprintf("%s/games", cli.server.URL),
				contentType: tc.contentType,
				body:        strings.NewReader(tc.req),
			}
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("unexpected erro %s", err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
	}

}