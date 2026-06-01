package methods

import (
	"context"
	"errors"
	"fmt"

	"github.com/faymndev/gator/internal/command"
	"github.com/faymndev/gator/internal/database"
	"github.com/faymndev/gator/internal/feed"
	"github.com/google/uuid"
)

func Aggregate(s *command.State, cmd command.Command) error {
	feed, err := feed.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	fmt.Printf("%+v\n", feed)
	return nil
}

func AddFeed(s *command.State, cmd command.Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return errors.New("must provide a name and a url")
	}

	ctx := context.Background()

	feed, err := s.Db.CreateFeed(ctx, database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   cmd.Args[0],
		Url:    cmd.Args[1],
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}

	_, err = s.Db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}

	fmt.Printf("%+v\n", feed)
	return nil
}

func ListFeeds(s *command.State, cmd command.Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list feeds: %w", err)
	}

	for _, feed := range feeds {
		if !feed.UserName.Valid {
			continue
		}
		fmt.Printf("%s\t%s by %s\n", feed.Name, feed.Url, feed.UserName.String)
	}

	return nil
}
