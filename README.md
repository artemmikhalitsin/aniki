# Aniki - Poker Hand Tracker

A cross-platform desktop application for tracking and analyzing poker hand histories, built with Go and Svelte using the Wails framework.

## About

Aniki automatically monitors poker site hand history directories, parses hand history files in real-time, and stores them in a local SQLite database for analysis. It provides insights into hands played, winnings/losses, rake paid, and other statistics.

**Name**: Aniki is lightly themed around gachimuchi culture.

**Development Note**: This project is very transparently mostly vibe-coded with AI assistance. Architecture and design patterns matter, but implementation is done with AI rather than by hand.

## Features

- **Automatic Hand History Monitoring**: Watches configured directories for new hand history files
- **Real-time Parsing**: Asynchronously processes hand histories as they're written by poker clients
- **Multi-Site Support**: Architecture supports multiple poker sites (currently implements PokerStars)
- **Cross-Platform**: Runs on Windows, Linux, and macOS
- **Local SQLite Storage**: File-based persistence without external database requirements
- **Statistics Dashboard**: View aggregate statistics including winnings, rake, win rates
- **Hand History Viewer**: Browse and inspect individual hands with full details
- **Duplicate Detection**: Automatically skips already-processed hands

## Technology Stack

### Backend (Go)

- **Framework**: [Wails v2](https://wails.io/) - Go + Web frontend desktop framework
- **ORM**: [GORM](https://gorm.io/) - Feature-rich ORM for Go
- **Database Driver**: [gorm.io/driver/sqlite](https://github.com/go-gorm/sqlite) - GORM SQLite driver
- **File Watching**: [fsnotify](https://github.com/fsnotify/fsnotify) - Cross-platform file system notifications
- **Architecture**: Repository pattern for clean data access layer
- **Go Version**: 1.24.0

### Frontend (Svelte + TypeScript)

- **Framework**: [Svelte](https://svelte.dev/) with TypeScript
- **Build Tool**: [Vite](https://vitejs.dev/)
- **Styling**: [Tailwind CSS](https://tailwindcss.com/)
- **Bindings**: Auto-generated TypeScript bindings for Go backend methods

## Project Structure

```
aniki/
├── internal/
│   ├── config/          # Platform-specific configuration management
│   │   └── config.go    # Config loading, default paths, OS detection
│   ├── database/        # Database models and connection
│   │   ├── models.go    # GORM models with tags (Hand, Site, Player, Action, Stats)
│   │   └── database.go  # Database initialization and auto-migration
│   ├── repository/      # Repository pattern implementation
│   │   ├── repository.go         # Repository interfaces
│   │   ├── site_repository.go    # Site data access
│   │   ├── hand_repository.go    # Hand data access with filtering
│   │   ├── player_repository.go  # Player data access
│   │   └── action_repository.go  # Action data access
│   ├── parser/          # Hand history parsing
│   │   ├── parser.go    # Parser interface and manager
│   │   └── pokerstars.go # PokerStars-specific parser implementation
│   └── watcher/         # File system monitoring
│       └── watcher.go   # Async file watching with worker pool
├── frontend/            # Svelte TypeScript frontend
│   ├── src/
│   │   ├── lib/         # Svelte components
│   │   │   ├── HandsList.svelte  # Hand history browser
│   │   │   ├── Stats.svelte      # Statistics dashboard
│   │   │   └── Settings.svelte   # Configuration UI
│   │   ├── App.svelte   # Main application component
│   │   └── main.ts      # Entry point
│   └── wailsjs/         # Auto-generated Go bindings
├── app.go               # Main Wails application struct with repositories
├── ORM and Data Access
- **GORM** chosen for type-safe database operations and automatic migrations
- **Repository Pattern** for clean separation of data access from business logic
- Interfaces enable easy mocking for unit tests
- GORM's auto-migration handles schema changes automatically

## Key Design Decisions

### Database Choice
- **modernc.org/sqlite** chosen over mattn/go-sqlite3 for pure Go implementation
- Avoids CGo dependencies, simplifying cross-compilation
- Database stored in platform-specific app data directory

### Async Processing
- 3-worker goroutine pool for concurrent hand history parsing
- 1-second debouncing to avoid processing incomplete files
- Duplicate detection via UNIQUE constraint on (site_id, hand_id)

### Configuration Locations
- **Windows**: `%APPDATA%\Aniki\`
- **macOS**: `~/Library/Application Support/Aniki/`
- **Linux**: `$XDG_CONFIG_HOME/aniki/` (or `~/.config/aniki/`)

### Default PokerStars Paths
- **Windows**: `%LOCALAPPDATA%\PokerStars\HandHistory\`
- **macOS**: `~/Library/Application Support/PokerStars/HandHistory/`
- **Linux**: `~/.wine/drive_c/.../PokerStars/HandHistory/` or `~/.local/share/PokerStars/HandHistory/`

## DGORM Models
All models use GORM struct tags for automatic schema management:

- **Site**: Poker site configurations with unique name constraint
  - Relationships: Has many hands (cascade delete)
- **Hand**: Parsed hand records with composite unique index on (site_id, hand_id)
  - Relationships: Belongs to site, has many players and actions (cascade delete)
  - Indexed: site_id, date_time, hero_name, game_type
- **Player**: Player information per hand
  - Relationships: Belongs to hand
  - Indexed: hand_id
- **Action**: Player actions per hand (ordered by sequence)
  - Relationships: Belongs to hand
  - Indexed: hand_id

### Repository Layer
Clean interfaces for data access:
- **SiteRepository**: CRUD operations for poker sites
- **HandRepository**: Advanced queries, filtering, statistics aggregation
- **PlayerRepository**: Player data access
- **ActionRepository**: Action data access with ordering
- date_time, hero_name, game_type for filtering
- Foreign key relationships for data integrity

## Future Enhancements

- **HUD Functionality**: Display player statistics (PFR, VPIP) in real-time overlay
- **Additional Sites**: Support for GGPoker, 888poker, partypoker
- **Advanced Statistics**: More detailed metrics, graphs, session tracking
- **Hand Replayer**: Visual hand replay with action sequences
- **Export/Import**: Hand history export and database backup/restore

## Development

### Prerequisites
- Go 1.24+
- Node.js 16+
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Platform-Specific Dependencies

**Linux (Debian/Ubuntu)**
```bash
sudo apt install gtk3-devel webkit2gtk-4.0-devel
```

**Linux (Fedora)**
```bash
sudo dnf install gtk3-devel webkit2gtk4.0-devel
```

**macOS**
```bash
# Xcode Command Line Tools required
xcode-select --install
```

**Windows**
- No additional dependencies required

### Live Development

```bash
wails dev
```

Runs Vite dev server with hot reload. Access Go methods via http://localhost:34115 in browser devtools.

### Building

```bash
# Build for current platform
wails build
gorm.io/gorm
gorm.io/driver
# Cross-platform builds
wails build -platform windows/amd64
wails build -platform darwin/amd64
wails build -platform darwin/arm64
wails build -platform linux/amd64
```

## Dependencies

### Go Dependencies
```
github.com/wailsapp/wails/v2
modernc.org/sqlite
github.com/fsnotify/fsnotify
```

### Frontend Dependencies
```
svelte
vite
tailwindcss
postcss
autoprefixer
```

## License

This project is licensed under the **Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License** (CC BY-NC-SA 4.0).

See the [LICENSE](LICENSE) file for details or visit: https://creativecommons.org/licenses/by-nc-sa/4.0/
