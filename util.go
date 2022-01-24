package main

import "math"

func GetRange(length float64) []int {
	rangeList := make([]int, int(length))
	for i := range rangeList {
		rangeList[i] = i
	}
	return rangeList
}

func AreEqual(one, another float64) bool {
	return math.Abs(one-another) < 0.0000001
}

func ConsoleLog(args ...interface{}) {
	window.Get("console").Call("log", args...)
}