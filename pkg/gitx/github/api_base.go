package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseAPI = "https://api.github.com"
)

func (g *GitHub) apiBase() string {
	if g.BaseApi != "" {
		return strings.TrimRight(g.BaseApi, "/")
	}
	return defaultBaseAPI
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

	g.setAPIHeaders(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return decodeJSONResponse(resp, out)
}

func (g *GitHub) getJSON(path string, query url.Values, out any) error {
	apiURL := g.apiBase() + path
	if len(query) > 0 {
		apiURL += "?" + query.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return err
	}

	g.setAPIHeaders(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return decodeJSONResponse(resp, out)
}

func (g *GitHub) setAPIHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "token "+g.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
}

func decodeJSONResponse(resp *http.Response, out any) error {
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
