package internal

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/DamienDuv/rssagg/internal/database"
	"github.com/google/uuid"
)

func StartScraping(
	db database.Querier,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Starting on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("Error while fetching feeds", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db database.Querier, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	log.Printf("Scraping %s", feed.Url)
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error while marking feed as fetched: %s", err)
		return
	}

	rssFeed, err := UrlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error while scraping %s: %s", feed.Url, err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description = sql.NullString{
				String: item.Description,
				Valid:  true,
			}
		}

		t, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error while parsing date %s: %s", item.PubDate, err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Description: description,
			PublishedAt: t,
			Url:         item.Link,
			FeedID:      feed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			
			log.Printf("Error while creating post: %s", err)
			continue
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
