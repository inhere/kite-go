package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/inhere/kite-go/pkg/gitx"
)

func TestGitHub_CreateAnnotatedTag(t *testing.T) {
	var requests []map[string]any

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Eq(t, http.MethodPost, r.Method)
		assert.Contains(t, r.Header.Get("Authorization"), "token test-token")
		assert.Contains(t, r.Header.Get("Accept"), "application/vnd.github+json")

		var body map[string]any
		err := json.NewDecoder(r.Body).Decode(&body)
		assert.NoErr(t, err)
		requests = append(requests, body)

		switch r.URL.Path {
		case "/repos/owner/repo/git/tags":
			assert.Eq(t, "v1.2.3", body["tag"])
			assert.Eq(t, "release message", body["message"])
			assert.Eq(t, "commit", body["type"])
			assert.Eq(t, "abc123", body["object"])

			tagger, ok := body["tagger"].(map[string]any)
			assert.True(t, ok)
			assert.Eq(t, "kite", tagger["name"])
			assert.Eq(t, "kite@example.com", tagger["email"])
			assert.NotEmpty(t, tagger["date"])

			_ = json.NewEncoder(w).Encode(map[string]any{
				"sha": "tag-object-sha",
			})
		case "/repos/owner/repo/git/refs":
			assert.Eq(t, "refs/tags/v1.2.3", body["ref"])
			assert.Eq(t, "tag-object-sha", body["sha"])

			_ = json.NewEncoder(w).Encode(map[string]any{
				"ref": "refs/tags/v1.2.3",
				"object": map[string]any{
					"sha": "tag-object-sha",
				},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	gh := New(&gitx.Config{HostUrl: "https://github.com"})
	gh.Token = "test-token"
	gh.BaseApi = srv.URL

	resp, err := gh.CreateAnnotatedTag(TagCreateInput{
		RepoPath: "owner/repo",
		Version:  "v1.2.3",
		Message:  "release message",
		Object:   "abc123",
		Tagger: Tagger{
			Name:  "kite",
			Email: "kite@example.com",
		},
	})

	assert.NoErr(t, err)
	assert.Eq(t, "refs/tags/v1.2.3", resp.Ref)
	assert.Eq(t, "tag-object-sha", resp.Object.SHA)
	assert.Len(t, requests, 2)
}

func TestGitHub_CreateAnnotatedTag_requireRepoPath(t *testing.T) {
	gh := New(&gitx.Config{})
	gh.Token = "test-token"

	_, err := gh.CreateAnnotatedTag(TagCreateInput{
		Version: "v1.2.3",
		Message: "release message",
		Object:  "abc123",
	})

	assert.Err(t, err)
	assert.Contains(t, err.Error(), "repo path")
}

func TestGitHub_CreateAnnotatedTag_requireToken(t *testing.T) {
	gh := New(&gitx.Config{})

	_, err := gh.CreateAnnotatedTag(TagCreateInput{
		RepoPath: "owner/repo",
		Version:  "v1.2.3",
		Message:  "release message",
		Object:   "abc123",
	})

	assert.Err(t, err)
	assert.True(t, strings.Contains(err.Error(), "token"))
}
