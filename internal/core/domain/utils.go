package domain

import (
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
		j.re = regexp.MustCompile(`((\\[|{)(.|\\s)*(}|\\]))`)
	})
}

func (j *jsonClearCompile) Clear(s string) string {
	j.init()
	var result = j.re.FindStringSubmatch(s)
	if result == nil {
		return s
	}

	return result[1]
}

func MustJSONClear(jsonString string) string {
	return jsonClearCompileInstance.Clear(jsonString)
}

func JSONClear(jsonString string) (string, error) {
	var result = jsonClearCompileInstance.Clear(jsonString)
	if result == "" {
		return "", ErrClearJson
	}
	return result, nil
}
