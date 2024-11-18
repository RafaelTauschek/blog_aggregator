package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RafaelTauschek/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 2 {
		return errors.New("not enough arguments provided")
	}

	name := cmd.arguments[0]
	url := cmd.arguments[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", feed)

	return nil
}
