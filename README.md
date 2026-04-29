# linkmgr - Terminal-Based Bookmark Manager

A fast, lightweight, terminal-based bookmark management application written in Go. Store and manage your bookmarks in a MySQL database with an intuitive command-line interface.

## Features

- **Fast & Lightweight**: Single binary with no external dependencies (except MySQL database)
- **Full CRUD Operations**: Create, read, update, and delete bookmarks
- **Tag Management**: Organize bookmarks with tags
- **Powerful Search**: Fulltext search across titles, excerpts, and content
- **Flexible Filtering**: Filter bookmarks by tags, public status, and more
- **Clean CLI**: Built with Cobra for a familiar command-line experience
- **Secure Configuration**: Store credentials securely with proper file permissions

## Installation

### From Source

```bash
git clone https://github.com/jvzantvoort/linkmgr
cd linkmgr
make build
sudo make install
```

### Using Go

```bash
go install github.com/jvzantvoort/linkmgr/cmd/linkmgr@latest
```

## Quick Start

### 1. Initialize Configuration

```bash
linkmgr init
```

This will prompt you for database connection settings and create a configuration file at `~/.linkmgr.yaml`.

### 2. Add Bookmarks

```bash
# Simple bookmark
linkmgr add https://golang.org

# With metadata
linkmgr add https://golang.org --title "Go Programming Language" --tags "programming,golang"

# With full details
linkmgr add https://example.com \
  --title "Example Site" \
  --tags "example,demo" \
  --excerpt "A demonstration website" \
  --public
```

### 3. List Bookmarks

```bash
# List all bookmarks
linkmgr list

# Filter by tag
linkmgr list --tag golang

# Limit results
linkmgr list --limit 20

# Pagination
linkmgr list --limit 10 --offset 10
```

### 4. Search Bookmarks

```bash
# Fulltext search
linkmgr search "kubernetes tutorial"

# Search with limit
linkmgr search "golang" --limit 20
```

## Usage

### Commands

#### init - Initialize Configuration

```bash
linkmgr init                    # Interactive configuration setup
linkmgr init --test-connection  # Test existing database connection
```

#### add - Add Bookmark

```bash
linkmgr add <url> [flags]

Flags:
  -t, --title string     Bookmark title
      --tags string      Comma-separated tags
  -e, --excerpt string   Brief description
  -a, --author string    Author name
  -p, --public           Mark as public
```

#### list - List Bookmarks

```bash
linkmgr list [flags]

Flags:
      --tag string       Filter by tag
  -n, --limit int        Maximum number of results
      --offset int       Result offset for pagination
      --public-only      Show only public bookmarks
```

#### search - Search Bookmarks

```bash
linkmgr search <query> [flags]

Flags:
  -n, --limit int        Maximum number of results
```

#### version - Version Information

```bash
linkmgr version
```

### Global Flags

```bash
--config string   Config file (default: $HOME/.linkmgr.yaml)
-v, --verbose     Verbose output
```

## Configuration

Configuration is stored in `~/.linkmgr.yaml` (or specified via `--config` flag).

Example configuration:

```yaml
database:
  host: localhost
  port: 3306
  user: your_database_user
  password: your_database_password
  database: shiori
  maxConnections: 10

display:
  defaultLimit: 50
  dateFormat: "2006-01-02 15:04"
```

### Environment Variables

Configuration can also be set via environment variables with the `LINKMGR_` prefix:

```bash
export LINKMGR_DATABASE_HOST=localhost
export LINKMGR_DATABASE_PORT=3306
export LINKMGR_DATABASE_USER=myuser
export LINKMGR_DATABASE_PASSWORD=mypassword
export LINKMGR_DATABASE_DATABASE=shiori
```

## Database Setup

linkmgr requires a MySQL or MariaDB database with the schema defined in `schema.sql`. The application is compatible with the [Shiori](https://github.com/go-shiori/shiori) bookmark manager schema.

### Creating the Database

```bash
# Create database
mysql -u root -p -e "CREATE DATABASE shiori CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"

# Import schema
mysql -u root -p shiori < schema.sql

# Create user (optional)
mysql -u root -p -e "CREATE USER 'linkmgr'@'localhost' IDENTIFIED BY 'your_password';"
mysql -u root -p -e "GRANT ALL PRIVILEGES ON shiori.* TO 'linkmgr'@'localhost';"
mysql -u root -p -e "FLUSH PRIVILEGES;"
```

## Development

### Prerequisites

- Go 1.21 or later
- MySQL 5.7+ or MariaDB 10.11+
- Make (optional, for build automation)

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Format code
make fmt

# Run linter
make lint

# All checks (format, vet, test)
make check
```

### Project Structure

```
linkmgr/
├── cmd/linkmgr/          # Application entry point
├── internal/
│   ├── cmd/              # CLI commands (Cobra)
│   ├── config/           # Configuration management (Viper)
│   ├── models/           # Data models
│   ├── repository/       # Database layer
│   ├── display/          # Output formatting
│   └── validator/        # Input validation
├── tests/
│   ├── integration/      # Integration tests
│   └── fixtures/         # Test data
├── Makefile              # Build automation
├── go.mod                # Go dependencies
└── README.md             # This file
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Format code (`make fmt`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## License

MIT License - see LICENSE file for details

## Acknowledgments

- Compatible with [Shiori](https://github.com/go-shiori/shiori) bookmark manager database schema
- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Configuration management via [Viper](https://github.com/spf13/viper)

## Support

For issues, questions, or contributions, please visit the [GitHub repository](https://github.com/jvzantvoort/linkmgr).
