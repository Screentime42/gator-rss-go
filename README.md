# 🐊 Gator RSS — A Minimalist RSS Aggregator in Go

Gator is a command-line RSS reader and feed aggregator written in Go. It lets you:

- register users
- add RSS feeds
- follow feeds
- periodically fetch posts
- browse your personalized feed

All data is stored in PostgreSQL, and the project uses `sqlc` for type-safe database access.

---

## 🛠️ Requirements

Before using Gator, make sure you have the following installed:

### Go
Download and install Go from:
https://go.dev/dl/

Verify installation:
```
go version
```

### PostgreSQL
Install PostgreSQL from:
https://www.postgresql.org/download/

Verify installation:
```
psql --version
```

Create a database for Gator:
```
CREATE DATABASE gator;
```

---

## 🚀 Installing the Gator CLI

Once Go is installed, you can install the Gator CLI globally using:

```
go install github.com/Screentime42/gator-rss-go/cmd/gator@latest
```

This places the `gator` binary in your Go bin directory, usually:

```
~/go/bin
```

Make sure this directory is on your PATH:

```
export PATH=$PATH:$(go env GOPATH)/bin
```

To make it permanent:

```
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

Verify installation:

```
gator --help
```

---

## ⚙️ Configuration

Gator uses a simple config file to store your database connection string and current user.

Create a config file at:

```
~/.gatorconfig.json
```

Example:

```json
{
  "db_url": "postgres://postgres:password@localhost:5432/gator?sslmode=disable",
  "current_user": ""
}
```

### Run migrations

Gator includes a reset command to rebuild the schema:

```
gator reset
```

This drops and recreates all tables.

---

## 🧑‍💻 Using the CLI

Once installed and configured, you can start using Gator.

### Register a user
```
gator register <username>
```

### Log in
```
gator login <username>
```

### Add a feed
```
gator addfeed "Go Blog" https://go.dev/blog/feed.atom
```

This also automatically follows the feed.

### Follow a feed manually
```
gator follow https://go.dev/blog/feed.atom
```

### List all feeds
```
gator feeds
```

### List who you follow
```
gator following
```

### Run the aggregator
Fetches posts every 30 seconds:
```
gator agg 30s
```

### Browse posts
```
gator browse
```

Or limit results:
```
gator browse 20
```

---

## 🧩 How It Works

### Database
The schema includes:

- users
- feeds
- feed_follows
- posts

Each feed belongs to a user (the creator).  
Each user can follow many feeds.  
Posts are linked to feeds.

### Aggregator
The aggregator:

1. selects the next feed to fetch  
2. downloads the RSS/Atom XML  
3. parses items  
4. inserts new posts  
5. updates `last_fetched_at`

### SQLC
All SQL lives in `sql/queries/`.

Generate code with:
```
sqlc generate
```

---

## 🧪 Development Workflow

### After editing SQL
```
sqlc generate
go clean -cache
go build
```

### After editing migrations
```
gator reset
```

### Debugging DB state
```
psql "$DB_URL"
SELECT * FROM feeds;
SELECT * FROM feed_follows;
SELECT * FROM posts ORDER BY published_at DESC;
```

---
