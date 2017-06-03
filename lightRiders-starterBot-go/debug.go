package main

import (
	"bufio"
	"os"
	"path"
	"time"

	"github.com/vendelin8/lightriders-starterbot-golang/utils"
	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	saveReplay   bool
	replayWriter *bufio.Writer
)

func debugInit() {
	log.New() //initializing logger
	handler := log.Must.FileHandler("log.txt", log.LogfmtFormat())
	log.Root().SetHandler(log.MultiHandler(
		log.LvlFilterHandler(log.LvlWarn, log.CallerStackHandler("%+v", handler)),
		log.MatchFilterHandler("lvl", log.LvlInfo, handler)))

}

func logI(msg string, ctx ...interface{}) {
	log.Info(msg, ctx)
}

func logE(msg string, ctx ...interface{}) {
	log.Error(msg, ctx)
}

func catchRuntimeErrors() { //startup error check
	if r := recover(); r != nil {
		logE("main", "error", r)
	}
}

func setSaveReplay(in bool) {
	saveReplay = true
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

func saveMovesToReplay() {
	writeReplayDirection(oppBot.LastMove)
	writeReplayDirection(ownBot.LastMove)
	replayWriter.Flush()
}
