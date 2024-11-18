package main

import (
	"context"
	"fmt"
)

func handlerAggregate(s *state, cmd command) error {

	url := "https://www.wagslane.dev/index.xml"

	data, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", data)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		name, err := s.db.GetUserName(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("%s at %s by %s\n", feed.Name, feed.Url, name)
	}

	return nil
}
