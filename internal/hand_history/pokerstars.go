package hand_history

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// PokerStarsParser parses PokerStars hand history files
type PokerStarsParser struct {
	handStart  *regexp.Regexp
	gameInfo   *regexp.Regexp
	tableInfo  *regexp.Regexp
	playerInfo *regexp.Regexp
	actionLine *regexp.Regexp
	holeCards  *regexp.Regexp
	boardLine  *regexp.Regexp
	potLine    *regexp.Regexp
	dateTime   *regexp.Regexp
}

// NewPokerStarsParser creates a new PokerStars parser
func NewPokerStarsParser() *PokerStarsParser {
	return &PokerStarsParser{
		handStart:  regexp.MustCompile(`PokerStars Hand #(\d+):`),
		gameInfo:   regexp.MustCompile(`:\s+(.+?)\s+-\s+Level`),
		tableInfo:  regexp.MustCompile(`Table '([^']+)'\s+(\d+)-max`),
		playerInfo: regexp.MustCompile(`Seat (\d+): ([^\(]+) \((\d+(?:\.\d+)?)\s+in chips\)`),
		actionLine: regexp.MustCompile(`^([^:]+):\s+(folds|checks|calls|bets|raises|posts small blind|posts big blind|posts the ante)\s*(\d+(?:\.\d+)?)?\s*(?:to\s+(\d+(?:\.\d+)?))?`),
		holeCards:  regexp.MustCompile(`Dealt to ([^\[]+)\s+\[([^\]]+)\]`),
		boardLine:  regexp.MustCompile(`\*\*\* (FLOP|TURN|RIVER) \*\*\*\s+\[([^\]]+)\]`),
		potLine:    regexp.MustCompile(`Total pot (\d+(?:\.\d+)?)\s*(?:\|\s*Rake\s+(\d+(?:\.\d+)?))?`),
		dateTime:   regexp.MustCompile(`(\d{4}/\d{2}/\d{2}) (\d{1,2}:\d{2}:\d{2})`),
	}
}

// GetSiteName returns "pokerstars"
func (p *PokerStarsParser) GetSiteName() string {
	return "pokerstars"
}

// CanParse checks if the content is from PokerStars
func (p *PokerStarsParser) CanParse(content string) bool {
	return strings.Contains(content, "PokerStars Hand #")
}

// ParseFile parses a PokerStars hand history file
func (p *PokerStarsParser) ParseFile(path string) ([]Hand, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return p.ParseContent(string(content))
}

