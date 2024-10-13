package internal

import (
	"fmt"
	"net/http"
	"reflect"
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

type resultMatrixRequest struct {
	MatchCountHome int `query:"match_count_home"`
	MatchCountAway int `query:"match_count_away"`
	HomeScored     int `query:"home_scored"`
	HomeConceded   int `query:"home_conceded"`
	AwayScored     int `query:"away_scored"`
	AwayConceded   int `query:"away_conceded"`
}

type ProbabilityWithOdds struct {
	Probability float64 `json:"probability"`
	Odds        float64 `json:"odds"`
}

type ResultMatrixResponse struct {
	HomeWin       ProbabilityWithOdds `json:"1"`
	Draw          ProbabilityWithOdds `json:"X"`
	AwayWin       ProbabilityWithOdds `json:"2"`
	HomeWinOrDraw ProbabilityWithOdds `json:"1X"`
	HomeWinOrAway ProbabilityWithOdds `json:"12"`
	AwayWinOrDraw ProbabilityWithOdds `json:"X2"`
	Over0_5Goals  ProbabilityWithOdds `json:"over_0.5"`
	Under0_5Goals ProbabilityWithOdds `json:"under_0.5"`
	Over1_5Goals  ProbabilityWithOdds `json:"over_1.5"`
	Under1_5Goals ProbabilityWithOdds `json:"under_1.5"`
	Over2_5Goals  ProbabilityWithOdds `json:"over_2.5"`
	Under2_5Goals ProbabilityWithOdds `json:"under_2.5"`
	Over3_5Goals  ProbabilityWithOdds `json:"over_3.5"`
	Under3_5Goals ProbabilityWithOdds `json:"under_3.5"`
	Over4_5Goals  ProbabilityWithOdds `json:"over_4.5"`
	Under4_5Goals ProbabilityWithOdds `json:"under_4.5"`
	Over5_5Goals  ProbabilityWithOdds `json:"over_5.5"`
	Under5_5Goals ProbabilityWithOdds `json:"under_5.5"`
	Over6_5Goals  ProbabilityWithOdds `json:"over_6.5"`
	Under6_5Goals ProbabilityWithOdds `json:"under_6.5"`
	Over7_5Goals  ProbabilityWithOdds `json:"over_7.5"`
	Under7_5Goals ProbabilityWithOdds `json:"under_7.5"`
	Goal          ProbabilityWithOdds `json:"goal"`
	NoGoal        ProbabilityWithOdds `json:"no_goal"`
	HomeGoal      ProbabilityWithOdds `json:"home_goal"`
	NoHomeGoal    ProbabilityWithOdds `json:"no_home_goal"`
	AwayGoal      ProbabilityWithOdds `json:"away_goal"`
	NoAwayGoal    ProbabilityWithOdds `json:"no_away_goal"`
	Result0_0     ProbabilityWithOdds `json:"0-0"`
	Result0_1     ProbabilityWithOdds `json:"0-1"`
	Result0_2     ProbabilityWithOdds `json:"0-2"`
	Result0_3     ProbabilityWithOdds `json:"0-3"`
	Result0_4     ProbabilityWithOdds `json:"0-4"`
	Result0_5     ProbabilityWithOdds `json:"0-5"`
	Result0_6     ProbabilityWithOdds `json:"0-6"`
	Result0_7     ProbabilityWithOdds `json:"0-7"`
	Result0_8     ProbabilityWithOdds `json:"0-8"`
	Result0_9     ProbabilityWithOdds `json:"0-9"`
	Result0_10    ProbabilityWithOdds `json:"0-10"`
	Result1_0     ProbabilityWithOdds `json:"1-0"`
	Result1_1     ProbabilityWithOdds `json:"1-1"`
	Result1_2     ProbabilityWithOdds `json:"1-2"`
	Result1_3     ProbabilityWithOdds `json:"1-3"`
	Result1_4     ProbabilityWithOdds `json:"1-4"`
	Result1_5     ProbabilityWithOdds `json:"1-5"`
	Result1_6     ProbabilityWithOdds `json:"1-6"`
	Result1_7     ProbabilityWithOdds `json:"1-7"`
	Result1_8     ProbabilityWithOdds `json:"1-8"`
	Result1_9     ProbabilityWithOdds `json:"1-9"`
	Result1_10    ProbabilityWithOdds `json:"1-10"`
	Result2_0     ProbabilityWithOdds `json:"2-0"`
	Result2_1     ProbabilityWithOdds `json:"2-1"`
	Result2_2     ProbabilityWithOdds `json:"2-2"`
	Result2_3     ProbabilityWithOdds `json:"2_3"`
	Result2_4     ProbabilityWithOdds `json:"2-4"`
	Result2_5     ProbabilityWithOdds `json:"2-5"`
	Result2_6     ProbabilityWithOdds `json:"2-6"`
	Result2_7     ProbabilityWithOdds `json:"2-7"`
	Result2_8     ProbabilityWithOdds `json:"2-8"`
	Result2_9     ProbabilityWithOdds `json:"2-9"`
	Result2_10    ProbabilityWithOdds `json:"2-10"`
	Result3_0     ProbabilityWithOdds `json:"3-0"`
	Result3_1     ProbabilityWithOdds `json:"3-1"`
	Result3_2     ProbabilityWithOdds `json:"3-2"`
	Result3_3     ProbabilityWithOdds `json:"3-3"`
	Result3_4     ProbabilityWithOdds `json:"3-4"`
	Result3_5     ProbabilityWithOdds `json:"3-5"`
	Result3_6     ProbabilityWithOdds `json:"3-6"`
	Result3_7     ProbabilityWithOdds `json:"3-7"`
	Result3_8     ProbabilityWithOdds `json:"3-8"`
	Result3_9     ProbabilityWithOdds `json:"3-9"`
	Result3_10    ProbabilityWithOdds `json:"3-10"`
	Result4_0     ProbabilityWithOdds `json:"4-0"`
	Result4_1     ProbabilityWithOdds `json:"4-1"`
	Result4_2     ProbabilityWithOdds `json:"4-2"`
	Result4_3     ProbabilityWithOdds `json:"4-3"`
	Result4_4     ProbabilityWithOdds `json:"4-4"`
	Result4_5     ProbabilityWithOdds `json:"4-5"`
	Result4_6     ProbabilityWithOdds `json:"4-6"`
	Result4_7     ProbabilityWithOdds `json:"4-7"`
	Result4_8     ProbabilityWithOdds `json:"4-8"`
	Result4_9     ProbabilityWithOdds `json:"4-9"`
	Result4_10    ProbabilityWithOdds `json:"4-10"`
	Result5_0     ProbabilityWithOdds `json:"5-0"`
	Result5_1     ProbabilityWithOdds `json:"5-1"`
	Result5_2     ProbabilityWithOdds `json:"5-2"`
	Result5_3     ProbabilityWithOdds `json:"5-3"`
	Result5_4     ProbabilityWithOdds `json:"5-4"`
	Result5_5     ProbabilityWithOdds `json:"5-5"`
	Result5_6     ProbabilityWithOdds `json:"5-6"`
	Result5_7     ProbabilityWithOdds `json:"5-7"`
	Result5_8     ProbabilityWithOdds `json:"5-8"`
	Result5_9     ProbabilityWithOdds `json:"5-9"`
	Result5_10    ProbabilityWithOdds `json:"5-10"`
	Result6_0     ProbabilityWithOdds `json:"6-0"`
	Result6_1     ProbabilityWithOdds `json:"6-1"`
	Result6_2     ProbabilityWithOdds `json:"6-2"`
	Result6_3     ProbabilityWithOdds `json:"6-3"`
	Result6_4     ProbabilityWithOdds `json:"6-4"`
	Result6_5     ProbabilityWithOdds `json:"6-5"`
	Result6_6     ProbabilityWithOdds `json:"6-6"`
	Result6_7     ProbabilityWithOdds `json:"6-7"`
	Result6_8     ProbabilityWithOdds `json:"6-8"`
	Result6_9     ProbabilityWithOdds `json:"6-9"`
	Result6_10    ProbabilityWithOdds `json:"6-10"`
	Result7_0     ProbabilityWithOdds `json:"7-0"`
	Result7_1     ProbabilityWithOdds `json:"7-1"`
	Result7_2     ProbabilityWithOdds `json:"7-2"`
	Result7_3     ProbabilityWithOdds `json:"7-3"`
	Result7_4     ProbabilityWithOdds `json:"7-4"`
	Result7_5     ProbabilityWithOdds `json:"7-5"`
	Result7_6     ProbabilityWithOdds `json:"7-6"`
	Result7_7     ProbabilityWithOdds `json:"7-7"`
	Result7_8     ProbabilityWithOdds `json:"7-8"`
	Result7_9     ProbabilityWithOdds `json:"7-9"`
	Result7_10    ProbabilityWithOdds `json:"7-10"`
	Result8_0     ProbabilityWithOdds `json:"8-0"`
	Result8_1     ProbabilityWithOdds `json:"8-1"`
	Result8_2     ProbabilityWithOdds `json:"8-2"`
	Result8_3     ProbabilityWithOdds `json:"8-3"`
	Result8_4     ProbabilityWithOdds `json:"8-4"`
	Result8_5     ProbabilityWithOdds `json:"8-5"`
	Result8_6     ProbabilityWithOdds `json:"8-6"`
	Result8_7     ProbabilityWithOdds `json:"8-7"`
	Result8_8     ProbabilityWithOdds `json:"8-8"`
	Result8_9     ProbabilityWithOdds `json:"8-9"`
	Result8_10    ProbabilityWithOdds `json:"8-10"`
	Result9_0     ProbabilityWithOdds `json:"9-0"`
	Result9_1     ProbabilityWithOdds `json:"9-1"`
	Result9_2     ProbabilityWithOdds `json:"9-2"`
	Result9_3     ProbabilityWithOdds `json:"9-3"`
	Result9_4     ProbabilityWithOdds `json:"9-4"`
	Result9_5     ProbabilityWithOdds `json:"9-5"`
	Result9_6     ProbabilityWithOdds `json:"9-6"`
	Result9_7     ProbabilityWithOdds `json:"9-7"`
	Result9_8     ProbabilityWithOdds `json:"9-8"`
	Result9_9     ProbabilityWithOdds `json:"9-9"`
	Result9_10    ProbabilityWithOdds `json:"9-10"`
	Result10_0    ProbabilityWithOdds `json:"10-0"`
	Result10_1    ProbabilityWithOdds `json:"10-1"`
	Result10_2    ProbabilityWithOdds `json:"10-2"`
	Result10_3    ProbabilityWithOdds `json:"10-3"`
	Result10_4    ProbabilityWithOdds `json:"10-4"`
	Result10_5    ProbabilityWithOdds `json:"10-5"`
	Result10_6    ProbabilityWithOdds `json:"10-6"`
	Result10_7    ProbabilityWithOdds `json:"10-7"`
	Result10_8    ProbabilityWithOdds `json:"10-8"`
	Result10_9    ProbabilityWithOdds `json:"10-9"`
	Result10_10   ProbabilityWithOdds `json:"10-10"`
}

func resultMatrixService(req resultMatrixRequest) map[string]ResultMatrixResponse {
	rm := NewResultMatrix(req.MatchCountHome, req.MatchCountAway, req.HomeScored, req.HomeConceded, req.AwayScored, req.AwayConceded)

	response := ResultMatrixResponse{
		HomeWin:       ProbabilityWithOdds{Probability: rm.GetHomeWinProbability(), Odds: AsOdds(rm.GetHomeWinProbability())},
		Draw:          ProbabilityWithOdds{Probability: rm.GetDrawProbability(), Odds: AsOdds(rm.GetDrawProbability())},
		AwayWin:       ProbabilityWithOdds{Probability: rm.GetAwayWinProbability(), Odds: AsOdds(rm.GetAwayWinProbability())},
		HomeWinOrDraw: ProbabilityWithOdds{Probability: rm.GetHomeWinOrDrawProbability(), Odds: AsOdds(rm.GetHomeWinOrDrawProbability())},
		HomeWinOrAway: ProbabilityWithOdds{Probability: rm.GetHomeWinOrAwayWinProbability(), Odds: AsOdds(rm.GetHomeWinOrAwayWinProbability())},
		AwayWinOrDraw: ProbabilityWithOdds{Probability: rm.GetAwayWinOrDrawProbability(), Odds: AsOdds(rm.GetAwayWinOrDrawProbability())},
		Over0_5Goals:  ProbabilityWithOdds{Probability: rm.GetOver0_5GoalsProbability(), Odds: AsOdds(rm.GetOver0_5GoalsProbability())},
		Under0_5Goals: ProbabilityWithOdds{Probability: rm.GetUnder0_5GoalsProbability(), Odds: AsOdds(rm.GetUnder0_5GoalsProbability())},
		Over1_5Goals:  ProbabilityWithOdds{Probability: rm.GetOver1_5GoalsProbability(), Odds: AsOdds(rm.GetOver1_5GoalsProbability())},
		Under1_5Goals: ProbabilityWithOdds{Probability: rm.GetUnder1_5GoalsProbability(), Odds: AsOdds(rm.GetUnder1_5GoalsProbability())},
		Over2_5Goals:  ProbabilityWithOdds{Probability: rm.GetOver2_5GoalsProbability(), Odds: AsOdds(rm.GetOver2_5GoalsProbability())},
		Under2_5Goals: ProbabilityWithOdds{Probability: rm.GetUnder2_5GoalsProbability(), Odds: AsOdds(rm.GetUnder2_5GoalsProbability())},
		Over3_5Goals:  ProbabilityWithOdds{Probability: rm.GetOver3_5GoalsProbability(), Odds: AsOdds(rm.GetOver3_5GoalsProbability())},
		Under3_5Goals: ProbabilityWithOdds{Probability: rm.GetUnder3_5GoalsProbability(), Odds: AsOdds(rm.GetUnder3_5GoalsProbability())},
		Over4_5Goals:  ProbabilityWithOdds{Probability: rm.GetOver4_5GoalsProbability(), Odds: AsOdds(rm.GetOver4_5GoalsProbability())},
		Under4_5Goals: ProbabilityWithOdds{Probability: rm.GetUnder4_5GoalsProbability(), Odds: AsOdds(rm.GetUnder4_5GoalsProbability())},
		Over5_5Goals:  ProbabilityWithOdds{Probability: rm.GetOver5_5GoalsProbability(), Odds: AsOdds(rm.GetOver5_5GoalsProbability())},
		Under5_5Goals: ProbabilityWithOdds{Probability: rm.GetUnder5_5GoalsProbability(), Odds: AsOdds(rm.GetUnder5_5GoalsProbability())},
		Over6_5Goals:  ProbabilityWithOdds{Probability: rm.GetOver6_5GoalsProbability(), Odds: AsOdds(rm.GetOver6_5GoalsProbability())},
		Under6_5Goals: ProbabilityWithOdds{Probability: rm.GetUnder6_5GoalsProbability(), Odds: AsOdds(rm.GetUnder6_5GoalsProbability())},
		Over7_5Goals:  ProbabilityWithOdds{Probability: rm.GetOver7_5GoalsProbability(), Odds: AsOdds(rm.GetOver7_5GoalsProbability())},
		Under7_5Goals: ProbabilityWithOdds{Probability: rm.GetUnder7_5GoalsProbability(), Odds: AsOdds(rm.GetUnder7_5GoalsProbability())},
		Goal:          ProbabilityWithOdds{Probability: rm.GetGoalProbability(), Odds: AsOdds(rm.GetGoalProbability())},
		NoGoal:        ProbabilityWithOdds{Probability: rm.GetNoGoalProbability(), Odds: AsOdds(rm.GetNoGoalProbability())},
		HomeGoal:      ProbabilityWithOdds{Probability: rm.GetHomeGoalProbability(), Odds: AsOdds(rm.GetHomeGoalProbability())},
		NoHomeGoal:    ProbabilityWithOdds{Probability: rm.GetNoHomeGoalProbability(), Odds: AsOdds(rm.GetNoHomeGoalProbability())},
		AwayGoal:      ProbabilityWithOdds{Probability: rm.GetAwayGoalProbability(), Odds: AsOdds(rm.GetAwayGoalProbability())},
		NoAwayGoal:    ProbabilityWithOdds{Probability: rm.GetNoAwayGoalProbability(), Odds: AsOdds(rm.GetNoAwayGoalProbability())},
	}

	// Add all results from 0-0 to 10-10
	for i := 0; i <= 10; i++ {
		for j := 0; j <= 10; j++ {
			prob := rm.GetResultProbability(i, j)
			odds := AsOdds(prob)
			field := fmt.Sprintf("Result%d_%d", i, j)
			reflect.ValueOf(&response).Elem().FieldByName(field).Set(reflect.ValueOf(ProbabilityWithOdds{Probability: prob, Odds: odds}))
		}
	}

	return map[string]ResultMatrixResponse{"result_matrix": response}
}

func ResultMatrixHandler(c echo.Context) error {
	req := resultMatrixRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, resultMatrixService(req))
}
