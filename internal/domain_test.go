package internal_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/giorgiovilardo/tksgo/internal"
)

func TestMatch_IdempotentKey(t *testing.T) {
	matchDate := time.Date(2023, 4, 15, 15, 0, 0, 0, time.UTC)
	match := internal.Match{
		League:    "Premier League",
		HomeTeam:  "Manchester United",
		AwayTeam:  "Liverpool",
		HomeGoals: 2,
		AwayGoals: 1,
		MatchDate: matchDate,
	}

	expected := "manchesterunited-liverpool-2023-04-15T15:00:00Z"
	result := match.IdempotentKey()

	if result != expected {
		t.Errorf("IdempotentKey() = %v, want %v", result, expected)
	}
}

func TestNormalizeName(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"Manchester United", "manchesterunited"},
		{"Real Madrid", "realmadrid"},
		{"Paris Saint-Germain", "parissaint-germain"},
		{"AFC Bournemouth", "afcbournemouth"},
		{"", ""},
		{"   Spaces   ", "spaces"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := internal.NormalizeName(tc.input)
			if result != tc.expected {
				t.Errorf("NormalizeName(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestAsOdds(t *testing.T) {
	testCases := []struct {
		input    float64
		expected float64
	}{
		{0.5, 2.0},
		{0.25, 4.0},
		{0.1, 10.0},
		{0.75, 1.3333333333333333},
		{1.0, 1.0},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%.2f", tc.input), func(t *testing.T) {
			result := internal.AsOdds(tc.input)
			if result != tc.expected {
				t.Errorf("AsOdds(%f) = %f, want %f", tc.input, result, tc.expected)
			}
		})
	}
}
