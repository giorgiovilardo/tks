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
	normalize := func(name string) string {
		return strings.ToLower(strings.ReplaceAll(name, " ", ""))
	}
	return fmt.Sprintf("%s-%s-%s", normalize(m.HomeTeam), normalize(m.AwayTeam), m.MatchDate.Format(time.RFC3339))
}
