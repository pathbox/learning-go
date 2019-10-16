package main

import "time"

type User struct {
	Id      int       `gorm:"column:id"`
	Name    string    `gorm:"column:name"`
	Status  string    `gorm:"column:status"`
	Created time.Time `gorm:"column:created"`
}

func (User) TableName() string {
	return "User"
}

func (User) SchemaName() string {
	return "Test"
}