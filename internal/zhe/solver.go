package zhe

import (
	"log"
	"math"
	"slices"
	"sync/atomic"
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

func (res *result) insertSolution(solution solution, maxSolution int) {
	// Find the index of the solution, starting with the first one (worst one)
	i := 0
	for ; i < len(res.solutions); i++ {
		if solution.score >= res.solutions[i].score {
			break
		}
	}

	// If the solution buffer is full and this one is the worst, nothing to do
	if len(res.solutions) >= maxSolution && i == 0 {
		return
	}

	// Insert the new solution
	res.solutions = slices.Insert(res.solutions, i, solution)

	// Trim the number of solution by removing the first one (worst one)
	if len(res.solutions) > maxSolution {
		res.solutions = slices.Delete(res.solutions, 0, 1)
	}
}
