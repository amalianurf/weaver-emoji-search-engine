package main

import (
	"context"
	"sort"
	"strings"

	"github.com/ServiceWeaver/weaver"
	"golang.org/x/exp/slices"
)

type Searcher interface {
	Search(context.Context, string) ([]string, error)
}

type searcher struct {
	weaver.Implements[Searcher]
}

func (s *searcher) Search(ctx context.Context, query string) ([]string, error) {
	s.Logger(ctx).Debug("Search", "query", query)
	words := strings.Fields(strings.ToLower(query))

	var result []string
	for emoji, labels := range emojis {
		if matches(labels, words) {
			result = append(result, emoji)
		}
	}

	sort.Strings(result)
	return result, nil
}

func matches(labels, words []string) bool {
	for _, word := range words {
		if !slices.Contains(labels, word) {
			return false
		}
	}
	return true
}
