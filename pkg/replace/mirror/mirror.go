package mirror

import (
	"fmt"
	"regexp"
)

var rules = make([]*rule, 0)

type rule struct {
	matcher *regexp.Regexp
	replace string
}

type Config struct {
	Module  string
	Replace string
}

func Register(module, replace string) {
	regex, err := regexp.Compile(module)
	if err != nil {
		panic(fmt.Errorf("failed to parse module: %v", err))
	}

	rules = append(rules, &rule{
		matcher: regex,
		replace: replace,
	})
}

func Find(module string) string {
	for _, r := range rules {
		if r.matcher.MatchString(module) {
			return r.matcher.ReplaceAllString(module, r.replace)
		}
	}
	return module
}
