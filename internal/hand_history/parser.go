package hand_history

import (
	"time"
)

// Parser defines the interface that all site-specific parsers must implement
type Parser interface {
	// CanParse checks if this parser can handle the given content
	CanParse(content string) bool

	// ParseFile parses a hand history file and returns all hands found
	ParseFile(path string) ([]Hand, error)

	// ParseContent parses hand history content and returns all hands found
	ParseContent(content string) ([]Hand, error)

	// GetSiteName returns the name of the poker site this parser handles
	GetSiteName() string
}

// Hand represents a parsed poker hand
type Hand struct {
	HandID    string
	SiteID    int
	GameType  string
	Stakes    string
	TableName string
	DateTime  time.Time
	HeroName  string
	Position  string
	HoleCards []string
	Board     []string
	Actions   []Action
	Players   []Player
	Result    float64
	Rake      float64
	TotalPot  float64
	RawText   string
}

// Action represents a player action in a hand
type Action struct {
	PlayerName string
	Action     string // fold, check, call, bet, raise
	Amount     float64
	Street     string // preflop, flop, turn, river
	Sequence   int
}

// Player represents a player at the table
type Player struct {
	Name     string
	Seat     int
	Stack    float64
	Position string
}

// Manager manages multiple parsers for different poker sites
type Manager struct {
	parsers map[string]Parser
}

// NewManager creates a new parser manager
func NewManager() *Manager {
	m := &Manager{
		parsers: make(map[string]Parser),
	}

	// Register parsers
	m.Register(NewPokerStarsParser())

	return m
}

// Register adds a parser to the manager
func (m *Manager) Register(parser Parser) {
	m.parsers[parser.GetSiteName()] = parser
}

// ParseFile attempts to parse a file using all registered parsers
func (m *Manager) ParseFile(path string) ([]Hand, string, error) {
	// Read file content first
	content, err := readFileContent(path)
	if err != nil {
		return nil, "", err
	}

	// Try each parser
	for siteName, parser := range m.parsers {
		if parser.CanParse(content) {
			hands, err := parser.ParseContent(content)
			return hands, siteName, err
		}
	}

	return nil, "", nil // No parser found
}

// ParseContent attempts to parse content using all registered parsers
func (m *Manager) ParseContent(content string) ([]Hand, string, error) {
	// Try each parser
	for siteName, parser := range m.parsers {
		if parser.CanParse(content) {
			hands, err := parser.ParseContent(content)
			return hands, siteName, err
		}
	}

	return nil, "", nil // No parser found
}

// readFileContent reads the entire file content
func readFileContent(path string) (string, error) {
	// Implementation will use os.ReadFile
	return "", nil
}
