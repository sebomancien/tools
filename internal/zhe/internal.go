package zhe

import "github.com/sebomancien/tools/pkg/expression"

type variable struct {
	values []float32
}

type constraint struct {
	exp    expression.Operation
	target float32
	min    float32
	max    float32
	weight float32
}

type config struct {
	variables   []variable
	constraints []constraint
}

type solution struct {
	values []float32
	score  float32
}

type result struct {
	solutions  []solution
	nbSolution uint64
}
