package main

import (
	"log"
	"math"
	"math/rand"
	"strconv"
	"syscall/js"
	"time"
)

var doc js.Value
var window js.Value

func init() {
	rand.Seed(time.Now().UnixMilli())
	window = js.Global()
	doc = window.Get("document")
}

type ContainerSize struct {
	Width  float64
	Height float64
}

type Game struct {
	Width              float64
	Height             float64
	Speed              float64
	InitialSnakeLenght float64
	InitialDirection   Vector
	Snake              []Vector
	Direction          Vector
	Food               Vector
	Movement           string
	Score              int
	LastUpdate         int64
	StopTime           int64
	BestScore          int
}

type Projector struct {
	ProjectDistance func(float64) float64
	ProjectPosition func(Vector) Vector
}

type GameManager struct {
	ContainerSize
	Game
	Projector
}

func getFood(width, height float64, snake []Vector) Vector {
	allPositions := make([]Vector, 0)
	for _, wv := range GetRange(width) {
		for _, hv := range GetRange(height) {
			x, y := wv, hv
			allPositions = append(allPositions, Vector{
				x: float64(x) + 0.5,
				y: float64(y) + 0.5,
			})
		}
	}
	segments := GetSegmentsFromVectors(snake)

	freePositions := make([]Vector, 0)
	for _, point := range allPositions {
		isNotInside := 0
		for _, segment := range segments {
			if !segment.IsPointInside(point) {
				isNotInside = 1
			}
		}
		if isNotInside == 1 {
			freePositions = append(freePositions, point)
		}
	}
	return GetRandomFrom(freePositions)
}

func getBestScore() int {
	best := window.Get("localStorage").Call("getItem", "bestScore")
	besti, err := strconv.Atoi(best.String())
	if err != nil {
		log.Println("fail to conver best score to interger")
	}
	return besti
}

func getGameInitialMgr(mgr *GameManager) *GameManager {
	mgr.Game = DEFAULT_GAME_CONFIG.Game
	head := Vector{
		x: math.Round(mgr.Game.Width/2) - 0.5,
		y: math.Round(mgr.Game.Height/2) - 0.5,
	}
	mgr.Game.Snake = []Vector{
		head.Subtract(mgr.Game.InitialDirection.ScaleBy(mgr.Game.InitialSnakeLenght)),
		head,
	}
	mgr.Game.Food = getFood(mgr.Game.Width, mgr.Game.Height, mgr.Game.Snake)
	mgr.Game.Direction = mgr.Game.InitialDirection
	mgr.Game.Score = 0
	mgr.Game.BestScore = getBestScore()
	return mgr
}

type Accumulator struct {
	distance float64
	tail     []Vector
}

func reduce(list []Vector, initialValue Accumulator, fn func(acc Accumulator, point Vector, index int) Accumulator) Accumulator {
	acc := initialValue
	for i, v := range list {
		acc = fn(acc, v, i)
	}
	return acc
}

func getNewTail(oldSnake []Vector, distance float64) []Vector {
	initialValue := Accumulator{
		distance: distance,
		tail:     []Vector{},
	}
	newTail := reduce(GetWithoutLastElement(oldSnake), initialValue, func(acc Accumulator, point Vector, index int) Accumulator {
		if len(acc.tail) != 0 {
			var tail []Vector
			tail = append(tail, acc.tail...)
			tail = append(tail, point)
			acc.tail = tail
			return acc
		}
		next := oldSnake[index+1]
		segment := Segment{
			start: point,
			end:   next,
		}
		length := segment.Length()
		if length >= initialValue.distance {
			vector := segment.GetVector().Normalize().ScaleBy(acc.distance)
			var tail []Vector
			tail = append(tail, acc.tail...)
			tail = append(tail, point.Add(vector))
			return Accumulator{
				distance: 0,
				tail:     tail,
			}
		} else {
			acc.distance = acc.distance - length
			return acc
		}
	}).tail

	return newTail
}

func getNewDirection(oldDirection Vector, movement string) Vector {
	var shouldChange bool
	var newDirection Vector
	if _, ok := DIRECTION[movement]; ok {
		newDirection = DIRECTION[movement]
		if !oldDirection.IsOpposite(newDirection) {
			shouldChange = true
		}
	} else {
		shouldChange = false
	}

	if shouldChange {
		return newDirection
	} else {
		return oldDirection
	}
}

