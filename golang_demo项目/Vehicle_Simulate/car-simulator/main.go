package main

import (
	"fmt"
	"math/rand"
	"time"
)

// SimulationConfig 模拟配置
type SimulationConfig struct {
	MapWidth           int
	MapHeight          int
	VehicleCount       int
	OrderCount         int
	RetryTime          int
	Delay              int
	HourlyDistribution []float64
}

func main() {
	//使用随机种子来仿真不同的时刻
	rand.Seed(time.Now().UnixNano())

	m := &Map{Width: 20, Height: 20}

	// 24小时订单分布
	// hourlyOrders := []int{245, 198, 152, 131, 173, 260, 377, 473, 505, 478, 477, 483, 509, 484, 470, 531, 607, 678, 650, 529, 508, 446, 369, 265}
	hourlyOrders := []int{20, 19, 15, 13, 17, 26, 37, 47, 50, 47, 47, 48, 50, 48, 47, 53, 60, 67, 65, 52, 50, 44, 36, 26}

	// 计算总订单数
	totalOrders := 0
	for _, count := range hourlyOrders {
		totalOrders += count
	}

	vehicleCount := 5
	fmt.Printf("部署 %d 辆车来处理订单\n", vehicleCount)

	vehicles := []*Vehicle{}
	for i := 0; i < vehicleCount; i++ {
		vehicles = append(vehicles, NewVehicle(fmt.Sprintf("V%d", i), m))
	}

	//生成订单
	orders := []*Order{}
	orderID := 0
	for hour := 0; hour < 24; hour++ {
		hourlyCount := hourlyOrders[hour]
		// 在当前小时内随机分配订单创建时间
		for i := 0; i < hourlyCount; i++ {
			// 在当前小时内随机生成秒数 (0-3599秒)
			secondsInHour := rand.Intn(360)
			createdAt := int64(hour*360 + secondsInHour)

			orders = append(orders, &Order{
				ID:               fmt.Sprintf("O%d", orderID),
				CreatedAt:        createdAt,
				Pickup:           m.RandomPosition(),
				Dropoff:          m.RandomPosition(),
				Status:           StatusCreated,
				NextScheduleTime: createdAt, // 下次可调度时间初始时等于创建时间
			})
			orderID++
		}
		fmt.Printf("第 %2d 小时: 生成 %d 个订单\n", hour+1, hourlyCount)
	}
	//模拟器
	sim := &Simulator{
		Map:       m,
		Vehicles:  vehicles,
		Orders:    orders,
		Algorithm: &NearestVehicleAlgorithm{},
	}

	sim.Run()

	var sum int64
	var count int
	for _, o := range orders {
		if o.Status == StatusDelivered {
			sum += o.PickupTime - o.CreatedAt + o.SchedulingDelay
			count++
		}
	}

	// 计算取消的订单数
	cancelledCount := 0
	for _, o := range orders {
		if o.Status == StatusCancelled {
			cancelledCount++
		}
	}

	fmt.Println("\n=== 最终统计报告 ===")
	fmt.Printf("完成订单数: %d/%d (%.1f%%)\n", count, len(orders), float64(count)/float64(len(orders))*100)
	fmt.Printf("取消订单数: %d\n", cancelledCount)
	if count > 0 {
		fmt.Printf("平均取货时间: %.2f秒\n", float64(sum)/float64(count))
	}

	// 详细的时间统计
	var deliveryTimes []int64
	var pickupTimes []int64
	for _, o := range orders {
		if o.Status == StatusDelivered {
			deliveryTimes = append(deliveryTimes, o.DeliveryTime-o.CreatedAt+o.SchedulingDelay)
			pickupTimes = append(pickupTimes, o.PickupTime-o.CreatedAt+o.SchedulingDelay)
		}
	}

	if len(deliveryTimes) > 0 {
		var totalDeliveryTime int64
		var totalPickupTime int64
		minDeliveryTime := deliveryTimes[0]
		maxDeliveryTime := deliveryTimes[0]
		minPickupTime := pickupTimes[0]
		maxPickupTime := pickupTimes[0]

		for i, dt := range deliveryTimes {
			totalDeliveryTime += dt
			totalPickupTime += pickupTimes[i]
			if dt < minDeliveryTime {
				minDeliveryTime = dt
			}
			if dt > maxDeliveryTime {
				maxDeliveryTime = dt
			}
			if pickupTimes[i] < minPickupTime {
				minPickupTime = pickupTimes[i]
			}
			if pickupTimes[i] > maxPickupTime {
				maxPickupTime = pickupTimes[i]
			}
		}

		fmt.Printf("总配送时间: 平均=%.2f秒, 最短=%d秒, 最长=%d秒\n",
			float64(totalDeliveryTime)/float64(len(deliveryTimes)), minDeliveryTime, maxDeliveryTime)
		fmt.Printf("取货时间: 平均=%.2f秒, 最短=%d秒, 最长=%d秒\n",
			float64(totalPickupTime)/float64(len(pickupTimes)), minPickupTime, maxPickupTime)
	}

	fmt.Println("\n模拟完成！")
}
