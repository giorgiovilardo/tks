package internal

import (
	"math"
)

type ResultMatrix struct {
	homeCoefficients []float64
	awayCoefficients []float64
	lambdaHome       float64
	lambdaAway       float64
}

func NewResultMatrix(matchCountHome, matchCountAway, homeScored, homeConceded, awayScored, awayConceded int) ResultMatrix {
	hc := make([]float64, 11)
	ac := make([]float64, 11)

	matchCountHomeFloat := float64(matchCountHome)
	matchCountAwayFloat := float64(matchCountAway)

	homeScoredAverage := float64(homeScored) / matchCountHomeFloat
	homeConcededAverage := float64(homeConceded) / matchCountHomeFloat
	awayScoredAverage := float64(awayScored) / matchCountAwayFloat
	awayConcededAverage := float64(awayConceded) / matchCountAwayFloat

	lambdaHome, lambdaAway := calcLambdas(homeScoredAverage, homeConcededAverage, awayScoredAverage, awayConcededAverage)

	for i := 0; i < 11; i++ {
		hc[i] = calcCoefficient(lambdaHome, i)
		ac[i] = calcCoefficient(lambdaAway, i)
	}

	return ResultMatrix{
		homeCoefficients: hc,
		awayCoefficients: ac,
		lambdaHome:       lambdaHome,
		lambdaAway:       lambdaAway,
	}
}

func (rm *ResultMatrix) GetTotalProbability() float64 {
	sum := 0.0
	for homeResult := 0; homeResult < 11; homeResult++ {
		for awayResult := 0; awayResult < 11; awayResult++ {
			sum += rm.GetResultProbability(homeResult, awayResult)
		}
	}
	return sum
}

func (rm *ResultMatrix) GetResultProbability(homeResult, awayResult int) float64 {
	var correctionFactor float64

	switch {
	case homeResult == 0 && awayResult == 0:
		correctionFactor = 1 + (rm.lambdaHome * rm.lambdaAway * 0.1)
	case homeResult == 1 && awayResult == 0:
		correctionFactor = 1 - (rm.lambdaAway * 0.1)
	case homeResult == 0 && awayResult == 1:
		correctionFactor = 1 - (rm.lambdaHome * 0.1)
	case homeResult == 1 && awayResult == 1:
		correctionFactor = 1.1
	default:
		correctionFactor = 1.0
	}

	return rm.homeCoefficients[homeResult] * rm.awayCoefficients[awayResult] * correctionFactor
}

func (rm *ResultMatrix) GetDrawProbability() float64 {
	sum := 0.0
	for i := 0; i < 11; i++ {
		sum += rm.GetResultProbability(i, i)
	}

	return sum
}

func (rm *ResultMatrix) GetHomeWinProbability() float64 {
	sum := 0.0
	for homeGoals := 1; homeGoals < 11; homeGoals++ {
		for awayGoals := 0; awayGoals < homeGoals; awayGoals++ {
			sum += rm.GetResultProbability(homeGoals, awayGoals)
		}
	}
	return sum
}

func (rm *ResultMatrix) GetAwayWinProbability() float64 {
	sum := 0.0
	for awayGoals := 1; awayGoals < 11; awayGoals++ {
		for homeGoals := 0; homeGoals < awayGoals; homeGoals++ {
			sum += rm.GetResultProbability(homeGoals, awayGoals)
		}
	}
	return sum
}

func (rm *ResultMatrix) GetHomeWinOrDrawProbability() float64 {
	return rm.GetHomeWinProbability() + rm.GetDrawProbability()
}

func (rm *ResultMatrix) GetAwayWinOrDrawProbability() float64 {
	return rm.GetAwayWinProbability() + rm.GetDrawProbability()
}

func (rm *ResultMatrix) GetHomeWinOrAwayWinProbability() float64 {
	return rm.GetHomeWinProbability() + rm.GetAwayWinProbability()
}

func (rm *ResultMatrix) GetOver0_5GoalsProbability() float64 {
	return rm.GetTotalProbability() - rm.GetResultProbability(0, 0)
}

func (rm *ResultMatrix) GetUnder0_5GoalsProbability() float64 {
	return rm.GetResultProbability(0, 0)
}

func (rm *ResultMatrix) GetOver1_5GoalsProbability() float64 {
	return rm.GetTotalProbability() - rm.GetResultProbability(0, 0) - rm.GetResultProbability(0, 1) - rm.GetResultProbability(1, 0)
}

