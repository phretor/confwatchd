package models

type EditionAttribute struct {
	EditionID   uint `gorm:"index"`
	AttributeID uint `gorm:"index"`
}
