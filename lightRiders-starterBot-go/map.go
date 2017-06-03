package main

import (
	"strconv"
	"strings"

	"github.com/vendelin8/lightriders-starterbot-golang/utils"
)

var (
	field *utils.Field
)

func atoi(in string) int {
	res, _ := strconv.Atoi(in)
	return res
}

func mapParseSetting(key, value string) {
	switch key {
	case "your_botid":
		ownBot.IdStr = value[0]
		ownBot.Id = atoi(value)
		oppBot.Id = (ownBot.Id + 1) % 2
		oppBot.IdStr = strconv.Itoa(oppBot.Id)[0]
		setSaveReplay(ownBot.Id != 1)
	case "field_width":
		field.Width = atoi(value)
	case "field_height":
		field.Height = atoi(value)
	}
}

func mapInit() {
	field = new(utils.Field)
}

func mapParse(in string) {
	if field.Rows == nil {
		field.Rows = make([][]byte, field.Height)
		k := 0
		for i := 0; i < field.Height; i++ {
			field.Rows[i] = make([]byte, field.Width)
			for j := 0; j < field.Width; j++ {
				ch := in[k]
				field.Rows[i][j] = ch
				k += 2
			}
		}
		MoveToMap(ownBot, in)
		MoveToMap(oppBot, in)
		createReplayFile()
	} else { //updating the map only
		//updating own position
		ownBot.MoveField(field)
		oppBot.SetAndMoveField(in, field)
		saveMovesToReplay()
	}
}

func MoveToMap(p *utils.Player, mapStr string) {
	pos := strings.Index(mapStr, string(p.IdStr)) / 2
	p.X = pos % field.Width
	p.Y = pos / field.Width
}

func getAllMoves(p *utils.Player, obj byte) []utils.Direction {
	//searching for all @obj next to the player
	result := make([]utils.Direction, 0)
	if p.X > 0 && field.Rows[p.Y][p.X-1] == obj {
		result = append(result, utils.Left)
	}
	if p.X+1 < field.Width && field.Rows[p.Y][p.X+1] == obj {
		result = append(result, utils.Right)
	}
	if p.Y > 0 && field.Rows[p.Y-1][p.X] == obj {
		result = append(result, utils.Up)
	}
	if p.Y+1 < field.Height && field.Rows[p.Y+1][p.X] == obj {
		result = append(result, utils.Down)
	}
	return result
}
