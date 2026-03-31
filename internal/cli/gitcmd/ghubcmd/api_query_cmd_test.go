package ghubcmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/inhere/kite-go/pkg/gitx"
	ghapi "github.com/inhere/kite-go/pkg/gitx/github"
)

func TestApiLatestCommitOptions_requireRepo(t *testing.T) {
	opts := apiLatestCommitOptions{}

	_, err := opts.toInput()
	assert.Err(t, err)
	assert.Contains(t, err.Error(), "repository path")
}

func TestApiTagListOptions_defaultLimit(t *testing.T) {
	opts := apiTagListOptions{
		RepoPath: "owner/repo",
	}

	in, err := opts.toInput()
	assert.NoErr(t, err)
	assert.Eq(t, 20, in.Limit)
}

func TestApiTagListOptions_capLimit(t *testing.T) {
	opts := apiTagListOptions{
		RepoPath: "owner/repo",
		Limit:    200,
	}

	in, err := opts.toInput()
	assert.NoErr(t, err)
	assert.Eq(t, 100, in.Limit)
}

func TestApiTagAddOptions_resolveLatestSHA(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Eq(t, http.MethodGet, r.Method)
		assert.Eq(t, "/repos/owner/repo/commits", r.URL.Path)
		assert.Eq(t, "1", r.URL.Query().Get("per_page"))

		_ = json.NewEncoder(w).Encode([]map[string]any{
			{
				"sha": "latest-sha-123",
				"commit": map[string]any{
					"message": "feat: latest",
					"author": map[string]any{
						"name":  "kite",
						"email": "kite@example.com",
						"date":  "2026-03-31T10:00:00Z",
					},
				},
			},
		})
	}))
	defer srv.Close()

	gh := ghapi.New(&gitx.Config{HostUrl: "https://github.com"})
	gh.Token = "test-token"
	gh.BaseApi = srv.URL

	opts := apiTagAddOptions{
		RepoPath: "owner/repo",
		Version:  "v1.2.3",
		Message:  "release message",
	}

	in, err := opts.resolveCreateInput(gh)
	assert.NoErr(t, err)
	assert.Eq(t, "latest-sha-123", in.Object)
}
