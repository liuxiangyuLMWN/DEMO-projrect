package main

import "math/rand"

type DispatchAlgorithm interface {
	SelectVehicle(order *Order, vehicles []*Vehicle) (*Vehicle, int64)
}

type NearestVehicleAlgorithm struct{}

func (a *NearestVehicleAlgorithm) SelectVehicle(order *Order, vehicles []*Vehicle) (*Vehicle, int64) {
	var best *Vehicle
	minDist := 1 << 30
	delay := int64(rand.Intn(10) + 1)
	for _, v := range vehicles {
		if !v.IsOnline || v.Status != StatusIdle {
			continue
		}
		d := EuclideanDistance(order.Pickup, v.Position)
		if d <= 10.0 && int(d) < minDist { // 欧几里得距离 <= 10
			minDist = int(d)
			best = v
		}
	}

	if best != nil {
		accepted := best.AcceptOrder()
		if accepted {
			return best, int64(delay)
		}
		// 如果不接受，返回nil和延迟信息（由调用方处理）
		return nil, int64(delay)
	}
	return nil, int64(delay)
}
