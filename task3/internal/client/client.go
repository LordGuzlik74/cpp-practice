package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"task3/internal/models"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseUrl    string
}

func New(baseUrl string) *Client {
	return &Client{
		baseUrl:    baseUrl,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetMovie(ctx context.Context, id int) (models.Movie, error) {
	url := fmt.Sprintf("%s/%d/info.0.json", c.baseUrl, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return models.Movie{}, err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return models.Movie{}, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return models.Movie{}, models.ErrNotFound
	}
	if res.StatusCode != http.StatusOK {
		return models.Movie{}, fmt.Errorf("error code %d", res.StatusCode)
	}
	var movie models.Movie
	if err := json.NewDecoder(res.Body).Decode(&movie); err != nil {
		return models.Movie{}, err
	}

	return movie, nil
}
