package main

import (
	"fmt"
	"os"
)

func flushOutput() {
	os.Stdout.Sync()
}

type Simulator struct {
	Map                 *Map
	Vehicles            []*Vehicle
	Orders              []*Order
	Algorithm           DispatchAlgorithm
	CurrentTime         int64
	totalActiveVehicles int // 累积的活跃车辆数
	totalHours          int // 统计的小时数
}

func (s *Simulator) Run() {
	totalSeconds := int64(24 * 360) // 将1小时设为360秒，总共8640秒

	fmt.Printf(" 开始车辆调度模拟: 24小时(%d秒), %d辆车, %d个订单\n",
		totalSeconds, len(s.Vehicles), len(s.Orders))
	for s.CurrentTime < totalSeconds {
		//主流程
		s.moveVehicles()
		s.handleOrders()

		// 累积车辆活跃状态（每秒统计一次）
		activeCount := 0
		for _, v := range s.Vehicles {
			if v.Status != StatusIdle {
				activeCount++
			}
		}
		s.totalActiveVehicles += activeCount
		s.totalHours++ // 这里totalHours实际上是总秒数

		// 每小时输出一次统计信息 (360秒 = 1小时)
		if s.CurrentTime > 0 && s.CurrentTime%360 == 0 {
			s.printHourlyStats()
		}
		s.CurrentTime++

	}
}

func (s *Simulator) moveVehicles() {
	for _, v := range s.Vehicles {
		if !v.IsOnline {
			continue
		}
		v.Move()
		// 到达取货点
		if v.Status == StatusPicking &&
			v.Target != nil &&
			v.Position.X == v.Target.X && v.Position.Y == v.Target.Y {
			// found := false
			for _, o := range s.Orders {
				if o.AssignedTo != nil && *o.AssignedTo == v.ID && o.Status == StatusScheduled {
					o.Status = StatusPicked
					o.PickupTime = s.CurrentTime
					v.Status = StatusDelivering
					target := o.Dropoff
					v.Target = &target
					// found = true
					fmt.Printf("[取货] 车辆 %s(位置:%d,%d) 在时间 %d 完成取货，订单 %s 状态: %s -> %s，取货耗时: %d秒，开始送货到 (%d,%d)\n",
						v.ID, v.Position.X, v.Position.Y, s.CurrentTime, o.ID, StatusScheduled, StatusPicked, o.PickupTime-o.CreatedAt+o.SchedulingDelay, o.Dropoff.X, o.Dropoff.Y)
					flushOutput()
					break
				}
			}
			// 如果没有找到匹配的订单，重置车辆状态（防止卡死）
			// if !found {
			// 	fmt.Printf("[警告] 车辆 %s 在时间 %d 到达取货点但未找到匹配订单，重置状态\n", v.ID, s.CurrentTime)
			// 	v.Status = StatusIdle
			// 	v.Target = nil
			// }
		}

		// 到达送达点
		if v.Status == StatusDelivering &&
			v.Target != nil &&
			v.Position.X == v.Target.X && v.Position.Y == v.Target.Y {

			// found := false
			for _, o := range s.Orders {
				if o.AssignedTo != nil && *o.AssignedTo == v.ID && o.Status == StatusPicked {
					o.Status = StatusDelivered
					o.DeliveryTime = s.CurrentTime
					v.Status = StatusIdle
					v.Target = nil
					// found = true
					fmt.Printf("✅ 车辆 %s(位置:%d,%d) 在时间 %d 完成送货，订单 %s 状态: %s -> %s，总耗时: %d秒\n",
						v.ID, v.Position.X, v.Position.Y, s.CurrentTime, o.ID, StatusPicked, StatusDelivered,
						o.DeliveryTime-o.CreatedAt+o.SchedulingDelay)
					flushOutput()
					break
				}
			}
			// 如果没有找到匹配的订单，重置车辆状态（防止卡死）
			// if !found {
			// 	fmt.Printf("[警告] 车辆 %s 在时间 %d 到达送货点但未找到匹配订单，重置状态\n", v.ID, s.CurrentTime)
			// 	v.Status = StatusIdle
			// 	v.Target = nil
			// }
		}

	}
}

