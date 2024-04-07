package store

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var (
	query           = flag.String("query", "Cricket", "Search term")
	maxResults      = flag.Int64("max-results", 50, "Max YouTube results")
	lastPublishedAt = time.Now().AddDate(0, 0, -7)
)

type Video struct {
	ID          int32
	Title       string
	Description string
	PublishedAt string
	Thumbnails  string
}

const developerKey = "AIzaSyABh7mwq9jLyiWaIjEWSX7gTsJXgg_8XZA"

func fetchLatestVideos() ([]*Video, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
	fmt.Println(lastPublishedAt.Format(time.RFC3339))
	fmt.Println(lastPublishedAt)
	call := service.Search.List([]string{"id", "snippet"}).
		Q(*query).
		Type("video"). // Search only for videos
		Order("date").
		MaxResults(*maxResults).                             // Order by date
		PublishedAfter(lastPublishedAt.Format(time.RFC3339)) // Fetch videos published after the specified date-time

	// Group video, channel, and playlist results in separate lists.
	var allVideos []*Video

	// Fetch all pages of search results
	for {
		response, err := call.Do()
		if err != nil {
			log.Fatalf("Error making search API call: %v", err)
		}

		// Iterate through each item and add it to the videos slice.
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
	//Fetch latest videos from YouTube API
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
