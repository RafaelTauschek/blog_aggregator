package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/RafaelTauschek/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return errors.New("no arguments provided")
	}

	url := cmd.arguments[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return err
	}

	fmt.Printf("%s now follows %s\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	following, err := s.db.GetFeedFollowForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, feed := range following {
		fmt.Printf("%s\n", feed.FeedsName)
	}

	return nil
}
