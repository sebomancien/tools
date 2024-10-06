package zhe

import (
	"io"

	"gopkg.in/yaml.v3"
)

func ReadYAML(reader io.Reader) (*Config, error) {
	decoder := yaml.NewDecoder(reader)
	var data Config
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (res *Result) WriteYAML(writer io.Writer) error {
	encoder := yaml.NewEncoder(writer)
	return encoder.Encode(res)
}
