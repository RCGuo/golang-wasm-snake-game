package main

import (
	"math"
	"math/rand"
)

type Vector struct {
	x float64
	y float64
}

func (v Vector) ScaleBy(number float64) Vector {
	return Vector{v.x * number, v.y * number}
}

func (v Vector) Add(input Vector) Vector {
	return Vector{v.x + input.x, v.y + input.y}
}

func (v Vector) Subtract(input Vector) Vector {
	return Vector{v.x - input.x, v.y - input.y}
}

func (v Vector) Length() float64 {
	return math.Hypot(v.x, v.y)
}

func (v Vector) Normalize() Vector {
	return v.ScaleBy(1 / v.Length())
}

func (v Vector) IsOpposite(input Vector) bool {
	newV := v.Add(input)
	return AreEqual(newV.x, 0) && AreEqual(newV.y, 0)
}

func (v Vector) EqualTo(input Vector) bool {
	return AreEqual(v.x, input.x) && AreEqual(v.y, input.y)
}

func GetRandomFrom(vList []Vector) Vector {
	return vList[int(math.Floor(rand.Float64()*float64(len(vList))))]
}

func GetWithoutLastElement(vList []Vector) []Vector {
	return vList[0:(len(vList) - 1)]
}

func GetLastElement(vList []Vector) Vector {
	return vList[len(vList)-1]
}
