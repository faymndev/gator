package methods

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/faymndev/gator/internal/command"
	"github.com/faymndev/gator/internal/database"
	"github.com/faymndev/gator/internal/feed"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func Aggregate(s *command.State, cmd command.Command) error {
	ctx := context.Background()

	duration, _ := time.ParseDuration("1m")
	ticker := time.NewTicker(duration)
	fmt.Printf("Collecting feeds every %s\n", duration)
	for ; ; <-ticker.C {
		next_feed, err := s.Db.GetNextFeedToFetch(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch next feed: %w", err)
		}

		rss_feed, err := feed.FetchFeed(ctx, next_feed.Url)
		fmt.Println(rss_feed.Channel.Title)
		for _, item := range rss_feed.Channel.Item {
			fmt.Printf("- %s\n", item.Title)

			publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				return fmt.Errorf("failed to parse published at: %w", err)
			}

			_, err = s.Db.CreatePost(ctx, database.CreatePostParams{
				Title:       item.Title,
				Url:         item.Link,
				Description: item.Description,
				PublishedAt: publishedAt,
				FeedID:      next_feed.ID,
			})
			if err != nil {
				// ignore error if it's just a duplicate url
				if pqErr, ok := errors.AsType[*pq.Error](err); ok && pqErr.Code != "23505" {
					return fmt.Errorf("failed to create post: %w", err)
				}
			}
		}

		updated_at := time.Now().UTC()
		err = s.Db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
			LastFetchedAt: sql.NullTime{Time: updated_at},
			UpdatedAt:     updated_at,
			ID:            next_feed.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to mark feed as fetched: %w", err)
		}
	}
}

func Browse(s *command.State, cmd command.Command, user database.User) error {
	limit := 2
	if len(cmd.Args) >= 1 {
		i, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit specified: %w", err)
		}
		limit = i
	}

	ctx := context.Background()
	feeds, err := s.Db.GetFollowing(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get user feeds: %w", err)
	}

	for _, feed := range feeds {
		posts, err := s.Db.GetPosts(ctx, database.GetPostsParams{
			FeedID: feed.ID,
			Limit:  int32(limit),
		})
		if err != nil {
			return fmt.Errorf("failed to get posts: %w", err)
		}

		fmt.Println(feed.Name)
		for _, post := range posts {
			fmt.Printf("- %s\n", post.Title)
		}
	}

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
