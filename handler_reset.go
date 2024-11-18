package main

import "context"

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}

	err = s.db.DeleteFeeds(context.Background())
	if err != nil {
		return err
	}

	err = s.db.DeleteFeedFollows(context.Background())
	if err != nil {
		return err
	}

	err = s.db.DeletePosts(context.Background())
	if err != nil {
		return err
	}

	return nil
}
