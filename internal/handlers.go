package internal

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type lastGoalsRequest struct {
	Team  string `query:"team"`
	Where string `query:"where"`
	Count int    `query:"count"`
	Type  string `query:"type"`
}
type lastGoals struct {
	Team      string `json:"team"`
	HomeGoals int    `json:"home_goals"`
	AwayGoals int    `json:"away_goals"`
}

func lastGoalsService(matches []Match, req lastGoalsRequest) lastGoals {
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
	homeGoals := lo.Sum(lo.Map(matchesToCheck, func(match Match, _ int) int {
		if req.Where == "home" {
			return match.HomeGoals
		}
		return match.AwayGoals
	}))
	awayGoals := lo.Sum(lo.Map(matchesToCheck, func(match Match, _ int) int {
		if req.Where == "home" {
			return match.AwayGoals
		}
		return match.HomeGoals
	}))

	return lastGoals{Team: req.Team, HomeGoals: homeGoals, AwayGoals: awayGoals}
}

func LastGoalsHandler(matches []Match) func(c echo.Context) error {
	return func(c echo.Context) error {
		req := lastGoalsRequest{}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, lastGoalsService(matches, req))
	}
}

func LastGoalsHtmlHandler(matches []Match) func(c echo.Context) error {
	return func(c echo.Context) error {
		req := lastGoalsRequest{}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		result := lastGoalsService(matches, req)
		if req.Type == "scored" {
			return c.HTML(http.StatusOK, fmt.Sprintf("%d", result.HomeGoals))
		}
		return c.HTML(http.StatusOK, fmt.Sprintf("%d", result.AwayGoals))
	}
}

type teamResponse struct {
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}

type teamsResponse struct {
	AllTeams []teamResponse `json:"all_teams"`
}

func teamsService(matches []Match) teamsResponse {
	allTeams := make([]string, 0)
	for _, match := range matches {
		allTeams = append(allTeams, match.HomeTeam)
		allTeams = append(allTeams, match.AwayTeam)
	}
	slices.SortFunc(allTeams, func(a, b string) int {
		return strings.Compare(NormalizeName(a), NormalizeName(b))
	})
	allTeams = lo.Uniq(allTeams)
	allTeamsStruct := lo.Map(allTeams, func(team string, _ int) teamResponse {
		return teamResponse{Name: team, ShortName: NormalizeName(team)}
	})
	return teamsResponse{AllTeams: allTeamsStruct}
}

func TeamsHandler(matches []Match) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, teamsService(matches))
	}
}

func TeamsHtmlHandler(matches []Match) func(c echo.Context) error {
	html := "<option value=\"%s\">%s</option>"
	return func(c echo.Context) error {
		result := teamsService(matches)
		htmlOptions := "<option value=\"\">Select Home Team</option>"
		for _, team := range result.AllTeams {
			htmlOptions += fmt.Sprintf(html, team.ShortName, team.Name)
		}
		return c.HTML(http.StatusOK, htmlOptions)
	}
}

type lastMatchesRequest struct {
	Team  string `query:"team"`
	Count int    `query:"count"`
	Where string `query:"where"`
}

// lastMatchesService returns the last `count` matches for the given team and location (home or away)
func lastMatchesService(matches []Match, req lastMatchesRequest) []Match {
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
		return a.MatchDate.Compare(b.MatchDate) * -1
	})
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
	return lo.Slice(matchesToCheck, 0, req.Count)
}

func LastMatchesHandler(matches []Match) func(c echo.Context) error {
	return func(c echo.Context) error {
		req := lastMatchesRequest{}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, lastMatchesService(matches, req))
	}
}
