// -----------------------------------------------------------------------------------------------------
// FILE 1: internal/auth/auth.go
// INSTRUCTIONS:
// Copy and paste the ENTIRETY of the code block below into your 'internal/auth/auth.go' file.
// This code exports the functions so the 'client' package can see them.
// -----------------------------------------------------------------------------------------------------
package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func FetchWithAuth(req *http.Request, token string) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("BootDev %s", token))
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	return buf.Bytes(), nil
}

func FetchWithAuthAndPayload(req *http.Request, token string, payload any) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("BootDev %s", token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set the payload
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(payload)
	if err != nil {
		return nil, fmt.Errorf("error encoding payload: %w", err)
	}
	req.Body = url.NopCloser(buf)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	buf = new(bytes.Buffer)
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	return buf.Bytes(), nil
}

// -----------------------------------------------------------------------------------------------------
// FILE 2: client/lessons.go
// INSTRUCTIONS:
// Copy and paste the ENTIRETY of the code block below into your 'client/lessons.go' file.
// This code adds the correct import and references the exported auth functions.
// -----------------------------------------------------------------------------------------------------
package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bootdotdev/bootdev/internal/auth" // Corrected import path
)

type Lesson struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	CourseID  string    `json:"courseId"`
	Path      string    `json:"path"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (c *Client) FetchLessons(token string) ([]Lesson, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/lessons", c.BaseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	body, err := auth.FetchWithAuth(req, token)
	if err != nil {
		return nil, fmt.Errorf("error fetching with auth: %w", err)
	}

	lessons := []Lesson{}
	err = json.Unmarshal(body, &lessons)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling lessons: %w", err)
	}

	return lessons, nil
}

func (c *Client) CreateLesson(token, path, title string) (Lesson, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/lessons", c.BaseURL), nil)
	if err != nil {
		return Lesson{}, fmt.Errorf("error creating request: %w", err)
	}

	body, err := auth.FetchWithAuthAndPayload(req, token, map[string]string{
		"path":  path,
		"title": title,
	})
	if err != nil {
		return Lesson{}, fmt.Errorf("error fetching with auth and payload: %w", err)
	}

	lesson := Lesson{}
	err = json.Unmarshal(body, &lesson)
	if err != nil {
		return Lesson{}, fmt.Errorf("error unmarshaling lesson: %w", err)
	}

	return lesson, nil
}

