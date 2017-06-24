package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/vendelin8/lightriders-starterbot-golang/utils"
	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	saveReplay   bool
	replayWriter *bufio.Writer
	round        int
)

func debugInit() {
	log.New() //initializing logger
	handler := log.Must.FileHandler("log.txt", log.LogfmtFormat())
	log.Root().SetHandler(log.MultiHandler(
		log.LvlFilterHandler(log.LvlWarn, log.CallerStackHandler("%+v", handler)),
		log.MatchFilterHandler("lvl", log.LvlInfo, handler)))
}

func logI(msg string, ctx ...interface{}) {
	log.Info(msg, ctx...)
}

func logE(msg string, ctx ...interface{}) {
	log.Error(msg, ctx...)
}

func catchRuntimeErrors() { //startup error check
	if r := recover(); r != nil {
		logE("main", "error", r)
	}
}

func setSaveReplay(in bool) {
	saveReplay = true
}

func createReplayFile() {
	os.Mkdir(utils.REPLAY_DIR, 0755)
	fp, err := os.Create(path.Join(utils.REPLAY_DIR, time.Now().Format("20060102150405.txt")))
	if err != nil {
		panic(err)
	}
	replayWriter = bufio.NewWriter(fp)

	rf := utils.ReplayFormat{field.Width, field.Height, ownBot.X, ownBot.Y, oppBot.X, oppBot.Y}
	b, err := json.Marshal(rf)
	if err != nil {
		panic(err)
	}
	replayWriter.Write(b)
	replayWriter.Flush()
}

func saveMovesToReplay() {
	replayWriter.WriteRune(utils.REPLAY_SEPARATOR)
	rm := utils.ReplayMove{ownBot.LastMove, oppBot.LastMove,
		[]string{strconv.Itoa(round)}} //example for debugging variables
	b, err := json.Marshal(rm)
	if err != nil {
		panic(err)
	}
	replayWriter.Write(b)
	replayWriter.Flush()
	round++
}
