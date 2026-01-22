package main

type OrderStatus string

const (
	StatusCreated   OrderStatus = "created"   //创建
	StatusScheduled OrderStatus = "scheduled" //被调度
	StatusPicked    OrderStatus = "picked"    //被拿起
	StatusDelivered OrderStatus = "delivered" //被送达
	StatusCancelled OrderStatus = "cancelled" //被取消
)

type Order struct {
	ID              string
	CreatedAt       int64 //创建时间
	Pickup          Position
	Dropoff         Position
	Status          OrderStatus
	AssignedTo      *string //分配给
	PickupTime      int64   //拿起时间
	DeliveryTime    int64   //完成时间
	SchedulingDelay int64   //调度时间
	RetryCount      int     //重复调度次数
	// LastRetryAt      int64
	NextScheduleTime int64 // 下次可调度时间
}
