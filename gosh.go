/*

    gosh - a toy command shell
    Written in 2021 by Mirko Perillo
    To the extent possible under law, the author(s) have dedicated all copyright and related and neighboring rights to this software to the public domain worldwide. This software is distributed without any warranty.
    You should have received a copy of the CC0 Public Domain Dedication along with this software. If not, see <http://creativecommons.org/publicdomain/zero/1.0/>. 
*/
package main

import (
	"encoding/csv"
	"fmt"
	"github.com/nsf/termbox-go"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

const configFolder = "./configs"

func printOnScreen(cursor *Cursor, msg string) {
	for _, c := range msg {
		if c != '\n' {
			termbox.SetCell(cursor.Col, cursor.Row, c, termbox.ColorDefault, termbox.ColorDefault)
			cursor.right()
		} else {
			cursor.newLine()
		}
	}
	termbox.Flush()
}

type Cursor struct {
	Row, Col int
	cmdLine  bool
}

func (c *Cursor) newLine() {
	c.Col = 0
	c.Row += 1
}

func (c *Cursor) right() {
	c.Col += 1
}

func (c *Cursor) left() {
	if c.Col > 0 {
		c.Col -= 1
	}
}

func (c *Cursor) deletePrevious() {
	var lastCol int
	if c.cmdLine {
		lastCol = 2
	} else {
		lastCol = 0
	}
	if c.Col > lastCol {
		c.Col--
		termbox.SetCell(c.Col, c.Row, rune(0), termbox.ColorDefault, termbox.ColorDefault)
		termbox.Flush()
	}
}

func (c *Cursor) cleanRow() {
	for c.Col > 0 {
		c.deletePrevious()
	}
}

func loadCompletions() map[string]*Trie {
	files, err := ioutil.ReadDir(configFolder)
	if err != nil {
		panic(err)
	}
	var suggestions = make(map[string]*Trie)
	for _, f := range files {
		fd, _ := os.Open(filepath.Join(configFolder, f.Name()))
		reader := csv.NewReader(fd)
		completions, _ := reader.Read()
		suggestions[f.Name()] = initTrie(completions)
	}

	return suggestions
}

func main() {
	cmdSuggestions := loadCompletions()
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	var inputBuffer = ""
	cursor := Cursor{cmdLine: true}
	printOnScreen(&cursor, "??")
	var suggestionCursor Cursor
	for {
		switch event := termbox.PollEvent(); event.Key {
		case termbox.KeyTab: //try to autocomplete command
			suggestionCursor.cleanRow()
			tokens := strings.Split(inputBuffer, " ")
			lastOption := tokens[len(tokens)-1]
			trie := cmdSuggestions[tokens[0]]
			suggestions := trie.lookup(lastOption)
			commandCursor := cursor
			cursor.newLine()
			cursor.cmdLine = false
			for _, suggestion := range suggestions {
				suggestion += " "
				printOnScreen(&cursor, suggestion)
			}
			suggestionCursor = cursor
			cursor = commandCursor
		case termbox.KeyEnter: //run the command
			suggestionCursor.cleanRow()
			inputBuffer = strings.TrimSpace(inputBuffer)
			command := inputBuffer
			if len(inputBuffer) > 0 {
				tokens := strings.Split(command, " ")
				isBuiltin := builtin(&cursor, tokens[0], tokens[1:]...)
				if !isBuiltin {
					out, err := exec.Command(tokens[0], tokens[1:]...).CombinedOutput()
					cursor.cmdLine = false
					if err != nil {
						cursor.newLine()
						printOnScreen(&cursor, " gosh: command not found: "+tokens[0])
					}
					cursor.newLine()
					printOnScreen(&cursor, string(out))
				}
			}
			inputBuffer = ""
			cursor.newLine()
			cursor.cmdLine = true
			printOnScreen(&cursor, "??")
		case termbox.KeyBackspace:
		case termbox.KeyBackspace2:
			cursor.cmdLine = true
			suggestionCursor.cleanRow()
			cursor.deletePrevious()
			if len(inputBuffer) > 0 {
				inputBuffer = inputBuffer[:len(inputBuffer)-1]
			}
		default:
			cursor.cmdLine = true
			var c = ""
			if event.Key == termbox.KeySpace {
				c = " "
			} else {
				c = string(event.Ch)
			}
			inputBuffer = inputBuffer + c
			printOnScreen(&cursor, c)
		}
	}
}

func builtin(cursor *Cursor, cmd string, args ...string) bool {
	switch cmd {
	case "bye":
		cursor.newLine()
		byeCommand()
		return true
	case "ciao":
		cursor.newLine()
		ciaoCommand(cursor)
		return true
	default:
		return false
	}
}

func byeCommand() {
	os.Exit(0)
}

func ciaoCommand(cursor *Cursor) {
	current, _ := user.Current()
	printOnScreen(cursor, fmt.Sprintf("ciao %v", current.Username))
}
