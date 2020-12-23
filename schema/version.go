package schema

import "time"

type VersionControl struct {
	Id              int32  `json:"id" gorm:"pk"`
	App             string `json:"app" gorm:"not null"`
	Platform        string `json:"platform" gorm:"not null"`
	Version         string `json:"version" gorm:"not null"`
	Url             string `json:"url" gorm:"not null"`
	Synchronization bool   `json:"synchronization" gorm:"not null"`
}

type UserVersion struct {
	Id        int32     `json:"id" gorm:"pk"`
	Address   string    `json:"address" gorm:"pk"`
	App       string    `json:"app" gorm:"not null"`
	Platform  string    `json:"platform" gorm:"not null"`
	Version   string    `json:"version" gorm:"not null"`
	Timestamp time.Time `json:"timestamp" gorm:"default:now()"`
}
