package internal

import (
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Leagues []League `koanf:"leagues"`
}

type League struct {
	Name string `koanf:"name"`
	URL  string `koanf:"url"`
}

type Match struct {
	League    string
	HomeTeam  string
	AwayTeam  string
	HomeGoals int
	AwayGoals int
	MatchDate time.Time
}

func (m Match) IdempotentKey() string {
	return fmt.Sprintf("%s-%s-%s", NormalizeName(m.HomeTeam), NormalizeName(m.AwayTeam), m.MatchDate.Format(time.RFC3339))
}

// NormalizeName removes all spaces and converts to lowercase
func NormalizeName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", ""))
}
