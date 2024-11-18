## Gator

Gator is a CLI blog aggregator.

## Installation

This CLI requires **Postgres** and **Go** installed to run the program.

- [Download the Postgres installer](https://www.postgresql.org/download/) from the official website and follow the setup instructions for your operating system.
- [Download the Go installer](https://golang.org/doc/install) suitable for your platform and follow the installation guide.

You can then install gator with:

```bash
go install ...
```

## Config

Create a .gatorconfig.json file in your home directory with the following structure:

```json
{
  "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```

Replace the vlaues with your database connection string

## Usage

Create a new user:

```bash
gator register <name>
```

Login

```bash
gator login <name>
```

There are few other commands you'll need:

- gator reset - Resets the database
- gator users - Lists all users
- gator agg <duration> - Start the aggregation
- gator addfeed <feed_name> <url> - Adds a feed
- gator feeds - Lists all feeds
- gator follow <url> - Follow a exsisting feed
- gator following - Lists all feeds the logged in user follows
- gator unfollow <url> - Unfollows a feed
- gator browse <limit> - Lists all posts


