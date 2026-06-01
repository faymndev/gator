package methods

import (
	"context"
	"errors"
	"fmt"

	"github.com/faymndev/gator/internal/command"
	"github.com/faymndev/gator/internal/database"
)

func FollowFeed(s *command.State, cmd command.Command) error {
	if len(cmd.Args) != 1 {
		return errors.New("must provide a url")
	}

	ctx := context.Background()

	user, err := s.Db.GetUser(ctx, s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	feed, err := s.Db.GetFeedByUrl(ctx, cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to get feed by url: %w", err)
	}

	feed_follow, err := s.Db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}

	fmt.Printf("%s\n%s\n", feed_follow.FeedName, feed_follow.UserName)
	return nil
}

func ListFollowing(s *command.State, cmd command.Command) error {
	ctx := context.Background()

	user, err := s.Db.GetUser(ctx, s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	feeds , err := s.Db.GetFollowing(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get follows: %w", err)
	}

	for _, feed := range feeds {
		fmt.Printf("%s\n", feed.Name)
	}
	return nil
}
