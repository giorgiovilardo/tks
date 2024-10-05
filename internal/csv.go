package internal

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func randomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36 Edg/98.0.1108.62",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36",
	}
	return userAgents[rand.Intn(len(userAgents))]
}

type leagueCSV struct {
	League League
	Data   string
}

func downloadCsvs(config Config) ([]leagueCSV, error) {
	client := &http.Client{}
	csvs := []leagueCSV{}

	for _, league := range config.Leagues {
		leagueName := league.Name
		req, err := http.NewRequest("GET", league.URL, nil)
		if err != nil {
			return []leagueCSV{}, fmt.Errorf("error creating request for %s: %w", leagueName, err)
		}

		req.Header.Set("User-Agent", randomUserAgent())

		resp, err := client.Do(req)
		if err != nil {
			return []leagueCSV{}, fmt.Errorf("error downloading CSV for %s: %w", leagueName, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return []leagueCSV{}, fmt.Errorf("unexpected status code %d for %s", resp.StatusCode, leagueName)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return []leagueCSV{}, fmt.Errorf("error reading response body for %s: %w", leagueName, err)
		}

		csvs = append(csvs, leagueCSV{League: league, Data: string(body)})
	}

	return csvs, nil
}

func parseCsvs(csvs []leagueCSV) ([]Match, error) {
	var matches []Match

	for _, csvData := range csvs {
		reader := csv.NewReader(strings.NewReader(csvData.Data))

		if _, err := reader.Read(); err != nil {
			return nil, fmt.Errorf("error reading CSV header: %w", err)
		}

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("error reading CSV record: %w", err)
			}

			match, err := parseMatch(record, csvData.League.Name)
			if err != nil {
				return nil, fmt.Errorf("error parsing match: %w", err)
			}

			matches = append(matches, match)
		}
	}

	return matches, nil
}

func parseMatch(record []string, leagueName string) (Match, error) {
	homeGoals, err := strconv.Atoi(record[5])
	if err != nil {
		return Match{}, fmt.Errorf("error parsing home goals: %w", err)
	}

	awayGoals, err := strconv.Atoi(record[6])
	if err != nil {
		return Match{}, fmt.Errorf("error parsing away goals: %w", err)
	}

	dateTimeStr := record[1] + " " + record[2]
	matchDate, err := time.Parse("02/01/2006 15:04", dateTimeStr)
	if err != nil {
		return Match{}, fmt.Errorf("error parsing match date and time: %w", err)
	}

	return Match{
		League:    leagueName,
		HomeTeam:  record[3],
		AwayTeam:  record[4],
		HomeGoals: homeGoals,
		AwayGoals: awayGoals,
		MatchDate: matchDate,
	}, nil
}

func GetMatchesFromCsv(config Config) ([]Match, error) {
	csvs, err := downloadCsvs(config)
	if err != nil {
		return []Match{}, fmt.Errorf("error downloading CSV for %s: %w", config.Leagues[0].Name, err)
	}

	return parseCsvs(csvs)
}
