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
	cache weaver.Ref[Cache]
}

func (s *searcher) Search(ctx context.Context, query string) ([]string, error) {
	s.Logger(ctx).Debug("Search", "query", query)

	if emojis, err := s.cache.Get().Get(ctx, query); err != nil {
		s.Logger(ctx).Error("cache.Get", "query", query, "err", err)
	} else if len(emojis) > 0 {
		return emojis, nil
	}

	words := strings.Fields(strings.ToLower(query))

	var result []string
	for emoji, labels := range emojis {
		if matches(labels, words) {
			result = append(result, emoji)
		}
	}

	sort.Strings(result)

	if err := s.cache.Get().Put(ctx, query, result); err != nil {
		s.Logger(ctx).Error("cache.Put", "query", query, "err", err)
	}

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
