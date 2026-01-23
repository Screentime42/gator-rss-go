package main

import (
	"context"
	"fmt"

	"github.com/Screentime42/gator-go/internal/database"
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
		fmt.Println(item.Title)
	}

	return nil
}
