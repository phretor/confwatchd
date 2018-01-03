package models

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"strings"
	"time"
)

const (
	EditionTypeConference = iota
	EditionTypeCamp
	EditionTypeTraining
)

type Edition struct {
	Slugable
	EventID uint
	Type    int `json:"type" gorm:"index"`

	Description string `json:"description" gorm:"not null;type:text"`
	Website     string `json:"website"`

	Country string `json:"country" gorm:"not null;index"`
	City    string `json:"city" gorm:"not null;index"`
	Address string `json:"address" gorm:"not null"`

	Starts time.Time `json:"starts" gorm:"not null; index"`
	Ends   time.Time `json:"ends" gorm:"not null; index"`

	Tags string `json:"tags" gorm:"type:text"`

	Attributes []Attribute `gorm:"many2many:edition_attributes;"`

	MetaAttributes []string `json:"attributes" gorm:"-"`
}

func EditionFromFile(filename string) (err error, edition Edition) {
	var raw []byte
	raw, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(raw, &edition)
	if err != nil {
		return
	}

	edition.UpdateSlug()
	return
}

func (e Edition) HasAttribute(a Attribute) bool {
	var ea EditionAttribute
	if err := db.Where("edition_id = ?", e.ID).Where("attribute_id = ?", a.ID).First(&ea).Error; err != nil {
		return false
	}
	return true
}

func (e Edition) AddAttribute(tx *gorm.DB, a Attribute) error {
	return tx.Create(&EditionAttribute{
		EditionID:   e.ID,
		AttributeID: a.ID,
	}).Error
}

func (e Edition) Equals(b Edition) bool {
	if e.Slug != b.Slug {
		return false
	} else if e.Title != b.Title {
		return false
	} else if e.Type != b.Type {
		return false
	} else if e.Description != b.Description {
		return false
	} else if e.Website != b.Website {
		return false
	} else if e.Country != b.Country {
		return false
	} else if e.City != b.City {
		return false
	} else if e.Address != b.Address {
		return false
	} else if e.Starts.Unix() != b.Starts.Unix() {
		return false
	} else if e.Ends.Unix() != b.Ends.Unix() {
		return false
	} else if e.Tags != b.Tags {
		return false
	}
	return true
}

func (e *Edition) UpdateFrom(b Edition) {
	e.Title = b.Title
	e.Type = b.Type
	e.Description = b.Description
	e.Website = b.Website
	e.Country = b.Country
	e.City = b.City
	e.Address = b.Address
	e.Starts = b.Starts
	e.Ends = b.Ends
	e.Tags = b.Tags
}

func IsValidEditionType(t int) bool {
	return t == EditionTypeConference || t == EditionTypeCamp || t == EditionTypeTraining
}

func (e *Edition) BeforeSave() (err error) {
	e.Slugable.BeforeSave()

	if IsValidEditionType(e.Type) == false {
		err = errors.New("Invalid type.")
	}

	if e.Starts.After(e.Ends) {
		err = errors.New("End date is before start date.")
	}

	unique := make([]string, 0)
	tmp := make(map[string]bool, 0)
	tags := strings.Split(e.Tags, ",")
	e.Tags = ""

	for _, t := range tags {
		t = strings.Trim(t, "\r\n\t ")
		if len(t) == 0 {
			continue
		}

		if found, _ := tmp[t]; found == false {
			tmp[t] = true
			unique = append(unique, t)
		}
	}

	e.Tags = strings.Join(unique, ",")

	return
}
