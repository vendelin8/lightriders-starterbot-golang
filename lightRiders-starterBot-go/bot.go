package main

import "math/rand"

var (
	ownId, oppId int
	lastMove     Direction
)

func botGetMove() string {
	moves := getAll(ownPosX, ownPosY, byte('.'))
	if len(moves) == 0 {
		lastMove = Up
	} else {
		lastMove = moves[rand.Intn(len(moves))]
	}
	return lastMove.String()
}
