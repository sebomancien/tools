package zhe

import (
	"log"
	"sync"
	"sync/atomic"
)

type Variable struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Values      []float64 `yaml:"values"`
	Unit        string    `yaml:"unit"`
}

type Constraint struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Formula     string  `yaml:"formula"`
	Target      float64 `yaml:"target"`
	Unit        string  `yaml:"unit"`
	Min         float64 `yaml:"min"`
	Max         float64 `yaml:"max"`
	Weight      float64 `yaml:"weight"`
}

type Config struct {
	Variables   map[string]Variable   `yaml:"variables"`
	Constraints map[string]Constraint `yaml:"constraints"`
}

type Solution struct {
	Variables   map[string]string `yaml:"variables"`
	Constraints map[string]string `yaml:"constraints"`
	Score       float64           `yaml:"score"`
}

type Result struct {
	NbSolution uint64     `yaml:"nb_solutions"`
	Solutions  []Solution `yaml:"solutions"`
}

type Progress struct {
	Counter uint64
	Total   uint64
}

type Solver struct {
	config          *Config
	subsolvers      []*subsolver
	nbPossibilities uint64
}

func NewSolver(config *Config) *Solver {
	return &Solver{
		config:          config,
		subsolvers:      nil,
		nbPossibilities: 0,
	}
}

func (s *Solver) Solve(maxSolution int) (*Result, error) {
	if len(s.config.Variables) == 0 || len(s.config.Constraints) == 0 {
		log.Fatal("No variables or constraints")
	}

	// Compute the total number of possibilities
	s.nbPossibilities = 1
	for _, variables := range s.config.Variables {
		s.nbPossibilities *= uint64(len(variables.Values))
	}

	// Convert to an internal config
	config, err := s.config.toInternal()
	if err != nil {
		return nil, err
	}

	// Create a channel to share solutions
	channel := make(chan solution, 100)

	// Start concurrent solver goroutines
	var wgSolver sync.WaitGroup
	for _, value := range config.variables[0].values {
		wgSolver.Add(1)
		subsolver := newSubsolver(config, channel, value)
		go func() {
			subsolver.solve(1)
			wgSolver.Done()
		}()
		s.subsolvers = append(s.subsolvers, subsolver)
	}

	result := newResult(maxSolution)

	// Start a goroutine to compile all solutions from the solver goroutines
	var wgResult sync.WaitGroup
	wgResult.Add(1)
	go func() {
		for solution := range channel {
			solution.computeScore(config)
			result.solutions.Insert(solution)
			result.nbSolution++
		}
		wgResult.Done()
	}()

	// Wait for all solver goroutines to complete
	wgSolver.Wait()

	// Close the channel and wait for the result routine to complete
	close(channel)
	wgResult.Wait()

	// Convert to an external result
	return result.toExternal(s.config)
}

func (s *Solver) GetProgress() Progress {
	var count uint64 = 0
	for _, subsolver := range s.subsolvers {
		count += atomic.LoadUint64(&subsolver.count)
	}
	return Progress{
		Counter: count,
		Total:   s.nbPossibilities,
	}
}
