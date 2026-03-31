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

func TestGitHub_GetLatestCommit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Eq(t, http.MethodGet, r.Method)
		assert.Eq(t, "/repos/owner/repo/commits", r.URL.Path)
		assert.Eq(t, "1", r.URL.Query().Get("per_page"))
		assert.Contains(t, r.Header.Get("Authorization"), "token test-token")

		_ = json.NewEncoder(w).Encode([]map[string]any{
			{
				"sha":      "commit-sha-1",
				"html_url": "https://github.com/owner/repo/commit/commit-sha-1",
				"commit": map[string]any{
					"message": "feat: latest commit",
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

	gh := New(&gitx.Config{HostUrl: "https://github.com"})
	gh.Token = "test-token"
	gh.BaseApi = srv.URL

	info, err := gh.GetLatestCommit("owner/repo")
	assert.NoErr(t, err)
	assert.Eq(t, "commit-sha-1", info.SHA)
	assert.Eq(t, "feat: latest commit", info.Commit.Message)
	assert.Eq(t, "kite", info.Commit.Author.Name)
}

func TestGitHub_ListTags(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Eq(t, http.MethodGet, r.Method)
		assert.Eq(t, "/repos/owner/repo/tags", r.URL.Path)
		assert.Eq(t, "5", r.URL.Query().Get("per_page"))
		assert.Contains(t, r.Header.Get("Authorization"), "token test-token")

		_ = json.NewEncoder(w).Encode([]map[string]any{
			{
				"name": "v1.2.3",
				"commit": map[string]any{
					"sha": "sha-123",
					"url": "https://api.github.com/repos/owner/repo/commits/sha-123",
				},
			},
			{
				"name": "v1.2.2",
				"commit": map[string]any{
					"sha": "sha-122",
					"url": "https://api.github.com/repos/owner/repo/commits/sha-122",
				},
			},
		})
	}))
	defer srv.Close()

	gh := New(&gitx.Config{HostUrl: "https://github.com"})
	gh.Token = "test-token"
	gh.BaseApi = srv.URL

	items, err := gh.ListTags(TagListInput{
		RepoPath: "owner/repo",
		Limit:    5,
	})
	assert.NoErr(t, err)
	assert.Len(t, items, 2)
	assert.Eq(t, "v1.2.3", items[0].Name)
	assert.Eq(t, "sha-123", items[0].Commit.SHA)
}
