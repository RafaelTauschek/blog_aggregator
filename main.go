package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RafaelTauschek/internal/config"
	"github.com/RafaelTauschek/internal/database"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "gator")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	exsits, ok := c.commands[cmd.name]

	if !ok {
		return errors.New("no command found")
	}

	err := exsits(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arguments) != 2 {
		return errors.New("not enough arguments provided")
	}

	name := cmd.arguments[0]
	url := cmd.arguments[1]

	userID, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
	if err != nil {
		return err
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    userID.ID,
	})

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", feed)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("no arguments provided")
	}

	name := cmd.arguments[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})

	if err != nil {
		return err
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}

	err = s.db.DeleteFeeds(context.Background())
	if err != nil {
		return err
	}
	return nil
}

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

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())

	if err != nil {
		return err
	}

	for _, user := range users {

		if user.Name == s.cfg.CurrentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("no argument provided")
	}
	username := cmd.arguments[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Println("user has been set")

	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	s := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := &commands{
		commands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAggregate)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Not enough arguments provided")
		return
	}

	cmdName := os.Args[1]
	arg := os.Args[2:]

	cmd := command{
		name:      cmdName,
		arguments: arg,
	}

	err = cmds.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
