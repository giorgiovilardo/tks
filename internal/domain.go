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
	League    string    `json:"league"`
	HomeTeam  string    `json:"home_team"`
	AwayTeam  string    `json:"away_team"`
	HomeGoals int       `json:"home_goals"`
	AwayGoals int       `json:"away_goals"`
	MatchDate time.Time `json:"match_date"`
}

func (m Match) IdempotentKey() string {
	return fmt.Sprintf("%s-%s-%s", NormalizeName(m.HomeTeam), NormalizeName(m.AwayTeam), m.MatchDate.Format(time.RFC3339))
}

// NormalizeName removes all spaces and converts to lowercase
func NormalizeName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", ""))
}

// AsOdds returns the probability in odd style
func AsOdds(value float64) float64 {
	return 1 / value
}
