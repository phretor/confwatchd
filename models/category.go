package models

import (
	"encoding/json"
	"io/ioutil"
	"sort"
)

type Category struct {
	Slugable
	Description string  `json:"description" gorm:"not null;type:text"`
	Events      []Event `gorm:"many2many:event_categories;"`
}

func CategoryBySlug(slug string) (err error, c Category) {
	err = db.Where("slug = ?", slug).First(&c).Error
	return
}

func Categories() (cats []Category) {
	if err := db.Find(&cats).Error; err != nil {
		cats = make([]Category, 0)
	}

	sort.Slice(cats, func(i, j int) bool {
		return cats[i].CountEvents() > cats[j].CountEvents()
	})

	return
}

func CategoryFromFile(filename string) (err error, c Category) {
	var raw []byte
	raw, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(raw, &c)
	if err != nil {
		return
	}

	c.UpdateSlug()
	return
}

func (c *Category) UpdateFrom(b Category) {
	c.Title = b.Title
	c.Description = b.Description
}

func (c Category) Equals(b Category) bool {
	return c.Title == b.Title && c.Slug == b.Slug && c.Description == b.Description
}

func (c Category) CountEvents() (count int) {
	db.Model(EventCategory{}).Where("category_id = ?", c.ID).Count(&count)
	return
}

func (c *Category) LoadEvents(limit int) error {
	err := db.Model(c).Related(&c.Events, "Events").Limit(limit).Order("events.ends desc").Error
	if err != nil {
		return err
	}

	for i, _ := range c.Events {
		err = c.Events[i].LoadCategories()
		if err != nil {
			return err
		}
	}
	return nil
}