func getMgrAfterMoveProcessing(mgr *GameManager, movement string, distance float64) *GameManager {
	newTail := getNewTail(mgr.Game.Snake, distance)
	oldHead := GetLastElement(mgr.Game.Snake)
	newHead := oldHead.Add(mgr.Game.Direction.ScaleBy(distance))
	newDirection := getNewDirection(mgr.Game.Direction, movement)

	if !mgr.Game.Direction.EqualTo(newDirection) {
		oldX, oldY := oldHead.x, oldHead.y
		oldXRounded, oldYRounded, newXRounded, newYRounded :=
			math.Round(oldX), math.Round(oldY), math.Round(newHead.x), math.Round(newHead.y)
		getMgrWithBrokenSnake := func(old, oldRounded, newRounded float64, getBreakpoint func(point float64) Vector) *GameManager {
			var roundedPoint float64
			if newRounded > oldRounded {
				roundedPoint = 0.5
			} else {
				roundedPoint = -0.5
			}
			breakpointComponent := oldRounded + roundedPoint
			breakpoint := getBreakpoint(breakpointComponent)
			vector := newDirection.ScaleBy(distance - math.Abs(old-breakpointComponent))
			head := breakpoint.Add(vector)

			mgr.Game.Direction = newDirection
			newTail = append(newTail, breakpoint, head)
			mgr.Game.Snake = newTail
			return mgr
		}

		if oldXRounded != newXRounded {
			return getMgrWithBrokenSnake(oldX, oldXRounded, newXRounded, func(x float64) Vector {
				return Vector{x, oldY}
			})
		}
		if oldYRounded != newYRounded {
			return getMgrWithBrokenSnake(oldY, oldYRounded, newYRounded, func(y float64) Vector {
				return Vector{oldX, y}
			})
		}
	}
	newMgr := mgr
	newMgr.Game.Snake = append(newTail, newHead)
	return newMgr
}

func getMgrAfterFoodProcessing(mgr *GameManager) *GameManager {
	headSegment := Segment{
		GetLastElement(GetWithoutLastElement(mgr.Game.Snake)),
		GetLastElement(mgr.Game.Snake),
	}

	if !headSegment.IsPointInside(mgr.Game.Food) {
		return mgr
	}

	var tailEnd, beforeTailEnd Vector
	var restOfSnake []Vector
	tailEnd, beforeTailEnd, restOfSnake =
		mgr.Game.Snake[0], mgr.Game.Snake[1], mgr.Game.Snake[2:]
	tailSegment := Segment{beforeTailEnd, tailEnd}
	newTailEnd := tailEnd.Add(tailSegment.GetVector().Normalize())

	var snake []Vector
	snake = append(snake, newTailEnd, beforeTailEnd)
	snake = append(snake, restOfSnake...)
	mgr.Game.Snake = snake
	mgr.Game.Score = mgr.Game.Score + 1
	mgr.Game.Food = getFood(mgr.Game.Width, mgr.Game.Height, snake)
	return mgr
}

func isGameOver(snake []Vector, width float64, height float64) bool {
	head := GetLastElement(snake)
	if head.x < 0 || head.x > width || head.y < 0 || head.y > height {
		return true
	}

	if len(snake) < 5 {
		return false
	}

	head, tail := snake[len(snake)-1], snake[:len(snake)-2]
	segments := GetSegmentsFromVectors(tail)
	for _, segment := range segments[2:] {
		projected := segment.GetProjectedPoint(head)
		if !segment.IsPointInside(projected) {
			return false
		} else {
			distance := Segment{head, projected}.Length()
			return distance < 0.5
		}
	}
	return false
}

func getNewGameManager(mgr *GameManager, movement string, timespan int64) *GameManager {
	distance := mgr.Game.Speed * float64(timespan)
	mgrAfterMove := getMgrAfterMoveProcessing(mgr, movement, distance)
	mgrAfterFood := getMgrAfterFoodProcessing(mgrAfterMove)
	if isGameOver(mgrAfterFood.Snake, mgrAfterFood.Game.Width, mgrAfterFood.Game.Width) {
		return getGameInitialMgr(mgr)
	}
	return mgrAfterFood
}

func getProjectors(containerSize ContainerSize, gameSize Game) Projector {
	widthRatio := containerSize.Width / gameSize.Width
	heightRatio := containerSize.Height / gameSize.Height
	unitOnScreen := math.Min(widthRatio, heightRatio)

	return Projector{
		ProjectDistance: func(distance float64) float64 {
			return distance * unitOnScreen
		},
		ProjectPosition: func(position Vector) Vector {
			return position.ScaleBy(unitOnScreen)
		},
	}
}

func getContainer() js.Value {
	return doc.Call("getElementById", "container")
}

func getContainerSize() ContainerSize {
	clientRect := getContainer().Call("getBoundingClientRect")
	width, height := clientRect.Get("width").Float(), clientRect.Get("height").Float()
	return ContainerSize{
		Width:  width,
		Height: height,
	}
}

func clearContainer() {
	container := getContainer()
	for container.Call("hasChildNodes").Equal(js.ValueOf(true)) {
		child := container.Get("lastElementChild")
		container.Call("removeChild", child)
	}
}

