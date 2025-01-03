package zhe

import (
	"github.com/sebomancien/tools/pkg/expression"
	"github.com/sebomancien/tools/pkg/utils"
)

type variable struct {
	values []float64
}

type constraint struct {
	exp    expression.Operation
	target float64
	min    float64
	max    float64
	weight float64
}

type config struct {
	variables   []variable
	constraints []constraint
}

type solution struct {
	values []float64
	score  float64
}

type result struct {
	solutions  *utils.SortedList[solution]
	nbSolution uint64
}
