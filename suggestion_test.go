/*

    gosh - a toy command shell
    Written in 2021 by Mirko Perillo
    To the extent possible under law, the author(s) have dedicated all copyright and related and neighboring rights to this software to the public domain worldwide. This software is distributed without any warranty.
    You should have received a copy of the CC0 Public Domain Dedication along with this software. If not, see <http://creativecommons.org/publicdomain/zero/1.0/>. 
*/

package main

import (
	"testing"
)

func TestFirst(t *testing.T) {
	var words = []string{"-l", "-a"}
	trie := initTrie(words)
	if len(trie.root.children) != 1 {
		t.Error("Expected one child")
	}
}

func TestAddFirstChild(t *testing.T) {
	var trie = initTrie([]string{"-l"})
	if len(trie.root.children) != 1 {
		t.Error()
	}
}

func TestTwoWords(t *testing.T) {
	var trie = initTrie([]string{"-l", "-a"})
	var rootRune = trie.root.c
	var firstLevelNode = trie.root.children['-']
	var secondLevelNodeL = trie.root.children['-'].children['l']
	var secondLevelNodeA = trie.root.children['-'].children['a']
	if rootRune != 0 {
		t.Error("Root node rune " + string(rootRune))
	}
	if firstLevelNode.c != '-' {
		t.Error("first level rune " + string(firstLevelNode.c))
	}

	if secondLevelNodeL.c != 'l' && secondLevelNodeL.endWord {
		t.Error("second level rune expected l " + string(secondLevelNodeL.c))
	}

	if secondLevelNodeA.c != 'a' && secondLevelNodeA.endWord {
		t.Error("second level rune expected a " + string(secondLevelNodeA.c))
	}
}

func TestLookupWord(t *testing.T) {
	var trie = initTrie([]string{"-l", "-a"})
	wordList := trie.lookup("-a")
	if len(wordList) > 1 || wordList[0] != "-a" {
		t.Error("Word not found: " + "-a")
	}
}

func TestSuggestions(t *testing.T) {
	var trie = initTrie([]string{"-l", "-list", "-lat", "-lng", "-a", "-bind"})
	wordList := trie.lookup("-l")
	expectedSuggestions := []string{"-l", "-list", "-lat", "-lng"}
	if len(wordList) != 4 || !assertSame(wordList, expectedSuggestions) {
		t.Error("Expected different words")
	}
}

func TestNoSuggestions(t *testing.T) {
	var trie = initTrie([]string{"-l", "-a", "-ls", "-ag"})
	wordList := trie.lookup("-an")
	if len(wordList) > 0 {
		t.Error("Expected no suggestions")
	}
}

func TestRichSuggestions(t *testing.T) {
	var trie = initTrie([]string{"-a", "--all", "-A", "--almost-all", "--author", "-b", "--escape", "--block-size", "-B", "--ignore-backups", "-c", "-C", "--color"})
	wordList := trie.lookup("--a")
	if len(wordList) != 3 && assertSame(wordList, []string{"--all", "--almost-all", "--author"}) {
		t.Error("Expected three suggestions")
	}
}

func TestMoreEndWordsOnSameLine(t *testing.T) {
	var trie = initTrie([]string{"-l", "-lin", "-line"})
	wordList := trie.lookup("-l")
	if !assertSame(wordList, []string{"-l", "-lin", "-line"}) {
		t.Error("Expected three suggestions")
	}
}

func assertSame(x []string, y []string) bool {
	var assertion = false
	for elemOfX := range x {
		for elemOfY := range y {
			if elemOfX == elemOfY {
				assertion = true
				break
			}
			assertion = false
		}
	}
	return assertion
}
