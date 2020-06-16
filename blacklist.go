package messagehook

import (
	"fmt"
	"regexp"
)

type Blacklist struct {
	regexes []*regexp.Regexp
}

func NewBlacklist(patterns []string) *Blacklist {
	b := &Blacklist{}
	for _, entry := range patterns {
		r, err := regexp.Compile(entry)
		if err != nil {
			fmt.Println(err)
			continue
		}
		b.regexes = append(b.regexes, r)
	}

	return b
}

func (b *Blacklist) Match(text string) bool {
	for _, r := range b.regexes {
		if r.MatchString(text) {
			return true
		}
	}
	return false
}
