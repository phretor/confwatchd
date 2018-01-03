package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type Slugable struct {
	gorm.Model
	Title string `json:"title" gorm:"not null"`
	Slug  string `json:"slug" gorm:"not null;unique_index"`
}

func (s *Slugable) UpdateSlug() {
	s.Slug = slug.Make(s.Title)
}

func (s *Slugable) BeforeSave() (err error) {
	s.UpdateSlug()
	return nil
}
