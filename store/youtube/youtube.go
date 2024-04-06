package youtube

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var (
	query          = flag.String("query", "Cricket", "Search term")
	maxResults     = flag.Int64("max-results", 50, "Max YouTube results")
	publishedAfter = flag.String("published-after", "2024-03-01T00:00:00Z", "Fetch videos published after this date-time (RFC3339)")
)

type Video struct {
	Title       string
	Description string
	PublishedAt string
	Thumbnails  string
}

const developerKey = "AIzaSyACfy3-YT_HoHy44aeGwu2BystV-4opKBk"

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
		Type("video"). // Search only for videos
		Order("date").
		MaxResults(*maxResults).        // Order by date
		PublishedAfter(*publishedAfter) // Fetch videos published after the specified date-time

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

	return allVideos, nil
}

func fetchAndStore() error {
	//Fetch latest videos from YouTube API
	videos, err := fetchLatestVideos()
	if err != nil {
		return err
	}

	for _, video := range videos {
		log.Printf("Title: %s\n", video.Title)
		log.Printf("Description: %s\n", video.Description)
		log.Printf("PublishedAt: %s\n", video.PublishedAt)
		log.Printf("Thumbnails: %s\n", video.Thumbnails)
		log.Println()
	}

	log.Println("Videos fetched and stored successfully.")
	return nil
}

func fetchVideoConcurrently(videoch chan *time.Ticker) {
	// Create a ticker to trigger fetching videos at regular intervals
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	if err := fetchAndStore(); err != nil {
		log.Println("Error fetching and storing videos:", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := fetchAndStore(); err != nil {
				log.Println("Error fetching and storing videos:", err)
			}
			videoch <- ticker
		}
	}

}

func FetchAndStoreVideos() {
	videoch := make(chan *time.Ticker)
	go fetchVideoConcurrently(videoch)

	for timer := range videoch {
		fmt.Println(timer)
	}
}
