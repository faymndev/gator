# gator

Boot.dev's guided aggregator CLI project

Gator lets you follow different RSS feeds in your terminal. 

## Development 

Gator is developed using the Go programming language, and is required to build from source.

PostgreSQL is required to track feeds and posts.

```bash
mise db_migrations && mise db_generate
mise db_migrations_down

# reset the database with 'gator reset'
```

## Usage

Start by creating the necessary configuration files.

```bash
gator init
```

Then, you can create your own user login.

```bash
gator register faymn 
gator login faymn

# list all users with 'gator users'
```

Add some feeds (you automatically follow feeds that you add) and browse their posts!

```bash
gator addfeed wagslane "https://www.wagslane.dev/index.xml" 
gator agg # start aggregating posts
gator browse 10 # see latest 10 posts from feeds that you follow
```
