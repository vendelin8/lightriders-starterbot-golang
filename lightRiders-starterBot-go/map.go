package main

import (
	"strconv"
	"strings"
	//	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	mapWidth, mapHeight int
	mapLines            [][]byte
	ownPosX, ownPosY    int
	oppPosX, oppPosY    int
	ownIdStr, oppIdStr  byte
)

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
	Wtf
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "up"
	case Right:
		return "right"
	case Down:
		return "down"
	case Left:
		return "left"
	}
	return "wtf"
}

const (
	REPLAY_INC = 32
)

func atoi(in string) int {
	res, _ := strconv.Atoi(in)
	return res
}

func mapParseSetting(key, value string) {
	switch key {
	case "your_botid":
		ownIdStr = byte(value[0])
		ownId = atoi(value)
		oppId = (ownId + 1) % 2
		oppIdStr = byte(strconv.Itoa(oppId)[0])
	case "field_width":
		mapWidth = atoi(value)
	case "field_height":
		mapHeight = atoi(value)
	}
}

func mapParse(in string) {
	if mapLines == nil {
		if ownId == 1 {
			createReplayFile()
		}
		mapLines = make([][]byte, mapHeight)
		k := 0
		for i := 0; i < mapHeight; i++ {
			mapLines[i] = make([]byte, mapWidth)
			for j := 0; j < mapWidth; j++ {
				ch := byte(in[k])
				mapLines[i][j] = ch
				k += 2
			}
		}
		ownPos := strings.Index(in, string(ownIdStr)) / 2
		ownPosX = ownPos % mapWidth
		ownPosY = ownPos / mapWidth
		oppPos := strings.Index(in, string(oppIdStr)) / 2
		oppPosX = oppPos % mapWidth
		oppPosY = oppPos / mapWidth
	} else { //updating the map only
		//updating own position
		mapLines[ownPosY][ownPosX] = byte('x') //for possible visualizing purposes
		switch lastMove {
		case Up:
			ownPosY -= 1
		case Down:
			ownPosY += 1
		case Left:
			ownPosX -= 1
		case Right:
			ownPosX += 1
		}
		mapLines[ownPosY][ownPosX] = ownIdStr

		mapLines[oppPosY][oppPosX] = byte('y') //for possible visualizing purposes
		oppPosX, oppPosY = whereIs(oppPosX, oppPosY, oppIdStr, in)
		mapLines[oppPosY][oppPosX] = oppIdStr
	}
	if ownId == 1 {
		replayWriter.WriteRune(REPLAY_INC + rune(ownPosX))
		replayWriter.WriteRune(REPLAY_INC + rune(ownPosY))
		replayWriter.WriteRune(REPLAY_INC + rune(oppPosX))
		replayWriter.WriteRune(REPLAY_INC + rune(oppPosY))
		replayWriter.WriteRune('\n')
		replayWriter.Flush()
	}
}

func whereIs(x, y int, obj byte, in string) (int, int) {
	//searching for obj next to pos
	if x > 0 && in[(y*mapWidth+x-1)*2] == obj {
		return x - 1, y
	}
	if x+1 < mapWidth && in[(y*mapWidth+x+1)*2] == obj {
		return x + 1, y
	}
	if y > 0 && in[((y-1)*mapWidth+x)*2] == obj {
		return x, y - 1
	}
	return x, y + 1
}

func getAll(x, y int, obj byte) []Direction {
	//searching for all obj next to pos
	result := make([]Direction, 0)
	if x > 0 && mapLines[y][x-1] == obj {
		result = append(result, Left)
	}
	if x+1 < mapWidth && mapLines[y][x+1] == obj {
		result = append(result, Right)
	}
	if y > 0 && mapLines[y-1][x] == obj {
		result = append(result, Up)
	}
	if y+1 < mapHeight && mapLines[y+1][x] == obj {
		result = append(result, Down)
	}
	return result
}
