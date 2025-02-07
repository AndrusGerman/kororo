package domain

import (
	"fmt"
	"regexp"
	"sync"
)

type jsonClearCompile struct {
	sync.Once
	re *regexp.Regexp
}

var jsonClearCompileInstance = new(jsonClearCompile)

func (j *jsonClearCompile) init() {
	j.Do(func() {
		j.re = regexp.MustCompile(`({|\[)(\w|\W)*(}|\])`)
	})
}

func (j *jsonClearCompile) Clear(s string) (string, error) {
	j.init()
	var result = j.re.FindStringSubmatch(s)
	if result == nil {
		return "", fmt.Errorf("%w: %s not valid json", ErrClearJson, s)
	}
	return result[0], nil

}

func MustJSONClear(jsonString string) string {
	result, err := jsonClearCompileInstance.Clear(jsonString)
	if err != nil {
		return ""
	}
	return result
}

func JSONClear(jsonString string) (string, error) {
	return jsonClearCompileInstance.Clear(jsonString)
}
