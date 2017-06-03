package main

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	rand.Seed(time.Now().Unix())

	defer catchRuntimeErrors()

	debugInit()
	botInit()
	mapInit()
	var text string
	var textParts []string
	running := true
	for running {
		text, _ = reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if len(text) == 0 {
			continue
		}
		textParts = strings.Split(text, " ")
		switch textParts[0] {
		case "action":
			if textParts[1] == "move" {
				writer.WriteString(botGetMove())
				writer.WriteRune('\n')
				writer.Flush()
			}
		case "update":
			if textParts[1] == "game" && textParts[2] == "field" {
				mapParse(textParts[3])
			}
		case "settings":
			mapParseSetting(textParts[1], textParts[2])
		case "quit", "end":
			running = false
		}
	}
}