func (s *Simulator) handleOrders() {
	const maxRetries = 3 // 最大重试次数
	const retryDelay = 5 // 每次重试延迟5秒

	for _, o := range s.Orders {
		// 只处理未分配的Created状态订单，且到达调度时间
		if o.Status == StatusCreated && o.AssignedTo == nil &&
			o.NextScheduleTime <= s.CurrentTime && o.RetryCount < maxRetries {
			v, delay := s.Algorithm.SelectVehicle(o, s.Vehicles)
			o.SchedulingDelay = delay
			if v == nil || v.Status != StatusIdle {
				// 调度失败：增加重试次数，设置下次调度时间
				o.RetryCount++
				// o.LastRetryAt = s.CurrentTime
				o.NextScheduleTime = s.CurrentTime + retryDelay + o.SchedulingDelay

				if o.RetryCount < maxRetries {
					fmt.Printf("[重试] 订单 %s 调度失败 (第%d次)，调度时间为%d秒，将在第 %d 秒后重试\n",
						o.ID, o.RetryCount, o.SchedulingDelay, o.NextScheduleTime)
				} else {
					fmt.Printf("❌ 订单 %s 调度失败，已达到最大重试次数 (%d)，取消订单\n",
						o.ID, maxRetries)
					o.Status = StatusCancelled
				}
				continue
			}

			// 调度成功：原子状态转换
			o.Status = StatusScheduled
			o.AssignedTo = &v.ID
			o.RetryCount = 0 // 重置重试次数
			v.Status = StatusPicking
			target := o.Pickup
			v.Target = &target

			fmt.Printf("[调度] 订单 %s 到车辆 %s(位置:%d,%d) (时间: %d秒，取货点: (%d,%d)，调度时间%d秒)\n",
				o.ID, v.ID, v.Position.X, v.Position.Y, s.CurrentTime, o.Pickup.X, o.Pickup.Y, o.SchedulingDelay)
			flushOutput()
		}
	}
}

// printHourlyStats 输出每小时的统计信息
func (s *Simulator) printHourlyStats() {
	hour := s.CurrentTime/360 + 1
	var activeOrders, completedOrders, scheduledOrders, pickedOrders, cancelledOrders int

	for _, o := range s.Orders {
		switch o.Status {
		case StatusCreated:
			// 还未调度的订单
		case StatusScheduled:
			scheduledOrders++
		case StatusPicked:
			pickedOrders++
		case StatusDelivered:
			completedOrders++
		case StatusCancelled:
			cancelledOrders++
		}
		if o.CreatedAt <= s.CurrentTime {
			activeOrders++
		}
	}

	var idleVehicles, pickingVehicles, deliveringVehicles int
	for _, v := range s.Vehicles {
		switch v.Status {
		case StatusIdle:
			idleVehicles++
		case StatusPicking:
			pickingVehicles++
		case StatusDelivering:
			deliveringVehicles++
		}
	}

	fmt.Printf("\n=== 第 %d 小时统计 (时间: %d秒) ===\n", hour, s.CurrentTime)
	fmt.Printf("订单状态: 活跃=%d, 已调度=%d, 已取货=%d, 已完成=%d, 已取消=%d\n",
		activeOrders, scheduledOrders, pickedOrders, completedOrders, cancelledOrders)
	fmt.Printf("车辆状态: 空闲=%d, 取货中=%d, 送货中=%d\n",
		idleVehicles, pickingVehicles, deliveringVehicles)
	fmt.Printf("完成率: %.1f%%\n", float64(completedOrders)/float64(len(s.Orders))*100)
	fmt.Println()
	flushOutput()
}