func (rm *ResultMatrix) GetUnder1_5GoalsProbability() float64 {
	return rm.GetTotalProbability() - rm.GetOver1_5GoalsProbability()
}

func (rm *ResultMatrix) GetOver2_5GoalsProbability() float64 {
	return calcGenericOverXGoalsProbability(rm, 3)
}

func (rm *ResultMatrix) GetUnder2_5GoalsProbability() float64 {
	return rm.GetTotalProbability() - rm.GetOver2_5GoalsProbability()
}

func (rm *ResultMatrix) GetOver3_5GoalsProbability() float64 {
	return calcGenericOverXGoalsProbability(rm, 4)
}

func (rm *ResultMatrix) GetUnder3_5GoalsProbability() float64 {
	return rm.GetTotalProbability() - rm.GetOver3_5GoalsProbability()
}

func (rm *ResultMatrix) GetOver4_5GoalsProbability() float64 {
	return calcGenericOverXGoalsProbability(rm, 5)
}

func (rm *ResultMatrix) GetUnder4_5GoalsProbability() float64 {
	return rm.GetTotalProbability() - rm.GetOver4_5GoalsProbability()
}

func (rm *ResultMatrix) GetOver5_5GoalsProbability() float64 {
	return calcGenericOverXGoalsProbability(rm, 6)
}

func (rm *ResultMatrix) GetUnder5_5GoalsProbability() float64 {
	return rm.GetTotalProbability() - rm.GetOver5_5GoalsProbability()
}

func (rm *ResultMatrix) GetOver6_5GoalsProbability() float64 {
	return calcGenericOverXGoalsProbability(rm, 7)
}

func (rm *ResultMatrix) GetUnder6_5GoalsProbability() float64 {
	return rm.GetTotalProbability() - rm.GetOver6_5GoalsProbability()
}

func (rm *ResultMatrix) GetOver7_5GoalsProbability() float64 {
	return calcGenericOverXGoalsProbability(rm, 8)
}

func (rm *ResultMatrix) GetUnder7_5GoalsProbability() float64 {
	return rm.GetTotalProbability() - rm.GetOver7_5GoalsProbability()
}

func (rm *ResultMatrix) GetGoalProbability() float64 {
	return rm.GetTotalProbability() - rm.GetResultProbability(0, 0)
}

func (rm *ResultMatrix) GetNoGoalProbability() float64 {
	return rm.GetResultProbability(0, 0)
}

func (rm *ResultMatrix) GetHomeGoalProbability() float64 {
	sum := 0.0
	for homeGoals := 1; homeGoals < 11; homeGoals++ {
		for awayGoals := 0; awayGoals < 11; awayGoals++ {
			sum += rm.GetResultProbability(homeGoals, awayGoals)
		}
	}
	return sum
}

func (rm *ResultMatrix) GetNoHomeGoalProbability() float64 {
	return rm.GetTotalProbability() - rm.GetHomeGoalProbability()
}

func (rm *ResultMatrix) GetAwayGoalProbability() float64 {
	sum := 0.0
	for homeGoals := 0; homeGoals < 11; homeGoals++ {
		for awayGoals := 1; awayGoals < 11; awayGoals++ {
			sum += rm.GetResultProbability(homeGoals, awayGoals)
		}
	}
	return sum
}

func (rm *ResultMatrix) GetNoAwayGoalProbability() float64 {
	return rm.GetTotalProbability() - rm.GetAwayGoalProbability()
}

// calcGenericOverXGoalsProbability calculates the probability of the total number of goals being greater than x.
// use for bigger overs, lowers are easy to do by exclusion of precise cells
func calcGenericOverXGoalsProbability(rm *ResultMatrix, x int) float64 {
	sum := rm.GetTotalProbability()

	for homeGoals := 0; homeGoals < x; homeGoals++ {
		for awayGoals := 0; awayGoals < x-homeGoals; awayGoals++ {
			sum -= rm.GetResultProbability(homeGoals, awayGoals)
		}
	}
	return sum
}

func calcLambdas(homeScoredAverage, homeConcededAverage, awayScoredAverage, awayConcededAverage float64) (float64, float64) {
	lambdaHome := (homeScoredAverage + awayConcededAverage) / 2
	lambdaAway := (awayScoredAverage + homeConcededAverage) / 2
	return lambdaHome, lambdaAway
}

func calcCoefficient(lambda float64, i int) float64 {
	return (math.Pow(lambda, float64(i)) * math.Exp(-lambda)) / float64(fact(i))
}

func fact(n int) int {
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}