func getContext(width, height float64) js.Value {
	canvas := doc.Call("getElementsByTagName", "canvas")
	if canvas.Get("length").Int() == 0 {
		canvas = doc.Call("createElement", "canvas")
		container := getContainer()
		container.Call("appendChild", canvas)
	}
	canvas = doc.Call("getElementsByTagName", "canvas").Get("0")
	canvas.Set("width", width)
	canvas.Set("height", height)
	context := canvas.Call("getContext", "2d")
	context.Call("clearRect", 0, 0, canvas.Get("width").Float(), canvas.Get("height").Float())
	return context
}

func renderCells(context js.Value, cellSide, width, height float64) {
	context.Set("globalAlpha", 0.2)
	for _, columnNum := range GetRange(width) {
		for _, rowNum := range GetRange(height) {
			if (columnNum+rowNum)%2 == 1 {
				context.Call("fillRect", (float64(columnNum) * cellSide), (float64(rowNum) * cellSide), cellSide, cellSide)
			}
		}
	}
	context.Set("globalAlpha", 1)
}

func renderFood(context js.Value, cellSide float64, point Vector) {
	context.Call("beginPath")
	context.Call("arc", point.x, point.y, cellSide/2.5, 0, 2*math.Pi)
	context.Set("fillStyle", "#e74c3c")
	context.Call("fill")
}

func renderSnake(context js.Value, cellSide float64, snake []Vector) {
	context.Set("lineWidth", cellSide)
	context.Set("strokeStyle", "#3498db")
	context.Call("beginPath")
	for _, v := range snake {
		context.Call("lineTo", v.x, v.y)
	}
	context.Call("stroke")
}

func renderScores(score, bestScore int) {
	doc.Call("getElementById", "current-score").Set("innerText", score)
	doc.Call("getElementById", "best-score").Set("innerText", bestScore)
}

func render(mgr *GameManager) {
	viewWidth, viewHeight :=
		mgr.Projector.ProjectDistance(mgr.Game.Width),
		mgr.Projector.ProjectDistance(mgr.Game.Height)
	context := getContext(viewWidth, viewHeight)
	cellSide := viewWidth / mgr.Game.Width
	renderCells(context, cellSide, mgr.Game.Width, mgr.Game.Height)
	renderFood(context, cellSide, mgr.Projector.ProjectPosition(mgr.Game.Food))
	projectedSnake := make([]Vector, len(mgr.Game.Snake))
	for i, v := range mgr.Game.Snake {
		projectedSnake[i] = mgr.Projector.ProjectPosition(v)
	}
	renderSnake(context, cellSide, projectedSnake)
	renderScores(mgr.Game.Score, mgr.Game.BestScore)
}

func getInitialMgr() *GameManager {
	containerSize := getContainerSize()
	mgr := getGameInitialMgr(&GameManager{})
	mgr.ContainerSize = containerSize
	mgr.Projector = getProjectors(containerSize, mgr.Game)
	mgr.Game.BestScore = getBestScore()
	return mgr
}

func getNewMgrPropsOnTick(oldMgr *GameManager) *GameManager {
	if oldMgr.Game.StopTime != 0 {
		return oldMgr
	}

	lastUpdate := time.Now().UnixMilli()
	if oldMgr.Game.LastUpdate != 0 {
		mgr := getNewGameManager(oldMgr, oldMgr.Game.Movement, lastUpdate-oldMgr.Game.LastUpdate)
		mgr.Game.LastUpdate = lastUpdate
		if mgr.Game.Score > oldMgr.Game.BestScore {
			window.Get("localStorage").Call("setItem", "bestScore", mgr.Game.Score)
			mgr.Game.BestScore = mgr.Game.Score
			return mgr
		}
		return mgr
	} else {
		oldMgr.Game.LastUpdate = lastUpdate
	}
	return oldMgr
}

func StartGame() {
	mgr := getInitialMgr()
	tick := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		mgr = getNewMgrPropsOnTick(mgr)
		render(mgr)
		return nil
	})
	window.Set("tick", tick)
	window.Get("addEventListener").Invoke("resize", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		clearContainer()
		containerSize := getContainerSize()
		mgr.ContainerSize = containerSize
		mgr.Projector = getProjectors(containerSize, mgr.Game)
		window.Call("tick")
		return nil
	}))

	window.Get("addEventListener").Invoke("keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		keyNum := args[0].JSValue().Get("keyCode").Int()
		if keyName, ok := MOVEMENT_KEYS[keyNum]; ok {
			mgr.Game.Movement = keyName
		} else {
			mgr.Game.Movement = ""
			if keyNum == STOP_KEY {
				now := time.Now().UnixMilli()
				if mgr.Game.StopTime != 0 {
					mgr.Game.StopTime = 0
					mgr.Game.LastUpdate = 0
				} else {
					mgr.Game.StopTime = now
				}
			}
		}
		return nil
	}))

	window.Call("setInterval", tick, UPDATE_EVERY)
}
