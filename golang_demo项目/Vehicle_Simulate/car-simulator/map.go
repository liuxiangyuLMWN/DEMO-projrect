package main

import (
	"math"
	"math/rand"
)

// Position 表示网格中的位置坐标
type Position struct {
	X, Y int
}

// Map 表示网格地图
type Map struct {
	Width, Height int
}

// Distance 计算两个位置之间的曼哈顿距离
func Distance(p1, p2 Position) int {
	dx := p1.X - p2.X
	if dx < 0 {
		dx = -dx
	}
	dy := p1.Y - p2.Y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

// EuclideanDistance 计算两个位置之间的欧几里得距离
func EuclideanDistance(p1, p2 Position) float64 {
	dx := float64(p1.X - p2.X)
	dy := float64(p1.Y - p2.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

// RandomPosition 在地图上生成一个随机位置
func (m *Map) RandomPosition() Position {
	return Position{
		X: rand.Intn(m.Width),
		Y: rand.Intn(m.Height),
	}
}
