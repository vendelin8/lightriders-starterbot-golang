package main

import (
	"math/rand"

	"github.com/vendelin8/lightriders-starterbot-golang/utils"
)

var (
	ownId, oppId int
	lastMove     utils.Direction
)

func botGetMove() string {
	moves := getAll(ownPosX, ownPosY, byte('.'))
	if len(moves) == 0 {
		lastMove = utils.Up
	} else {
		lastMove = moves[rand.Intn(len(moves))]
	}
	return lastMove.String()
}
