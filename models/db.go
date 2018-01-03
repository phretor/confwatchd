package models

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/mssql"
	// _ "github.com/jinzhu/gorm/dialects/mysql"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/ConfWatch/confwatchd/config"
	"github.com/ConfWatch/confwatchd/log"
	"github.com/ConfWatch/confwatchd/utils"
)

var err error
var db *gorm.DB

func Setup(c config.DatabaseConfig) error {
	db, err = gorm.Open(c.Type, c.Connection)
	if err != nil {
		return err
	}

	db.AutoMigrate(&Event{})
	db.AutoMigrate(&Edition{})
	db.AutoMigrate(&EditionAttribute{})
	db.AutoMigrate(&Attribute{})

	return nil
}

func seedAttributes(folder string) (err error, attributes []uint) {
	log.Debugf("Importing attributes from %s ...", folder)

	matches, err := filepath.Glob(filepath.Join(folder, "*.json"))
	if err != nil {
		return
	}

	tx := db.Begin()

	attributes = make([]uint, 0)

	for _, filename := range matches {
		var attr Attribute

		log.Debugf("Loading %s ...", filename)

		err, attr = AttributeFromFile(filename)
		if err != nil {
			tx.Rollback()
			return
		}

		var existing Attribute
		err, existing = AttributeBySlug(attr.Slug)
		if err == nil {
			if existing.Equals(attr) == false {
				log.Infof("Updating attribute %s.", log.Bold(attr.Slug))

				existing.UpdateFrom(attr)

				err = tx.Save(&existing).Error
				if err != nil {
					tx.Rollback()
					return
				}
			} else {
				log.Debugf("Attribute %s already exists.", attr.Slug)
			}

			attributes = append(attributes, existing.ID)
		} else {
			log.Infof("Creating attribute %s.", log.Bold(attr.Slug))

			err = tx.Create(&attr).Error
			if err != nil {
				tx.Rollback()
				return
			}

			attributes = append(attributes, attr.ID)
		}
	}

	err = tx.Commit().Error
	return
}

func seedEditions(tx *gorm.DB, folder string, event *Event, editions *[]uint) (err error) {
	if utils.IsFolder(folder) == false {
		return nil
	}

	log.Debugf("Importing %s editions from %s ...", log.Bold(event.Slug), folder)

	matches, err := filepath.Glob(filepath.Join(folder, "*.json"))
	if err != nil {
		return
	}

	attributes := make([]uint, 0)

	for _, filename := range matches {
		var edition Edition

		log.Debugf("Loading %s ...", filename)

		err, edition = EditionFromFile(filename)
		if err != nil {
			return
		}

		pedition := (*Edition)(nil)

		var existing Edition
		err, existing = event.EditionBySlug(edition.Slug)
		if err == nil {
			if existing.Equals(edition) == false {
				log.Infof("Updating edition %s ...", log.Bold(edition.Slug))

				existing.UpdateFrom(edition)

				err = tx.Save(&existing).Error
				if err != nil {
					return
				}
			} else {
				log.Debugf("Edition %s already exists.", edition.Slug)
			}

			pedition = &existing
			*editions = append(*editions, existing.ID)
		} else {
			log.Infof("Creating edition %s for event %s ...", log.Bold(edition.Slug), log.Bold(event.Slug))

			edition.EventID = event.ID

			err = tx.Create(&edition).Error
			if err != nil {
				return
			}

			pedition = &edition
			*editions = append(*editions, edition.ID)
		}

		for _, attributeName := range edition.MetaAttributes {
			var attr Attribute
			err, attr = AttributeBySlug(attributeName)
			if err != nil {
				log.Errorf("Attribute %s not found.", log.Bold(attributeName))
				return
			}

			if pedition.HasAttribute(attr) == false {
				log.Infof("Adding attribute %s to %s", log.Bold(attributeName), log.Bold(pedition.Slug))
				err = pedition.AddAttribute(tx, attr)
				if err != nil {
					return
				}
			}

			attributes = append(attributes, attr.ID)
		}

		var eaToPrune []EditionAttribute
		err = db.Where("edition_id = ?", pedition.ID).Not("attribute_id", attributes).Find(&eaToPrune).Error
		if err != nil {
			tx.Rollback()
			return
		}

		for _, ea := range eaToPrune {
			log.Infof("Unsetting attribute %d from edition %s.", ea.AttributeID, log.Bold(edition.Slug))
			tx.Where("edition_id = ?", ea.EditionID).Where("attribute_id = ?", ea.AttributeID).Delete(&ea)
		}
	}

	return nil
}

func seedEvents(folder string) (err error, events []uint) {
	log.Debugf("Importing events from %s ...", folder)

	matches, err := filepath.Glob(filepath.Join(folder, "*/event.json"))
	if err != nil {
		return
	}

	tx := db.Begin()

	editions := make([]uint, 0)
	events = make([]uint, 0)

	for _, filename := range matches {
		var event Event

		log.Debugf("Loading %s ...", filename)

		err, event = EventFromFile(filename)
		if err != nil {
			tx.Rollback()
			return
		}

		pevent := (*Event)(nil)

		var existing Event
		err, existing = EventBySlug(event.Slug)
		if err == nil {
			if existing.Equals(event) == false {
				log.Infof("Updating event %s ...", log.Bold(event.Slug))

				existing.UpdateFrom(event)

				err = tx.Save(&existing).Error
				if err != nil {
					tx.Rollback()
					return
				}
			} else {
				log.Debugf("Event %s already exists.", event.Slug)
			}

			pevent = &existing
			events = append(events, existing.ID)
		} else {
			log.Infof("Creating event %s ...", log.Bold(event.Slug))

			err = tx.Create(&event).Error
			if err != nil {
				tx.Rollback()
				return
			}

			pevent = &event
			events = append(events, event.ID)
		}

		editionsFolder := filepath.Join(folder, event.Slug, "editions")

		err = seedEditions(tx, editionsFolder, pevent, &editions)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	var editionsToPrune []Edition
	err = db.Not("id", editions).Find(&editionsToPrune).Error
	if err != nil {
		tx.Rollback()
		return
	}

	for _, edition := range editionsToPrune {
		log.Infof("Deleting edition %s.", log.Bold(edition.Slug))
		tx.Delete(&edition)
	}

	err = tx.Commit().Error
	return
}

func Seed(folder string) (err error) {
	folder, err = utils.ExpandPath(folder)
	if err != nil {
		return err
	}

	attrsFolder := path.Join(folder, "attributes")
	eventsFolder := path.Join(folder, "events")

	if utils.IsFolder(attrsFolder) == false {
		return fmt.Errorf("Folder %s does not exist.", attrsFolder)
	} else if utils.IsFolder(eventsFolder) == false {
		return fmt.Errorf("Folder %s does not exist.", eventsFolder)
	}

	err, attributes := seedAttributes(attrsFolder)
	if err != nil {
		return
	}

	var attrsToPrune []Attribute
	err = db.Not("id", attributes).Find(&attrsToPrune).Error
	if err != nil {
		return
	}

	for _, attr := range attrsToPrune {
		log.Infof("Deleting attribute %s.", log.Bold(attr.Slug))
		db.Delete(&attr)
	}

	err, events := seedEvents(eventsFolder)
	if err != nil {
		return
	}

	var eventsToPrune []Event
	err = db.Not("id", events).Find(&eventsToPrune).Error
	if err != nil {
		return
	}

	for _, event := range eventsToPrune {
		log.Infof("Deleting event %s.", log.Bold(event.Slug))
		db.Delete(&event)
	}

	return
}

func Close() {
	db.Close()
}
