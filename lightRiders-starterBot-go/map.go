package main

import (
	"strconv"
	"strings"

	"github.com/vendelin8/lightriders-starterbot-golang/utils"
	//	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	mapWidth, mapHeight int
	mapLines            [][]byte
	ownPosX, ownPosY    int
	oppPosX, oppPosY    int
	ownIdStr, oppIdStr  byte
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
		ownPosX, ownPosY = lastMove.NewPos(ownPosX, ownPosY)
		mapLines[ownPosY][ownPosX] = ownIdStr

		mapLines[oppPosY][oppPosX] = byte('y') //for possible visualizing purposes
		oppPosX, oppPosY = whereIs(oppPosX, oppPosY, oppIdStr, in)
		mapLines[oppPosY][oppPosX] = oppIdStr
	}
	if ownId == 1 {
		replayWriter.WriteRune(utils.REPLAY_INC + rune(ownPosX))
		replayWriter.WriteRune(utils.REPLAY_INC + rune(ownPosY))
		replayWriter.WriteRune(utils.REPLAY_INC + rune(oppPosX))
		replayWriter.WriteRune(utils.REPLAY_INC + rune(oppPosY))
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

func getAll(x, y int, obj byte) []utils.Direction {
	//searching for all obj next to pos
	result := make([]utils.Direction, 0)
	if x > 0 && mapLines[y][x-1] == obj {
		result = append(result, utils.Left)
	}
	if x+1 < mapWidth && mapLines[y][x+1] == obj {
		result = append(result, utils.Right)
	}
	if y > 0 && mapLines[y-1][x] == obj {
		result = append(result, utils.Up)
	}
	if y+1 < mapHeight && mapLines[y+1][x] == obj {
		result = append(result, utils.Down)
	}
	return result
}
