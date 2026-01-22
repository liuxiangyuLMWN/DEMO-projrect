package main

import "math/rand"

type VehicleStatus string

const (
	StatusIdle       VehicleStatus = "idle"
	StatusPicking    VehicleStatus = "picking"
	StatusDelivering VehicleStatus = "delivering"
)

type Vehicle struct {
	ID       string
	Position Position
	Status   VehicleStatus
	IsOnline bool
	Target   *Position
}

func NewVehicle(id string, m *Map) *Vehicle {
	return &Vehicle{
		ID:       id,
		Position: m.RandomPosition(),
		Status:   StatusIdle,
		IsOnline: true,
	}
}

// 只负责“走一格”
func (v *Vehicle) Move() {
	if v.Target == nil {
		return
	}

	// 计算下一步位置
	newX := v.Position.X
	newY := v.Position.Y

	if newX < v.Target.X {
		newX++
	} else if newX > v.Target.X {
		newX--
	} else if newY < v.Target.Y { // 只有X对齐后才移动Y
		newY++
	} else if newY > v.Target.Y {
		newY--
	}

	// 边界检查（虽然理论上Target应该在地图内，但为了安全）
	if newX >= 0 && newX < 20 && newY >= 0 && newY < 20 { // 假设地图大小为20x20
		v.Position.X = newX
		v.Position.Y = newY
	}
}

// AcceptOrder 车辆决定是否接受订单，返回 (是否接受, 决策延迟秒数)
func (v *Vehicle) AcceptOrder() bool {
	// 90% 概率接受订单
	accept := rand.Float64() < 0.9
	return accept
}
