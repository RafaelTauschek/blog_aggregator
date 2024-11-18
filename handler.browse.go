package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/RafaelTauschek/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2

	if len(cmd.arguments) == 1 {
		cmdLimit, err := strconv.Atoi(cmd.arguments[0])
		if err != nil {
			return err
		}
		limit = cmdLimit
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Println("***********************")
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Description: %s\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Printf("From: %s\n", post.PublishedAt.Time)
		fmt.Println("***********************")
	}

	return nil
}
