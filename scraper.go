package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"database/sql"

	"github.com/Screentime42/gator-go/internal/database"
	"github.com/google/uuid"
)

func scrapeFeeds(ctx context.Context, q *database.Queries) error {
	nextFeed, err := q.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("failed to get next feed: %w", err)
	}

	err = q.MarkFeedFetched(ctx, nextFeed.ID) 
	if err != nil {
		return fmt.Errorf("failed to mark feed as fetched: %w", err)
		}
	
	feedData, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}
	
	for _, item := range feedData.Channel.Item {
		publishedAt, err := parsePubDate(item.PubDate)
		if err != nil {
			log.Printf("failed to parse pubDate %q: %v", item.PubDate, err)
			continue
		}

		params := database.CreatePostParams{
			CreatedAt:		time.Now().UTC(),
			UpdatedAt: 		time.Now().UTC(),
			Title:			item.Title,
			Url:				item.Link,
			Description: sql.NullString{ 
				String: item.Description, 
				Valid: item.Description != "", 
			},
			PublishedAt:	publishedAt,
			FeedID:			nextFeed.ID,

		}
		post, err := q.CreatePost(ctx, params)
		if err != nil {
			log.Printf("unexpected error creating post: %v", err)
			continue
		}

		if post.ID == uuid.Nil {
			continue
		}
	}

	return nil
}


func parsePubDate(raw string) (time.Time, error) {
    layouts := []string{
        time.RFC1123Z,
        time.RFC1123,
        time.RFC822Z,
        time.RFC822,
        time.RFC3339,
        "Mon, 02 Jan 2006 15:04:05 MST",
        "02 Jan 2006 15:04:05 MST",
    }

    for _, layout := range layouts {
        if t, err := time.Parse(layout, raw); err == nil {
            return t.UTC(), nil
        }
    }

    return time.Time{}, fmt.Errorf("unrecognized time format: %q", raw)
}
