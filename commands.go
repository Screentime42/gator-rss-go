package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"strconv"

	"github.com/Screentime42/gator-go/internal/database"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	cmd map[string]func(*state, command) error
}


func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.cmd[cmd.Name]

	if ok {
		err := handler(s, cmd)
		if err != nil {
			return fmt.Errorf("Command failed: %w", err)
		} 
		return nil
	}
	return fmt.Errorf("Command not found")
}


func (c *commands) register(name string, f func(*state, command) error) {
	c.cmd[name] = f
}

// handlers hooked up to SQL queries
func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no username input")
	}
	
	if len(cmd.Args) > 1 {
		return fmt.Errorf("more than one username input - command only accepts one username")
	}

	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("user: %q doesn't exist", cmd.Args[0])
	}
	
	err = s.cfg.SetUser(cmd.Args[0]) 
	if err != nil {
			return fmt.Errorf("set user unsuccessful: %w", err)
	}
	
	fmt.Println("set user successful!")
	fmt.Printf("current user from state: %s\n", s.cfg.CurrentUserName)
	return nil
	
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no username found")
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("too many arguments")
	}

	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err == nil {
		return fmt.Errorf("user: %q already exists", cmd.Args[0])
	}

	params := database.CreateUserParams{
		CreatedAt:	time.Now().UTC(),
		UpdatedAt:	time.Now().UTC(),
		Name:		 	cmd.Args[0],
	}

	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("user creation unsucessful: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("set user unsuccessful: %w", err)
	}

	fmt.Printf("user: %q created successfully\n", user.Name)
	return nil
}


func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("db reset unsuccessful: %w", err)
	}

	fmt.Println("db reset successful")
	return nil
}


func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("get users unsuccessful: %w", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
	}
	}
	return nil
}


func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
        return fmt.Errorf("incorrect usage - usage: agg <time_between_reqs>")
    }

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
    if err != nil {
        return fmt.Errorf("invalid duration: %w", err)
    }

	 log.Printf("Collecting feeds every %s", timeBetweenRequests)
	
	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		if err := scrapeFeeds(context.Background(), s.db); err != nil {
			log.Println("error scraping feeds:", err)
		}
	}
}


func handlerAddFeed(s *state, cmd command, user database.User) error {
	name, url := cmd.Args[0], cmd.Args[1]
	

	params := database.CreateFeedParams{
		CreatedAt:	time.Now().UTC(),
		UpdatedAt:	time.Now().UTC(),
		Name:			name,
		Url:			url,
		UserID:		user.ID,
	}
	
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		CreatedAt:	time.Now().UTC(),
		UpdatedAt:	time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow feed %w", err)
	}

	fmt.Println("Created feed:", feed.Name)
	return nil
}

func handlerViewFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetUserFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get feeds list: %w", err)
	}

	//Table header
	fmt.Printf("%-25s %-40s %-20s\n", "Feed Name", "URL", "User") 
	fmt.Printf("%s\n", strings.Repeat("-", 90)) 
	
	//Table rows
	for _, f := range feeds { 
		fmt.Printf("%-25s %-40s %-20s\n", f.FeedName, f.FeedUrl, f.UserName) 
	}

	return nil
}


func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("incorrect usage - usage: follow <feed-url>")
	}

	feedURL := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("failed to get feed: %w", err)
	}

	params := database.CreateFeedFollowParams{
		CreatedAt:	time.Now().UTC(),
		UpdatedAt:	time.Now().UTC(),
		UserID:		user.ID,
		FeedID:		feed.ID,
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("failed to create follow: %w", err)
	}
	

	fmt.Printf("You (%s) are now following %s\n", follow.UserName, follow.FeedName)
	return nil
}


func handlerFollowing(s *state, cmd command, user database.User) error {

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Failed to get followed feeds list: %w", err)
	}

	//Table header
	fmt.Printf("%-25s %-40s %-20s\n", "Feed Name", "URL", "User") 
	fmt.Printf("%s\n", strings.Repeat("-", 90)) 
	
	//Table rows
	for _, f := range feeds { 
		fmt.Printf("%-25s %-40s %-20s\n", f.FeedName, f.FeedUrl, f.UserName) 
	}
	return nil
}


func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("invalid usage - usage: unfollow <feed-url>")
	}
	feedURL := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("failed to get feed: %w", err)
	}

	params := database.UnfollowParams {
		UserID: user.ID,
		FeedID: feed.ID,
	}

	if err := s.db.Unfollow(context.Background(), params); err != nil {
		return fmt.Errorf("unable to unfollow: %w", err)
	}
	return nil
}


func handlerBrowse(s *state, cmd command) error {
	// Get current user
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to load current user: %w", err)
	}

	// Limit defined with default being 2
	limit := 2
	if len(cmd.Args) > 0 {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = n
	}

	// Query posts for the user
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID:	user.ID,
		Limit:	int32(limit),
	})
	if err != nil {
		return fmt.Errorf("failed to fetch posts: %w", err)
	}

	// Print posts
	// Table header 
	fmt.Printf("%-30s %-50s %-25s\n", "Title", "URL", "Published") 
	fmt.Printf("%s\n", strings.Repeat("-", 110)) 
	
	// Table rows 
	for _, p := range posts {
		fmt.Printf( "%-30s %-50s %-25s\n", 
		p.Title, p.Url, 
		p.PublishedAt.Format("2006-01-02 15:04"), 
		) 
	}

	return nil
}


