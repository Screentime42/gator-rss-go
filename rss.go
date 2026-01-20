package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {	
		Title				string		`xml:"title"`
		Link				string		`xml:"link"`
		Description		string		`xml:"description"`
		Item				[]RSSItem	`xml:"item"`
	}	`xml:"channel"`
}

type RSSItem struct {
	Title			string	`xml:"title"`
	Link			string	`xml:"link"`
	Description	string	`xml:"description"`
	PubDate		string	`xml:"pubDate"`
}



func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// create http client 
   client := http.Client{
		Timeout: 10 * time.Second,
	}

   // build request with context
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
   
	// set User-Agent header
	req.Header.Set("User-Agent", "gator")

   // do the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

   // read all from resp.Body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

   // unmarshal into an RSSFeed
	var rssFeed RSSFeed
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return nil, err
	}

   // unescape titles/descriptions
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i] = item
	}

   // return the feed
	return &rssFeed, nil
}