package guide

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	apiKey string
	http   *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey, &http.Client{
		Timeout: 2500 * time.Millisecond,
	}}
}

type Hunt struct {
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	Destinations []Destination `json:"destinations"`
}

type Destination struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Description  string   `json:"description"`
	Website      string   `json:"webSite"`
	BannerImages []string `json:"bannerImages"`
	Latitude     float64  `json:"latitude"`
	Longitude    float64  `json:"longitude"`
	Street       string   `json:"street"`
	Suburb       string   `json:"suburb"`
	State        string   `json:"state"`
	PostCode     string   `json:"postcode"`
}

func (c *Client) FindHunt(ctx context.Context, companyKey, region, hunt string) (Hunt, error) {
	url := fmt.Sprintf("https://guide.app/api/v1/regions/%s/hunts/%s?type=Collection", region, hunt)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Hunt{}, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", c.apiKey)
	req.Header.Set("companyKey", companyKey)

	res, err := c.http.Do(req)
	if err != nil {
		return Hunt{}, err
	}
	defer res.Body.Close()

	var h Hunt
	if err := json.NewDecoder(res.Body).Decode(&h); err != nil {
		return Hunt{}, err
	}
	return h, nil
}
