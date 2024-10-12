package internal_test

import (
	"math"
	"testing"

	"github.com/giorgiovilardo/tksgo/internal"
)

func TestResultMatrixProbabilities(t *testing.T) {
	rm := internal.NewResultMatrix(5, 5, 6, 5, 7, 5)

	testCases := []struct {
		name     string
		method   func() float64
		expected float64
	}{
		{"TotalProbability", rm.GetTotalProbability, 1.0},
		{"HomeWinProbability", rm.GetHomeWinProbability, 0.32045},
		{"DrawProbability", rm.GetDrawProbability, 0.30970},
		{"AwayWinProbability", rm.GetAwayWinProbability, 0.36983},
		{"HomeWinOrDrawProbability", rm.GetHomeWinOrDrawProbability, 0.63016},
		{"HomeWinOrAwayWinProbability", rm.GetHomeWinOrAwayWinProbability, 0.69029},
		{"AwayWinOrDrawProbability", rm.GetAwayWinOrDrawProbability, 0.67954},
		{"Over0_5GoalsProbability", rm.GetOver0_5GoalsProbability, 0.88650},
		{"Under0_5GoalsProbability", rm.GetUnder0_5GoalsProbability, 0.11349},
		{"Over1_5GoalsProbability", rm.GetOver1_5GoalsProbability, 0.68237},
		{"Under1_5GoalsProbability", rm.GetUnder1_5GoalsProbability, 0.31762},
		{"Over2_5GoalsProbability", rm.GetOver2_5GoalsProbability, 0.40396},
		{"Under2_5GoalsProbability", rm.GetUnder2_5GoalsProbability, 0.59603},
		{"Over3_5GoalsProbability", rm.GetOver3_5GoalsProbability, 0.20065},
		{"Under3_5GoalsProbability", rm.GetUnder3_5GoalsProbability, 0.79934},
		{"Over4_5GoalsProbability", rm.GetOver4_5GoalsProbability, 0.08375},
		{"Under4_5GoalsProbability", rm.GetUnder4_5GoalsProbability, 0.91624},
		{"Over5_5GoalsProbability", rm.GetOver5_5GoalsProbability, 0.02997},
		{"Under5_5GoalsProbability", rm.GetUnder5_5GoalsProbability, 0.97002},
		{"Over6_5GoalsProbability", rm.GetOver6_5GoalsProbability, 0.00936},
		{"Under6_5GoalsProbability", rm.GetUnder6_5GoalsProbability, 0.99063},
		{"Over7_5GoalsProbability", rm.GetOver7_5GoalsProbability, 0.00258},
		{"Under7_5GoalsProbability", rm.GetUnder7_5GoalsProbability, 0.99741},
		{"GoalProbability", rm.GetGoalProbability, 0.88650},
		{"NoGoalProbability", rm.GetNoGoalProbability, 0.11349},
		{"HomeGoalProbability", rm.GetHomeGoalProbability, 0.66712},
		{"NoHomeGoalProbability", rm.GetNoHomeGoalProbability, 0.33287},
		{"AwayGoalProbability", rm.GetAwayGoalProbability, 0.69880},
		{"NoAwayGoalProbability", rm.GetNoAwayGoalProbability, 0.30119},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.method()
			if math.Abs(result-tc.expected) > 1e-4 {
				t.Errorf("%s: expected %.4f, but got %.4f", tc.name, tc.expected, result)
			}
		})
	}
}

func TestGetResultProbability(t *testing.T) {
	rm := internal.NewResultMatrix(5, 5, 6, 5, 7, 5)

	testCases := []struct {
		homeResult int
		awayResult int
		expected   float64
	}{
		{0, 0, 0.11349},
		{1, 0, 0.09705},
		{0, 1, 0.10707},
		{1, 1, 0.14557},
		{2, 1, 0.07278},
		{1, 2, 0.07940},
		{3, 2, 0.01601},
		{2, 3, 0.01746},
	}

	for _, tc := range testCases {
		result := rm.GetResultProbability(tc.homeResult, tc.awayResult)
		if math.Abs(result-tc.expected) > 1e-4 {
			t.Errorf("GetResultProbability(%d, %d): expected %.4f, but got %.4f", tc.homeResult, tc.awayResult, tc.expected, result)
		}
	}
}
