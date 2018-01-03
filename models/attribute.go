package models

import (
	"encoding/json"
	"io/ioutil"
)

type Attribute struct {
	Slugable
	Description string    `json:"description" gorm:"not null;type:text"`
	Editions    []Edition `gorm:"many2many:edition_attributes;"`
}

func AttributeBySlug(slug string) (err error, attr Attribute) {
	err = db.Where("slug = ?", slug).First(&attr).Error
	return
}

func AttributeFromFile(filename string) (err error, attr Attribute) {
	var raw []byte
	raw, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(raw, &attr)
	if err != nil {
		return
	}

	attr.UpdateSlug()
	return
}

func (a *Attribute) UpdateFrom(b Attribute) {
	a.Title = b.Title
	a.Description = b.Description
}

func (a Attribute) Equals(b Attribute) bool {
	return a.Title == b.Title && a.Slug == b.Slug && a.Description == b.Description
}
