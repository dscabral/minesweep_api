package api_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
    "num_line": 1,
    "num_column": 1,
    "num_mine": 1
}`
	validClickJson = `{
    "line": 1,
    "column": 1,
    "flag": false
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

func TestCreateGame(t *testing.T) {
	cli := newClientServer(t)
	defer cli.server.Close()

	cases := map[string]struct {
		req         string
		contentType string
		status      int
	}{
		"add a valida game": {
			req:         validJson,
			contentType: contentType,
			status:      http.StatusCreated,
		},
		"add a valid game with a invalid json": {
			req:         invalidJson,
			contentType: contentType,
			status:      http.StatusBadRequest,
		},
		"add a valid game without a content type": {
			req:         validJson,
			contentType: "",
			status:      http.StatusUnsupportedMediaType,
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

func TestClickCell(t *testing.T) {
	cli := newClientServer(t)
	defer cli.server.Close()

	id, err := createGame(t, &cli)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		id          string
		req         string
		contentType string
		status      int
	}{
		"click on a valid cell": {
			id:          id,
			req:         validClickJson,
			contentType: contentType,
			status:      http.StatusOK,
		},
		"click on a valid cell with a invalid json": {
			id:          id,
			req:         invalidJson,
			contentType: contentType,
			status:      http.StatusBadRequest,
		},
		"click on a valid cell without a content type": {
			id:          id,
			req:         validClickJson,
			contentType: "",
			status:      http.StatusUnsupportedMediaType,
		},
		"click on a cell in a non-existing board": {
			id:          "wrong",
			req:         validClickJson,
			contentType: contentType,
			status:      http.StatusNotFound,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      cli.server.Client(),
				method:      http.MethodPost,
				url:         fmt.Sprintf("%s/games/%s/click", cli.server.URL, tc.id),
				contentType: tc.contentType,
				body:        strings.NewReader(tc.req),
			}
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("unexpected erro %s", err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
	}
}

func createGame(t *testing.T, cli *clientServer) (string, error) {
	t.Helper()

	ag, err := cli.service.Start(1, 1, 1)
	if err != nil {
		return "", err
	}
	return ag, nil
}
