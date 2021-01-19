package route

import "time"

type RouteAddressMap struct {
	RAMID     int `json:"ramID"`
	RouteID   int `json:"routeID"`
	RequestID int `json:"requestID"`
}

type Route struct {
	RouteID          int               `json:"routeID"`
	FieldWorkerID    int               `json:"fieldWorkerID"`
	CreateRouteTime  time.Time         `json:"createRouteTime"`
	IsDone           bool              `json:"isDone"`
	IsStart          bool              `json:"isStart"`
	RouteAddressMaps []RouteAddressMap `json:"routeAddressMaps"`
}
