package types

import (
	"github.com/mattn/go-shellwords"
)

func ParseEntrypoint(target string) ([]*string, error) {
	tokens, err := shellwords.Parse(target)
	if err != nil {
		return []*string{}, err
	}

	result := []*string{}
	for _, token := range tokens {
		s := token
		result = append(result, &s)
	}
	return result, nil
}
