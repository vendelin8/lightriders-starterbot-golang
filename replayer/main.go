//lightriders visualizer with termbox console drawing
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
	"github.com/vendelin8/lightriders-starterbot-golang/utils"
	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	opts        Options
	w, h        int //width and height of the termbox
	index       int //the current step
	bg, fg      termbox.Attribute
	mapW, mapH  int //width and height of the field
	replayReder *bufio.Reader
	ticker      *time.Ticker
	tickerC     <-chan time.Time
	top, left   int //left and top of the centered field
	p0, p1      *Player
	running     bool
)

type Options struct {
	File  string `short:"f" long:"file" description:"file to replay, last one will used if empty"`
	Delay int    `short:"d" long:"delay" description:"delay in millisec between turns" default:"600"`
}

type Player struct {
	utils.Player
	Moves []utils.Direction
}

func (p *Player) Fill() {
	//fills the player's position with it's id
	termbox.SetCell(left+p.X, top+p.Y, rune(p.IdStr), fg, bg)
}

func (p *Player) Move2() {
	//moves, fills and updates last move
	p.Move()
	p.Fill()
	p.LastMove = p.Moves[index]
}

func (p *Player) MoveFill() {
	//draws "x" to current position, player's id to the next one
	termbox.SetCell(left+p.X, top+p.Y, 'x', fg, bg)
	p.Move2()
}

func (p *Player) UnMove() {
	//undos one move
	termbox.SetCell(left+p.X, top+p.Y, '.', fg, bg)
	p.LastMove = p.Moves[index].Reverse()
	p.Move2()
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

	if len(opts.File) == 0 { //searching for last replay
		fileInfos, err := ioutil.ReadDir(utils.REPLAY_DIR)
		if err != nil {
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

	fp, err := os.Open(opts.File) //reading replay
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	replayReder = bufio.NewReader(fp)
	mapW = readReplayInt() //first 2 runes: width and height
	mapH = readReplayInt()

	err = termbox.Init() //termbox console drawing init
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	bg = termbox.ColorWhite
	fg = termbox.ColorBlack
	w, h = termbox.Size()

	//info to the screen
	printStr(0, 0, "quit            : Ctrl+C, Esc")
	printStr(0, 1, "back 1 turn     : Left")
	printStr(0, 2, "back 10 turns   : Down")
	printStr(0, 3, "forward 1 turn  : Right")
	printStr(0, 4, "forward 10 turns: Up")
	printStr(0, 5, "play/pause      : Space")
	printStr(0, 6, "-10% speed      : [")
	printStr(0, 7, "+10% speed      : ]")
	printStr(0, 8, "half speed      : {")
	printStr(0, 9, "double speed    : }")

	left = (w - mapW) / 2 //init map with dots
	top = (h - mapH) / 2
	for i := 0; i < mapH; i++ {
		for j := 0; j < mapW; j++ {
			termbox.SetCell(left+j, top+i, '.', fg, bg)
		}
	}

	p0 = new(Player) //init players
	p0.X = readReplayInt()
	p0.Y = readReplayInt()
	p0.IdStr = '0'
	p0.Moves = make([]utils.Direction, 0)
	p0.Fill()
	p1 = new(Player)
	p1.X = readReplayInt()
	p1.Y = readReplayInt()
	p1.IdStr = '1'
	p1.Moves = make([]utils.Direction, 0)
	p1.Fill()
	termbox.Flush()
	var d utils.Direction
	for { //load moves
		d = readReplayDirection()
		if d == utils.Wtf {
			break
		}
		p0.Moves = append(p0.Moves, d)
		p1.Moves = append(p1.Moves, readReplayDirection())
	}
	p0.LastMove = p0.Moves[0]
	p1.LastMove = p1.Moves[0]

	running = false
	index = 0

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

func readReplayInt() int {
	r, _, _ := replayReder.ReadRune()
	return int(r) - utils.REPLAY_INC
}

func readReplayDirection() utils.Direction {
	r, _, err := replayReder.ReadRune()
	if err != nil {
		return utils.Wtf
	}
	return utils.Direction(r - utils.REPLAY_INC)
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

func forward(turnCount int) bool {
	//go forward @turnCount turns, returns true if ended
	for ; turnCount > 0; turnCount-- {
		if index+1 >= len(p0.Moves) {
			p0.MoveFill()
			p1.MoveFill()
			printCenter(top+mapH, "Press any key to quit...")
			termbox.PollEvent()
			return true
		}
		index++
		p0.MoveFill()
		p1.MoveFill()
		termbox.Flush()
	}
	return false
}

func backward(turnCount int) {
	//go back @turnCount turns
	for ; turnCount > 0 && index > 0; turnCount-- {
		index--
		p0.UnMove()
		p1.UnMove()
		termbox.Flush()
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
