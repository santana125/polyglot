package repository

import (
	"fmt"
	"sort"

	"github.com/filhodanuvem/polyglot/language"
)

var excludeList = map[string]bool{
	"":                          true,
	"Text":                      true,
	"Markdown":                  true,
	"Ignore List":               true,
	"JSON":                      true,
	"Git Config":                true,
	"XML Property List counter": true,
	"Gradle":                    true,
	"XML":                       true,
}

type Statistics struct {
	counters   []Counter
	langs      map[string]int
	reposCount int
}

type Counter struct {
	Lang    string `json:"language"`
	Counter int    `json:"counter"`
}

// Implementing Sort interface
func (s *Statistics) Len() int {
	return len(s.counters)
}

func (s *Statistics) Less(i, j int) bool {
	return s.counters[i].Counter > s.counters[j].Counter
}

func (s *Statistics) Swap(i, j int) {
	s.counters[i], s.counters[j] = s.counters[j], s.counters[i]
}

func (s *Statistics) String() string {
	return fmt.Sprintf("%+v", s.counters)
}

func (s *Statistics) Length() int {
	return s.reposCount
}

func (s *Statistics) FirstLanguages(length int) []Counter {
	sort.Sort(s)

	if length > len(s.counters) {
		length = len(s.counters)
	}
	return s.counters[0:length]
}

func (s *Statistics) Merge(stats *Statistics) {
	s.reposCount++
	if s.langs == nil {
		s.langs = make(map[string]int)
	}
	for i := range stats.counters {
		lang := stats.counters[i].Lang
		if _, ok := s.langs[lang]; !ok {
			s.langs[lang] = len(s.counters)
			s.counters = append(s.counters, stats.counters[i])
			continue
		}
		s.counters[s.langs[lang]].Counter += stats.counters[i].Counter
	}
}

func GetStatistics(files []string) (Statistics, error) {
	var stats Statistics
	stats.langs = make(map[string]int)
	for i := range files {
		lang, err := language.DetectByFile(files[i])
		if err != nil {
			return stats, err
		}
		if _, ok := excludeList[lang]; ok {
			continue
		}
		if _, ok := stats.langs[lang]; !ok {
			stats.langs[lang] = len(stats.counters)
			c := Counter{
				Lang:    lang,
				Counter: 0,
			}
			stats.counters = append(stats.counters, c)
			continue
		}
		stats.counters[stats.langs[lang]].Counter++
	}

	stats.reposCount = 1
	return stats, nil
}
