package main

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	replayWriter *bufio.Writer
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	rand.Seed(time.Now().Unix())

	log.New() //initializing logger
	handler := log.Must.FileHandler("log.txt", log.LogfmtFormat())
	log.Root().SetHandler(log.MultiHandler(
		log.LvlFilterHandler(log.LvlWarn, log.CallerStackHandler("%+v", handler)),
		log.MatchFilterHandler("lvl", log.LvlInfo, handler)))

	defer func() { //startup error check
		if r := recover(); r != nil {
			log.Error("main", "error", r)
		}
	}()

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

func createReplayFile() {
	fp, err := os.Create("replays/" + time.Now().Format("20060102150405.txt"))
	if err != nil {
		panic(err)
	}
	replayWriter = bufio.NewWriter(fp)
	replayWriter.WriteRune(rune(mapWidth))
	replayWriter.WriteRune(rune(mapHeight))
	replayWriter.WriteRune('\n')
}
