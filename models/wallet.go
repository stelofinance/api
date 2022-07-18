package models

import ()

type Wallet struct {
	ID      uint   `gorm:"primaryKey"`
	Address string `gorm:"unique;not null"`
	UserID  uint
	Webhook string
	Key     string
}
