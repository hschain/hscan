package models

type Frame struct {
	UsersNumber            int32   `json:"usersNumber"`
	Tps                    int32   `json:"tps"`
	CurrentDayProvisions   float64 `json:"current_day_provisions"`
	TotalCirculationSupply int64   `json:"total_circulation_supply"`
}
