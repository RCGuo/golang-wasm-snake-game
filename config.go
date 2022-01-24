package main

var STOP_KEY = 32
var MaxFPS = 120
var UPDATE_EVERY = (1000 / MaxFPS)
var MOVEMENT_KEYS = map[int]string{
	38: "TOP",
	87: "TOP",
	37: "LEFT",
	65: "LEFT",
	39: "RIGHT",
	68: "RIGHT",
	40: "DOWN",
	83: "DOWN",
}
var DIRECTION = map[string]Vector{
	"TOP":   {0, -1},
	"RIGHT": {1, 0},
	"DOWN":  {0, 1},
	"LEFT":  {-1, 0},
}
var DEFAULT_GAME_CONFIG = GameManager{
	Game: Game{
		Width:              17,
		Height:             17,
		Speed:              0.006,
		InitialSnakeLenght: 3,
		InitialDirection:   DIRECTION["RIGHT"],
		Score:              0,
	},
}
var SPEED_OPTIONS = map[string]float64{
	"0.5x": 0.5,
	"1x": 1,
	"1.5x": 1.5,
	"2x": 2,
	"2.5x": 2.5,
}