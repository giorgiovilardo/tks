package internal

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

func LastGoalsHandler(matches []Match) func(c echo.Context) error {
	type lastGoalsRequest struct {
		Team  string `query:"team"`
		Where string `query:"where"`
		Count int    `query:"count"`
	}
	type lastGoals struct {
		Team  string `json:"team"`
		Goals int    `json:"goals"`
	}
	normalizedMatches := lo.Map(matches, func(match Match, _ int) Match {
		return Match{
			League:    match.League,
			HomeTeam:  NormalizeName(match.HomeTeam),
			AwayTeam:  NormalizeName(match.AwayTeam),
			HomeGoals: match.HomeGoals,
			AwayGoals: match.AwayGoals,
			MatchDate: match.MatchDate,
		}
	})
	slices.SortFunc(normalizedMatches, func(a, b Match) int {
		return a.MatchDate.Compare(b.MatchDate)
	})
	return func(c echo.Context) error {
		req := lastGoalsRequest{}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		matchesToCheck := make([]Match, 0)
		if req.Where == "home" {
			matchesToCheck = lo.Filter(normalizedMatches, func(match Match, _ int) bool {
				return match.HomeTeam == req.Team
			})
		} else if req.Where == "away" {
			matchesToCheck = lo.Filter(normalizedMatches, func(match Match, _ int) bool {
				return match.AwayTeam == req.Team
			})
		}

		slices.Reverse(matchesToCheck)
		matchesToCheck = lo.Slice(matchesToCheck, 0, req.Count)
		goals := lo.Sum(lo.Map(matchesToCheck, func(match Match, _ int) int {
			if req.Where == "home" {
				return match.HomeGoals
			}
			return match.AwayGoals
		}))

		fmt.Printf("Matches to check for team %s playing %s, %d count: %d matches - %+v\n", req.Team, req.Where, req.Count, len(matchesToCheck), matchesToCheck)

		return c.JSON(http.StatusOK, lastGoals{Team: req.Team, Goals: goals})
	}
}

func TeamsHandler(matches []Match) func(c echo.Context) error {
	type teamCoso struct {
		Name      string `json:"name"`
		ShortName string `json:"short_name"`
	}
	type teamsResponse struct {
		AllTeams []teamCoso `json:"all_teams"`
	}
	allTeams := make([]string, 0)
	for _, match := range matches {
		allTeams = append(allTeams, match.HomeTeam)
		allTeams = append(allTeams, match.AwayTeam)
	}
	slices.SortFunc(allTeams, func(a, b string) int {
		return strings.Compare(NormalizeName(a), NormalizeName(b))
	})
	allTeams = lo.Uniq(allTeams)
	allTeamsStruct := lo.Map(allTeams, func(team string, _ int) teamCoso {
		return teamCoso{Name: team, ShortName: NormalizeName(team)}
	})
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, teamsResponse{AllTeams: allTeamsStruct})
	}
}
