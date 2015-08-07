package wordstore

import (
	"slowteetoe.com/autocompleter/Godeps/_workspace/src/github.com/derekparker/trie"
)

type Wordstore struct{}

var t *trie.Trie

func init() {
	t = trie.New()
	t.Add("cow", nil)
	t.Add("coffee", nil)
	t.Add("cake", nil)
}

func (w Wordstore) PrefixSearch(searchTerm string) []string {
	results := []string{}
	if r := t.PrefixSearch(searchTerm); r != nil {
		results = r
	}
	return results
}
