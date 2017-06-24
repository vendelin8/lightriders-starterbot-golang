//lightriders visualizer with termbox console drawing
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/nsf/termbox-go"
	"github.com/vendelin8/lightriders-starterbot-golang/utils"
	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	opts         Options
	w, h         int //width and height of the termbox
	index        int //the current step
	bg, fg       termbox.Attribute
	replayReder  *bufio.Reader
	ticker       *time.Ticker
	tickerC      <-chan time.Time
	top, left    int //left and top of the centered field
	p0, p1       *Player
	running      bool
	rf           utils.ReplayFormat
	maxDebugLens []int
	debugValues  [][]string
	debugMiddle  int
)

type Options struct {
	File   string `short:"f" long:"file" description:"file to replay, last one will used if empty"`
	Delay  int    `short:"d" long:"delay" description:"delay in millisec between turns" default:"600"`
	InvCol bool   `short:"i" long:"invert-colors" description:"invert colors for dark terminal backgrounds"`
}

type Player struct {
	utils.Player
	Moves []utils.Direction
	Fg    termbox.Attribute
}

func (p *Player) Fill() {
	//fills the player's position with it's id
	termbox.SetCell(left+p.X, top+p.Y, rune(p.IdStr), p.Fg, bg)
}

func (p *Player) Move2() {
	//moves, fills and updates last move
	p.Move()
	p.Fill()
	p.LastMove = p.Moves[index]
}

func (p *Player) MoveFill(indexDiff int) {
	//draws "x" to current position, player's id to the next one
	var toDraw rune
	if index > 1 {
		lastLast := p.Moves[index+indexDiff-2]
		switch p.LastMove {
		case utils.Up:
			if lastLast == utils.Right {
				toDraw = '┘'
			} else if lastLast == utils.Left {
				toDraw = '└'
			} else {
				toDraw = '│'
			}
		case utils.Left:
			if lastLast == utils.Up {
				toDraw = '┐'
			} else if lastLast == utils.Down {
				toDraw = '┘'
			} else {
				toDraw = '─'
			}
		case utils.Down:
			if lastLast == utils.Right {
				toDraw = '┐'
			} else if lastLast == utils.Left {
				toDraw = '┌'
			} else {
				toDraw = '│'
			}
		case utils.Right:
			if lastLast == utils.Up {
				toDraw = '┌'
			} else if lastLast == utils.Down {
				toDraw = '└'
			} else {
				toDraw = '─'
			}
		}
	} else {
		toDraw = 'x'
	}
	termbox.SetCell(left+p.X, top+p.Y, toDraw, p.Fg, bg)
	p.Move2()
}

func (p *Player) UnMove() {
	//undos one move
	termbox.SetCell(left+p.X, top+p.Y, '.', fg, bg)
	p.LastMove = p.Moves[index].Reverse()
	p.Move2()
}

