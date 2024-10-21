package utils

import (
	"fmt"
)

type unitScale struct {
	value  float64
	prefix string
}

func FormatUnit(value float64, unit string) string {
	var scales = []unitScale{
		{value: 1e+15, prefix: "P"},
		{value: 1e+12, prefix: "T"},
		{value: 1e+09, prefix: "G"},
		{value: 1e+06, prefix: "M"},
		{value: 1e+03, prefix: "k"},
		{value: 1e+00, prefix: ""},
		{value: 1e-03, prefix: "m"},
		{value: 1e-06, prefix: "µ"},
		{value: 1e-09, prefix: "n"},
		{value: 1e-12, prefix: "p"},
	}

	for _, s := range scales {
		if value >= s.value {
			tmp := fmt.Sprintf("%g%s%s", value/s.value, s.prefix, unit)
			return tmp
		}
	}

	return fmt.Sprintf("%g%s", value, unit)
}
