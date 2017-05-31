package main

import (
	"math/rand"

	"github.com/vendelin8/lightriders-starterbot-golang/utils"
)

var (
	ownBot, oppBot *utils.Player
)

func botInit() {
	ownBot = new(utils.Player)
	oppBot = new(utils.Player)
}

func botGetMove() string {
	moves := getAllMoves(ownBot, '.') //directions of empty fields next to the player
	if len(moves) == 0 {              //game over
		ownBot.SetMove(utils.Up)
	} else { //entry point: now a random is given
		ownBot.SetMove(moves[rand.Intn(len(moves))])
	}
	return ownBot.LastMove.String()
}