func main() {
	var err error
	if _, err = flags.ParseArgs(&opts, os.Args); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
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

	if len(opts.File) == 0 { //searching for last replay
		var fileInfos []os.FileInfo
		if fileInfos, err = ioutil.ReadDir(utils.REPLAY_DIR); err != nil {
			panic(err)
		}
		latest := time.Now().AddDate(-100, 0, 0)
		for _, info := range fileInfos { //loading languages
			name := path.Join(utils.REPLAY_DIR, info.Name())
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

	var b []byte //parsing replay file
	if b, err = ioutil.ReadFile(opts.File); err != nil {
		panic(err)
	}
	bs := bytes.Split(b, []byte{utils.REPLAY_SEPARATOR})
	if err = json.Unmarshal(bs[0], &rf); err != nil {
		panic(err)
	}

	if err = termbox.Init(); err != nil { //termbox console drawing init
		panic(err)
	}
	defer termbox.Close()
	bg = termbox.ColorDefault
	if opts.InvCol {
		fg = termbox.ColorWhite
	} else {
		fg = termbox.ColorBlack
	}
	w, h = termbox.Size()

	//key info to the screen
	keys := []string{"quit", "Ctrl+C, Esc",
		"back 1 turn", "Left",
		"back 10 turns", "Down",
		"forward 1 turn", "Right",
		"forward 10 turns", "Up",
		"play/pause", "Space",
		"-10% speed", "[",
		"+10% speed", "]",
		"half speed", "{",
		"double speed", "}"}
	var maxLeft, maxRight, i, j int
	for i = 0; i < len(keys); i += 2 { //getting max lengths for padding
		if len(keys[i]) > maxLeft {
			maxLeft = len(keys[i])
		}
		if len(keys[i+1]) > maxRight {
			maxRight = len(keys[i+1])
		}
	}
	left = w - maxLeft - maxRight - 2
	middle := w - maxRight
	for i = 0; i < len(keys); i += 2 {
		printStr(left, i/2, keys[i])
		termbox.SetCell(middle-2, i/2, ':', fg, bg)
		printStr(middle, i/2, keys[i+1])
	}

	var key string
	for _, key = range utils.DEBUG_VARS { //getting padding for debug vars
		if len(key) > debugMiddle {
			debugMiddle = len(key)
		}
	}
	for i, key = range utils.DEBUG_VARS { //printing debug vars' label
		printStr(i, 0, key)
		termbox.SetCell(debugMiddle, i, ':', fg, bg)
		printStr(middle, i/2, keys[i+1])
	}
	debugMiddle += 2

	rounds := len(bs) - 1
	p0 = new(Player) //init players
	p0.X = rf.OppX
	p0.Y = rf.OppY
	p0.IdStr = '0'
	p0.Moves = make([]utils.Direction, rounds)
	p0.Fg = termbox.ColorGreen
	p1 = new(Player)
	p1.X = rf.OwnX
	p1.Y = rf.OwnY
	p1.IdStr = '1'
	p1.Moves = make([]utils.Direction, rounds)
	p1.Fg = termbox.ColorRed
	var rm utils.ReplayMove
	maxDebugLens = make([]int, len(utils.DEBUG_VARS))
	debugValues = make([][]string, rounds)
	for i = 0; i < rounds; i++ { //load moves
		if err = json.Unmarshal(bs[i+1], &rm); err != nil {
			panic(err)
		}
		p0.Moves[i] = rm.OppMove
		p1.Moves[i] = rm.OwnMove
		debugValues[i] = make([]string, len(utils.DEBUG_VARS))
		for j, key = range rm.Others {
			if len(key) > maxDebugLens[j] {
				maxDebugLens[j] = len(key)
			}
			debugValues[i][j] = key
		}
	}
	p0.LastMove = p0.Moves[0]
	p1.LastMove = p1.Moves[0]

	left = 0 //init map with dots to the middle
	for _, i = range maxDebugLens {
		if i > left {
			left = i
		}
	}
	left = left + debugMiddle + (w-rf.FieldWidth-left-maxLeft-maxRight-debugMiddle)/2
	top = (h - rf.FieldHeight) / 2
	for i := 0; i < rf.FieldHeight; i++ {
		for j := 0; j < rf.FieldWidth; j++ {
			termbox.SetCell(left+j, top+i, '.', fg, bg)
		}
	}

	running = false
	index = 0
	p0.Fill()
	p1.Fill()
	writeDebugValues()

	defer func() { //init main loop
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

	startCounter() //start main loop
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

func printStr(x, y int, text string) {
	for i, value := range text {
		termbox.SetCell(x+i, y, rune(value), fg, bg)
	}
}

func printCenter(y int, text string) {
	printStr((w-len(text))/2, y, text)
	termbox.Flush()
}

func writeDebugValues() {
	for i, key := range debugValues[index] {
		printStr(debugMiddle, i, key)
		for j := len(key); j < maxDebugLens[i]; j++ {
			termbox.SetCell(debugMiddle+j, i, ' ', fg, bg)
		}
	}
	termbox.Flush()
}

func forward(turnCount int) bool {
	//go forward @turnCount turns, returns true if ended
	for ; turnCount > 0; turnCount-- {
		if index+1 >= len(p0.Moves) {
			p0.MoveFill(1)
			p1.MoveFill(1)
			printCenter(top+rf.FieldHeight+1, "Press any key to quit...")
			termbox.PollEvent()
			return true
		}
		index++
		p0.MoveFill(0)
		p1.MoveFill(0)
		writeDebugValues()
	}
	return false
}

func backward(turnCount int) {
	//go back @turnCount turns
	for ; turnCount > 0 && index > 0; turnCount-- {
		index--
		p0.UnMove()
		p1.UnMove()
		writeDebugValues()
	}
}

func startCounter() {
	ticker = time.NewTicker(time.Millisecond * time.Duration(opts.Delay))
	tickerC = ticker.C
	running = true
}
func stopCounter() {
	ticker.Stop()
	running = false
}

func changeSpeed(delayMult float32) {
	stopCounter()
	opts.Delay = int(float32(opts.Delay) * delayMult)
	startCounter()
}
