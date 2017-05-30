package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/nsf/termbox-go"
	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	opts Options
)

const (
	replayDir  = "replays"
	REPLAY_INC = 32
)

type Options struct {
	File  string `short:"f" long:"file" description:"file to replay, last one will used if empty"`
	Delay int    `short:"d" long:"delay" description:"delay in millisec between turns" default:"1000"`
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	log.New() //initializing logger
	handler := log.Must.FileHandler("log.txt", log.LogfmtFormat())
	log.Root().SetHandler(log.MultiHandler(
		log.LvlFilterHandler(log.LvlWarn, log.CallerStackHandler("%+v", handler)),
		log.MatchFilterHandler("lvl", log.LvlInfo, handler)))

	defer func() { //startup error check
		if r := recover(); r != nil {
			log.Error("main", "error", r)
			panic(r)
		}
	}()

	if len(opts.File) == 0 {
		fileInfos, err := ioutil.ReadDir(replayDir)
		if err != nil {
			panic(err)
		}
		latest := time.Now().AddDate(-100, 0, 0)
		for _, info := range fileInfos { //loading languages
			name := path.Join(replayDir, info.Name())
			s, _ := os.Stat(name)
			if s.IsDir() {
				continue
			}
			if s.ModTime().After(latest) {
				opts.File = name
				latest = s.ModTime()
			}
		}
	}

	fp, err := os.Open(opts.File)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	scanner.Scan()
	t := scanner.Text()
	mapWidth := int(t[0])
	mapHeight := int(t[1])

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	bg := termbox.ColorWhite
	fg := termbox.ColorBlack
	w, h := termbox.Size()

	printStr := func(x, y int, text string, flush bool) {
		for i, value := range text {
			termbox.SetCell(x+i, y, rune(value), fg, bg)
		}
		if flush {
			termbox.Flush()
		}
	}

	printStr(0, 0, "quit            : Ctrl+C, Esc", false)
	printStr(0, 1, "back 1 turn     : Left", false)
	printStr(0, 2, "back 10 turns   : Down", false)
	printStr(0, 3, "forward 1 turn  : Right", false)
	printStr(0, 4, "forward 10 turns: Up", false)
	printStr(0, 5, "play/pause      : Space", false)
	printStr(0, 6, "-10% speed      : [", false)
	printStr(0, 7, "+10% speed      : ]", false)
	printStr(0, 8, "half speed      : {", false)
	printStr(0, 9, "double speed    : }", false)
	w -= 29

	left := (w - mapWidth) / 2
	top := (h - mapHeight) / 2
	for i := 0; i < mapHeight; i++ {
		for j := 0; j < mapWidth; j++ {
			termbox.SetCell(left+j, top+i, '.', fg, bg)
		}
	}

	index := 0
	lines := make([][]int, 0)
	var line []int
	for scanner.Scan() {
		t = scanner.Text()
		linesTmp := make([]int, 4)
		linesTmp[0] = left + int(t[0]) - REPLAY_INC
		linesTmp[1] = top + int(t[1]) - REPLAY_INC
		linesTmp[2] = left + int(t[2]) - REPLAY_INC
		linesTmp[3] = top + int(t[3]) - REPLAY_INC
		lines = append(lines, linesTmp)
	}

	if err = scanner.Err(); err != nil {
		panic(err)
	}

	updateSingle := func(a, b rune) {
		termbox.SetCell(line[0], line[1], a, fg, bg)
		termbox.SetCell(line[2], line[3], b, fg, bg)
	}

	line = lines[index]
	updateSingle('1', '0')
	termbox.Flush()

	printCenter := func(y int, text string) {
		printStr((w-len(text))/2, y, text, true)
	}

	update := func(diff int) {
		if diff == 1 {
			updateSingle('x', 'y')
		} else {
			updateSingle('.', '.')
		}
		index += diff
		line = lines[index]
		updateSingle('1', '0')
		termbox.Flush()
	}

	forward := func(turnCount int) bool {
		for ; turnCount > 0; turnCount-- {
			if index+1 >= len(lines) {
				printCenter(top+mapHeight, "Press any key to quit...")
				termbox.PollEvent()
				return true
			}
			update(1)
		}
		return false
	}

	backward := func(turnCount int) {
		for ; turnCount > 0 && index > 0; turnCount-- {
			update(-1)
		}
	}

	var ticker *time.Ticker
	var tickerC <-chan time.Time
	running := false

	startCounter := func() {
		ticker = time.NewTicker(time.Millisecond * time.Duration(opts.Delay))
		tickerC = ticker.C
		running = true
	}
	stopCounter := func() {
		ticker.Stop()
		running = false
	}

	defer func() {
		if running {
			stopCounter()
		}
	}()
	eventQueue := make(chan termbox.Event)

	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	changeSpeed := func(delayMult float32) {
		stopCounter()
		opts.Delay = int(float32(opts.Delay) * delayMult)
		startCounter()
	}

	startCounter()
	for {
		select {
		case <-tickerC:
			if forward(1) {
				return
			}
		case e := <-eventQueue:
			if e.Type == termbox.EventKey {
				switch e.Key {
				case termbox.KeyCtrlC, termbox.KeyEsc:
					return
				case termbox.KeyArrowDown: //-10 turns
					backward(10)
				case termbox.KeyArrowUp: //+10 turns
					if forward(10) {
						return
					}
				case termbox.KeyArrowLeft: //-1 turn
					backward(1)
				case termbox.KeyArrowRight: //+1 turn
					if forward(1) {
						return
					}
				case termbox.KeySpace:
					if running {
						stopCounter()
					} else {
						startCounter()
					}
				default:
					switch e.Ch {
					case 91: //[, -10%
						changeSpeed(1.1)
					case 93: //], +10%
						changeSpeed(0.9)
					case 123: //{, half speed
						changeSpeed(2)
					case 125: //}, double speed
						changeSpeed(0.5)
					}
				}
			}
		}
	}
}
