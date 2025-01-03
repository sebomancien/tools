package zhe

import (
	"fmt"
	"strings"

	"github.com/sebomancien/tools/pkg/expression"
	"github.com/sebomancien/tools/pkg/utils"
)

func (in Variable) toInternal() variable {
	out := variable{
		values: make([]float64, len(in.Values)),
	}
	copy(out.values, in.Values)
	return out
}

func (in Constraint) toInternal() (constraint, error) {
	exp, err := expression.Parse(in.Formula)
	if err != nil {
		return constraint{}, err
	}

	return constraint{
		exp:    exp,
		target: in.Target,
		min:    in.Min,
		max:    in.Max,
		weight: in.Weight,
	}, nil
}

func (in *Config) toInternal() (*config, error) {
	out := config{
		variables:   make([]variable, len(in.Variables)),
		constraints: make([]constraint, len(in.Constraints)),
	}

	// Converts variables
	for i, v := range utils.IterateMapSorted(in.Variables) {
		// Replace the variable key by its index in each constraint
		for k2 := range in.Constraints {
			copy := in.Constraints[k2]
			copy.Formula = strings.ReplaceAll(copy.Formula, fmt.Sprintf("{%s}", i.Key), fmt.Sprintf("{%d}", i.Index))
			in.Constraints[k2] = copy
		}

		out.variables[i.Index] = v.toInternal()
	}

	// Converts constraints
	for i, c := range utils.IterateMapSorted(in.Constraints) {
		constraint, err := c.toInternal()
		if err != nil {
			return nil, err
		}

		out.constraints[i.Index] = constraint
	}

	return &out, nil
}

func (in solution) toExternal(config *Config) (Solution, error) {
	out := Solution{
		Variables:   make(map[string]string),
		Constraints: make(map[string]string),
		Score:       in.score,
	}

	if len(in.values) != (len(config.Variables) + len(config.Constraints)) {
		return Solution{}, fmt.Errorf("size mismatch")
	}

	// Converts variables
	for i, v := range utils.IterateMapSorted(config.Variables) {
		out.Variables[i.Key] = utils.FormatUnit(in.values[i.Index], v.Unit)
	}

	// Converts constraints
	for i, v := range utils.IterateMapSorted(config.Constraints) {
		out.Constraints[i.Key] = utils.FormatUnit(in.values[len(config.Variables)+i.Index], v.Unit)
	}

	return out, nil
}

func (in *result) toExternal(config *Config) (*Result, error) {
	out := Result{
		NbSolution: in.nbSolution,
		Solutions:  []Solution{},
	}

	for _, s := range in.solutions.All() {
		solution, err := s.toExternal(config)
		if err != nil {
			return nil, err
		}
		out.Solutions = append(out.Solutions, solution)
	}

	return &out, nil
}
