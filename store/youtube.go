package store

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var (
	query           = flag.String("query", "Cricket", "Search term")
	maxResults      = flag.Int64("max-results", 50, "Max YouTube results")
	lastPublishedAt = time.Now().AddDate(0, 0, -7)
	developerKey    = os.Getenv("YOUTUBE_API_KEY")
)

type Video struct {
	ID          int32
	Title       string
	Description string
	PublishedAt string
	Thumbnails  string
}

func fetchLatestVideos() ([]*Video, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	call := service.Search.List([]string{"id", "snippet"}).
		Q(*query).
		Type("video").
		Order("date").
		MaxResults(*maxResults).
		PublishedAfter(lastPublishedAt.Format(time.RFC3339))

	var allVideos []*Video

	// Fetch all pages of search results
	for {
		response, err := call.Do()
		if err != nil {
			return nil, err
		}

		for _, item := range response.Items {
			if item.Id.Kind == "youtube#video" {
				video := &Video{
					Title:       item.Snippet.Title,
					Description: item.Snippet.Description,
					PublishedAt: item.Snippet.PublishedAt,
					Thumbnails:  item.Snippet.Thumbnails.Default.Url,
				}
				allVideos = append(allVideos, video)

			}
		}

		if response.NextPageToken == "" {
			break
		}

		call.PageToken(response.NextPageToken)
	}

	lastPublishedAt = time.Now()
	return allVideos, nil
}

func FetchAndStore(s *Store) error {
	videos, err := fetchLatestVideos()
	if err != nil {
		return err
	}

	if err := s.driver.StoreVideosInDB(videos); err != nil {
		return err
	}

	log.Println("Videos fetched and stored successfully.")
	return nil
}

func (s *Store) GetVideosFromDB(ctx context.Context, pageNum int, pageSize int) ([]*Video, error) {
	videos, err := s.driver.GetVideo(ctx, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (s *Store) SearchInVideos(ctx context.Context, query string) ([]*Video, error) {
	videos, err := s.driver.SearchVideo(ctx, query)
	if err != nil {
		return nil, err
	}
	return videos, nil
}
