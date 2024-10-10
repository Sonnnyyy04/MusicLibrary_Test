package service

import (
	"MusicLibrary_Test/internal/models"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
)

type APIService struct {
	ExternalAPIURL string
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func NewAPIService() *APIService {
	return &APIService{
		ExternalAPIURL: os.Getenv("EXTERNAL_API_URL"),
	}
}

func (s *APIService) GetSongDetails(group, song string) (*models.Song, error) {
	endpoint := fmt.Sprintf("%s/info", s.ExternalAPIURL)
	params := url.Values{}
	params.Add("group", group)
	params.Add("song", song)
	reqURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	logrus.Debugf("Requesting external API: %s", reqURL)
	client := http.Client{}
	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var songDetail SongDetail
	err = json.NewDecoder(resp.Body).Decode(&songDetail)
	if err != nil {
		return nil, err
	}
	return &models.Song{
		GroupName:   group,
		SongName:    song,
		Text:        songDetail.Text,
		DateRelease: songDetail.ReleaseDate,
		Link:        songDetail.Link,
	}, nil
}
