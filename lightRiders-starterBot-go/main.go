package main

import (
	"bufio"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/vendelin8/lightriders-starterbot-golang/utils"
	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	opts         Options
	replayWriter *bufio.Writer
)

type Options struct {
	OmitReplay bool `short:"R" long:"omit-replay" description:"don't save replay"`
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		panic(err)
	}

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

func writeReplayInt(value int) {
	replayWriter.WriteRune(utils.REPLAY_INC + rune(value))
}

func writeReplayDirection(value utils.Direction) {
	replayWriter.WriteRune(utils.REPLAY_INC + rune(value))
}

func createReplayFile() {
	os.Mkdir(utils.REPLAY_DIR, 0755)
	fp, err := os.Create(path.Join(utils.REPLAY_DIR, time.Now().Format("20060102150405.txt")))
	if err != nil {
		panic(err)
	}
	replayWriter = bufio.NewWriter(fp)

	//first two runes are the width and height of the replay
	writeReplayInt(field.Width)
	writeReplayInt(field.Height)

	//then x, y positions of the first and second players
	writeReplayInt(oppBot.X)
	writeReplayInt(oppBot.Y)
	writeReplayInt(ownBot.X)
	writeReplayInt(ownBot.Y)
	replayWriter.Flush()
}
