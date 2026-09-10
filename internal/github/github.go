// Package github предоставляет доступ к GitHub Issues API
// Для упр. Донован 4.10
package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // текст в формате Markdown
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

// SearchIssues отправляет запрос к GitHub Issues API.
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))

	req, err := http.NewRequest("GET", IssuesURL+"?q="+q, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search issues: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

func GroupByAge(issues []*Issue, now time.Time) (recent, thisYear, old []*Issue) {
	monthAgo := now.AddDate(0, -1, 0)
	yearAgo := now.AddDate(-1, 0, 0)

	for _, issue := range issues {

		if issue == nil {
			continue
		}

		if issue.CreatedAt.After(monthAgo) {
			recent = append(recent, issue)
		} else if issue.CreatedAt.After(yearAgo) {
			thisYear = append(thisYear, issue)
		} else {
			old = append(old, issue)
		}
	}
	return recent, thisYear, old

}
