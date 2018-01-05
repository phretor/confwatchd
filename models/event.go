package models

import (
	"encoding/json"
	"fmt"
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

	currentEdition *Edition `gorm:"-"`
}

func Events() (events []Event) {
	rows, err := db.Raw("SELECT e.* FROM events e INNER JOIN editions d on e.id = d.event_id GROUP BY e.id ORDER BY d.starts ASC").Rows()
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var event Event
		db.ScanRows(rows, &event)
		event.LoadCategories()
		events = append(events, event)
	}

	return
}

func CountEvents() (count int) {
	db.Model(&Event{}).Count(&count)
	return
}

func NextEvents(limit int) (events []Event) {
	rows, err := db.Raw("SELECT e.* FROM events e INNER JOIN editions d on e.id = d.event_id AND d.starts > ? GROUP BY e.id ORDER BY d.starts ASC LIMIT "+fmt.Sprintf("%d", limit), time.Now()).Rows()
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var event Event
		db.ScanRows(rows, &event)
		event.LoadCategories()
		events = append(events, event)
	}

	return
}

func EventsByCountry(c string, limit int) (events []Event, err error) {
	rows, err := db.Raw("SELECT e.* FROM events e INNER JOIN editions d on e.id = d.event_id AND d.country = ? GROUP BY e.id", c).Rows()
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var event Event
		db.ScanRows(rows, &event)
		event.LoadCategories()
		events = append(events, event)
	}

	return
}

func EventBySlug(slug string) (err error, event Event) {
	err = db.Where("slug = ?", slug).First(&event).Error
	if err == nil {
		err = event.LoadCategories()
	}
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

func (e *Event) LoadCategories() error {
	return db.Model(e).Related(&e.Categories, "Categories").Error
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

func (e *Event) CurrentEdition() (edition Edition) {
	if e.currentEdition == nil {
		db.Where("event_id = ?", e.ID).Order("strftime('%s', 'now') - ends ASC").First(&edition)
		edition.LoadAttributes()
		e.currentEdition = &edition
	}
	return *e.currentEdition
}

func (e *Event) CurrentEditionByCountry(c string) (edition Edition) {
	if e.currentEdition == nil {
		db.Where("event_id = ?", e.ID).Where("country = ?", c).Order("strftime('%s', 'now') - ends ASC").First(&edition)
		edition.LoadAttributes()
		e.currentEdition = &edition
	}
	return *e.currentEdition
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
	if err == nil {
		edition.LoadAttributes()
	}
	return
}

func (e Event) ByTime(starts string, ends string, order string, limit int) (editions []Edition) {
	now := time.Now()
	db.Where("event_id = ?", e.ID).Where(starts, now).Where(ends, now).Order(order).Find(&editions).Limit(limit)
	if err == nil {
		for i, _ := range editions {
			editions[i].LoadAttributes()
		}
	}
	return
}

func (e Event) Past(limit int) (past []Edition) {
	return e.ByTime("", "ends < ?", "starts desc", limit)
}

func (e Event) Present(limit int) (present []Edition) {
	return e.ByTime("starts < ?", "ends > ?", "starts desc", limit)
}

func (e Event) Future(limit int) (future []Edition) {
	return e.ByTime("starts > ?", "", "starts asc", limit)
}
