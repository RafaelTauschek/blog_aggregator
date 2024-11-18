package main

import (
	"context"
	"errors"

	"github.com/RafaelTauschek/internal/database"
)

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return errors.New("no arguments provided")
	}

	url := cmd.arguments[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	err = s.db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}

	return nil
}
