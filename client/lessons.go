package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bootdotdev/bootdev/internal/auth"
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
