package models

type EventCategory struct {
	EventID    uint `gorm:"index"`
	CategoryID uint `gorm:"index"`
}
