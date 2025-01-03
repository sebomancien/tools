package zhe

import (
	"log"
	"math"
	"sync/atomic"

	"github.com/sebomancien/tools/pkg/utils"
)

type subsolver struct {
	config  *config
	values  []float64
	channel chan<- solution
	count   uint64

	nbVariables   int
	nbConstraints int
}

func newSubsolver(config *config, channel chan<- solution, values ...float64) *subsolver {
	sub := subsolver{
		config:        config,
		values:        make([]float64, len(config.variables)+len(config.constraints)),
		nbVariables:   len(config.variables),
		nbConstraints: len(config.constraints),
		channel:       channel,
		count:         0,
	}
	copy(sub.values, values)
	return &sub
}

func newResult(capacity int) *result {
	return &result{
		nbSolution: 0,
		solutions: utils.NewSortedList[solution](capacity, func(a, b solution) int {
			switch {
			case a.score < b.score:
				return 1
			case a.score > b.score:
				return -1
			default:
				return 0
			}
		}),
	}
}

func (s *subsolver) solve(depth int) {
	if depth == s.nbVariables-1 {
		for _, s.values[depth] = range s.config.variables[depth].values {
			valid := s.evaluate()
			if valid {
				solution := solution{
					values: make([]float64, s.nbVariables+s.nbConstraints),
				}
				copy(solution.values, s.values)
				s.channel <- solution
			}
		}
		atomic.AddUint64(&s.count, uint64(len(s.config.variables[depth].values)))
	} else {
		for _, s.values[depth] = range s.config.variables[depth].values {
			s.solve(depth + 1)
		}
	}
}

func (s *subsolver) evaluate() bool {
	// Evaluate all constraints and check their validity
	for i, c := range s.config.constraints {
		value, err := c.exp.Evaluate(s.values)
		if err != nil {
			log.Fatal(err)
		}
		if value < c.min || value > c.max {
			return false
		}
		s.values[s.nbVariables+i] = value
	}

	return true
}

func (s *solution) computeScore(config *config) {
	// Compute the solution score
	s.score = 0
	for i, c := range config.constraints {
		diff := math.Abs(s.values[len(config.variables)+i] - c.target)
		if c.target != 0 {
			diff /= c.target
		}
		s.score += diff * c.weight
	}
}
