package wordstore

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"slowteetoe.com/autocompleter/Godeps/_workspace/src/github.com/derekparker/trie"
	"strings"
)

type Wordstore struct{}

var t *trie.Trie

func init() {
	t = trie.New()
	t.Add("cow", nil)
	t.Add("coffee", nil)
	t.Add("cake", nil)
	PopulateDictionary(t)
}

func (w Wordstore) PrefixSearch(searchTerm string) []string {
	results := []string{}

	if searchTerm == "" {
		log.Println("No search term provided.")
		return results
	}
	if r := t.PrefixSearch(searchTerm); r != nil {
		results = r
	}
	return results
}

func PopulateDictionary(t *trie.Trie) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Println(pwd)
	wordsfile, err := filepath.Abs("./data/words")
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadFile(wordsfile)
	if err != nil {
		panic(err)
	}
	words := strings.Split(string(b[:]), "\n")

	for _, element := range words {
		t.Add(element, nil)
	}
}
