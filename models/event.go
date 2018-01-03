package models

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"time"
)

type Event struct {
	Slugable
	Description  string   `json:"description" gorm:"not null;type:text"`
	Website      string   `json:"website"`
	MetaEditions []string `json:"editions" gorm:"-"`
	Editions     []Edition

	Categories     []Category `gorm:"many2many:event_categories;"`
	MetaCategories []string   `json:"categories" gorm:"-"`
}

func Events() (events []Event) {
	if err := db.Find(&events).Error; err != nil {
		events = make([]Event, 0)
	}
	return
}

func EventBySlug(slug string) (err error, event Event) {
	err = db.Where("slug = ?", slug).First(&event).Error
	return
}

func EventFromFile(filename string) (err error, event Event) {
	var raw []byte
	raw, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(raw, &event)
	if err != nil {
		return
	}

	event.UpdateSlug()
	return
}

func (e Event) Equals(b Event) bool {
	if e.Slug != b.Slug {
		return false
	} else if e.Title != b.Title {
		return false
	} else if e.Website != b.Website {
		return false
	} else if e.Description != b.Description {
		return false
	}
	return true
}

func (e *Event) UpdateFrom(b Event) {
	e.Title = b.Title
	e.Website = b.Website
	e.Description = b.Description
}

func (e Event) HasCategory(c Category) bool {
	var ec EventCategory
	if err := db.Where("event_id = ?", e.ID).Where("category_id = ?", c.ID).First(&ec).Error; err != nil {
		return false
	}
	return true
}

func (e Event) AddCategory(tx *gorm.DB, c Category) error {
	return tx.Create(&EventCategory{
		EventID:    e.ID,
		CategoryID: c.ID,
	}).Error
}

func (e Event) Tags() string {
	var edition Edition
	if err := db.Where("event_id = ?", e.ID).Order("ends desc").First(&edition).Error; err == nil {
		return edition.Tags
	}
	return ""
}

func (e *Event) EditionBySlug(slug string) (err error, edition Edition) {
	err = db.Where("event_id = ?", e.ID).Where("slug = ?", slug).First(&edition).Error
	return
}

func (e *Event) Past(limit int) (past []Edition) {
	db.Where("event_id = ?", e.ID).Where("ends < ?", time.Now()).Order("ends asc").Find(&past).Limit(limit)
	return
}

func (e *Event) Present(limit int) (past []Edition) {
	now := time.Now()
	db.Where("event_id = ?", e.ID).Where("starts < ?", now).Where("ends > ?", now).Order("starts desc").Find(&past).Limit(limit)
	return
}

func (e *Event) Future(limit int) (past []Edition) {
	db.Where("event_id = ?", e.ID).Where("starts > ?", time.Now()).Order("starts asc").Find(&past).Limit(limit)
	return
}