// ParseContent parses PokerStars hand history content
func (p *PokerStarsParser) ParseContent(content string) ([]Hand, error) {
	var hands []Hand
	var currentHand *Hand
	var currentStreet string
	var actionSequence int
	var rawHandText strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()

		// Check for new hand
		if matches := p.handStart.FindStringSubmatch(line); matches != nil {
			// Save previous hand if exists
			if currentHand != nil {
				currentHand.RawText = rawHandText.String()
				hands = append(hands, *currentHand)
			}

			// Start new hand
			currentHand = &Hand{
				HandID:  matches[1],
				Actions: []Action{},
				Players: []Player{},
			}
			currentStreet = "preflop"
			actionSequence = 0
			rawHandText.Reset()
			rawHandText.WriteString(line + "\n")

			// Parse game type from the same line
			if gameMatches := p.gameInfo.FindStringSubmatch(line); gameMatches != nil {
				currentHand.GameType = strings.TrimSpace(gameMatches[1])
			}

			// Parse date/time
			if dateMatches := p.dateTime.FindStringSubmatch(line); dateMatches != nil {
				dateStr := dateMatches[1] + " " + dateMatches[2]
				parsedTime, err := time.Parse("2006/01/02 15:04:05", dateStr)
				if err == nil {
					currentHand.DateTime = parsedTime
				}
			}

			// Parse stakes (simple extraction)
			if strings.Contains(line, "$") {
				stakesRegex := regexp.MustCompile(`\$?([\d.]+)/\$?([\d.]+)`)
				if stakesMatches := stakesRegex.FindStringSubmatch(line); stakesMatches != nil {
					currentHand.Stakes = stakesMatches[0]
				}
			}

			continue
		}

		if currentHand == nil {
			continue
		}

		rawHandText.WriteString(line + "\n")

		// Parse table info
		if matches := p.tableInfo.FindStringSubmatch(line); matches != nil {
			currentHand.TableName = matches[1]
		}

		// Parse player info
		if matches := p.playerInfo.FindStringSubmatch(line); matches != nil {
			seat, _ := strconv.Atoi(matches[1])
			stack, _ := strconv.ParseFloat(matches[3], 64)
			player := Player{
				Name:  strings.TrimSpace(matches[2]),
				Seat:  seat,
				Stack: stack,
			}
			currentHand.Players = append(currentHand.Players, player)
		}

		// Parse hole cards (identifies hero)
		if matches := p.holeCards.FindStringSubmatch(line); matches != nil {
			currentHand.HeroName = strings.TrimSpace(matches[1])
			cards := strings.Fields(matches[2])
			currentHand.HoleCards = cards
		}

		// Parse streets
		if strings.Contains(line, "*** FLOP ***") {
			currentStreet = "flop"
		} else if strings.Contains(line, "*** TURN ***") {
			currentStreet = "turn"
		} else if strings.Contains(line, "*** RIVER ***") {
			currentStreet = "river"
		} else if strings.Contains(line, "*** SHOW DOWN ***") {
			currentStreet = "showdown"
		}

		// Parse board cards
		if matches := p.boardLine.FindStringSubmatch(line); matches != nil {
			cards := strings.Fields(matches[2])
			currentHand.Board = cards
		}

		// Parse actions
		if matches := p.actionLine.FindStringSubmatch(line); matches != nil {
			playerName := strings.TrimSpace(matches[1])
			actionType := matches[2]
			amount := 0.0

			if len(matches) > 3 && matches[3] != "" {
				amount, _ = strconv.ParseFloat(matches[3], 64)
			}
			if len(matches) > 4 && matches[4] != "" {
				// For raises, use the "to" amount
				amount, _ = strconv.ParseFloat(matches[4], 64)
			}

			action := Action{
				PlayerName: playerName,
				Action:     actionType,
				Amount:     amount,
				Street:     currentStreet,
				Sequence:   actionSequence,
			}
			currentHand.Actions = append(currentHand.Actions, action)
			actionSequence++
		}

		// Parse pot and rake
		if matches := p.potLine.FindStringSubmatch(line); matches != nil {
			pot, _ := strconv.ParseFloat(matches[1], 64)
			currentHand.TotalPot = pot

			if len(matches) > 2 && matches[2] != "" {
				rake, _ := strconv.ParseFloat(matches[2], 64)
				currentHand.Rake = rake
			}
		}

		// Calculate result for hero
		if currentHand.HeroName != "" && strings.Contains(line, currentHand.HeroName) {
			if strings.Contains(line, "collected") {
				collectRegex := regexp.MustCompile(`collected (\d+(?:\.\d+)?)`)
				if collectMatches := collectRegex.FindStringSubmatch(line); collectMatches != nil {
					collected, _ := strconv.ParseFloat(collectMatches[1], 64)
					currentHand.Result += collected
				}
			}
		}
	}

	// Don't forget the last hand
	if currentHand != nil {
		currentHand.RawText = rawHandText.String()

		// Calculate final result (amount won minus amount invested)
		invested := 0.0
		for _, action := range currentHand.Actions {
			if action.PlayerName == currentHand.HeroName && action.Amount > 0 {
				invested += action.Amount
			}
		}
		currentHand.Result -= invested

		hands = append(hands, *currentHand)
	}

	return hands, scanner.Err()
}
