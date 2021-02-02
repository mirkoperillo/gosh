/*

    gosh - a toy command shell
    Written in 2021 by Mirko Perillo
    To the extent possible under law, the author(s) have dedicated all copyright and related and neighboring rights to this software to the public domain worldwide. This software is distributed without any warranty.
    You should have received a copy of the CC0 Public Domain Dedication along with this software. If not, see <http://creativecommons.org/publicdomain/zero/1.0/>. 
*/

package main

type Trie struct {
	root *Node
}

type Node struct {
	c        rune
	endWord  bool
	children map[rune]Node
}

func (t *Trie) add(word string) {
	n := t.root
	for pos, r := range word {
		if child, exists := n.children[r]; !exists {
			newNode := Node{r, pos == len(word)-1, make(map[rune]Node)}
			n.children[r] = newNode
			n = &newNode
		} else {
			n = &child
		}
	}
}

func initTrie(flags []string) *Trie {
	var trie = Trie{}
	trie.root = &Node{children: make(map[rune]Node)}
	for _, word := range flags {
		trie.add(word)
	}
	return &trie
}

func (node *Node) print(level int) string {
	var treeAsString = ""
	var indent = "  "
	var identation = ""
	for i := 0; i < level; i++ {
		identation = identation + indent
	}
	isEnd := " 0 "
	if node.endWord {
		isEnd = " 1 "
	}
	treeAsString += identation + string(node.c) + isEnd + "\n"
	for _, child := range node.children {
		treeAsString += child.print(level + 1)
	}
	return treeAsString
}

func (node *Node) explore(prefix string, suggestions []string) []string {
	if len(node.children) > 0 {
		for c, child := range node.children {
			suggestions = child.explore(prefix+string(c), suggestions)
		}
	} else {
		suggestions = append(suggestions, prefix)
	}
	return suggestions
}

func (t *Trie) lookup(prefix string) []string {
	wordList := []string{}
	node := t.root
	for pos, c := range prefix {
		if child, exists := node.children[c]; exists {
			if pos == len(prefix)-1 {
				if child.endWord {
					wordList = append(wordList, prefix)
				}
				if len(child.children) > 0 {
					wordList = child.explore(prefix, wordList)
				}
			}
			node = &child
		} else {
			break
		}
	}
	return wordList
}
