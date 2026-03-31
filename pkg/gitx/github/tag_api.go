package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultBaseAPI = "https://api.github.com"
)

// Tagger info for create annotated tag.
type Tagger struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Date  string `json:"date,omitempty"`
}

// IsZero checks whether the tagger is empty.
func (t Tagger) IsZero() bool {
	return t.Name == "" && t.Email == "" && t.Date == ""
}

// WithNowDate fills current time when date is empty.
func (t Tagger) WithNowDate() Tagger {
	if t.Date == "" {
		t.Date = time.Now().Format(time.RFC3339)
	}
	return t
}

// TagCreateInput for creating an annotated tag by GitHub API.
type TagCreateInput struct {
	RepoPath string
	Version  string
	Message  string
	Object   string
	Type     string
	Tagger   Tagger
}

type tagObjectRequest struct {
	Tag     string  `json:"tag"`
	Message string  `json:"message"`
	Object  string  `json:"object"`
	Type    string  `json:"type"`
	Tagger  *Tagger `json:"tagger,omitempty"`
}

type tagObjectResponse struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

type refObject struct {
	Type string `json:"type,omitempty"`
	SHA  string `json:"sha"`
	URL  string `json:"url,omitempty"`
}

type refCreateRequest struct {
	Ref string `json:"ref"`
	SHA string `json:"sha"`
}

// TagCreateResponse for created tag ref result.
type TagCreateResponse struct {
	Ref    string    `json:"ref"`
	URL    string    `json:"url,omitempty"`
	Object refObject `json:"object"`
}

func (g *GitHub) apiBase() string {
	if g.BaseApi != "" {
		return strings.TrimRight(g.BaseApi, "/")
	}
	return defaultBaseAPI
}

// CreateAnnotatedTag creates annotated tag and tag ref by GitHub Git Database API.
func (g *GitHub) CreateAnnotatedTag(in TagCreateInput) (*TagCreateResponse, error) {
	if strings.TrimSpace(in.RepoPath) == "" {
		return nil, fmt.Errorf("github: repo path is required")
	}
	if strings.TrimSpace(in.Version) == "" {
		return nil, fmt.Errorf("github: tag version is required")
	}
	if strings.TrimSpace(in.Message) == "" {
		return nil, fmt.Errorf("github: tag message is required")
	}
	if strings.TrimSpace(in.Object) == "" {
		return nil, fmt.Errorf("github: target object sha is required")
	}
	if strings.TrimSpace(g.Token) == "" {
		return nil, fmt.Errorf("github: token is required")
	}

	objType := in.Type
	if objType == "" {
		objType = "commit"
	}

	reqBody := tagObjectRequest{
		Tag:     in.Version,
		Message: in.Message,
		Object:  in.Object,
		Type:    objType,
	}

	if !in.Tagger.IsZero() {
		tagger := in.Tagger.WithNowDate()
		reqBody.Tagger = &tagger
	}

	var tagObj tagObjectResponse
	err := g.postJSON("/repos/"+strings.Trim(in.RepoPath, "/")+"/git/tags", reqBody, &tagObj)
	if err != nil {
		return nil, err
	}

	refReq := refCreateRequest{
		Ref: "refs/tags/" + in.Version,
		SHA: tagObj.SHA,
	}

	var created TagCreateResponse
	err = g.postJSON("/repos/"+strings.Trim(in.RepoPath, "/")+"/git/refs", refReq, &created)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (g *GitHub) postJSON(path string, body any, out any) error {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, g.apiBase()+path, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "token "+g.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := strings.TrimSpace(string(respBody))
		if msg == "" {
			msg = resp.Status
		}
		return fmt.Errorf("github api request failed: %s", msg)
	}

	if out == nil || len(respBody) == 0 {
		return nil
	}
	return json.Unmarshal(respBody, out)
}
