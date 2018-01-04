package controllers

import (
	"fmt"
	"github.com/ConfWatch/confwatchd/config"
	"github.com/ConfWatch/confwatchd/models"
	"github.com/gin-gonic/gin"
	"strings"
)

func ShowEdition(c *gin.Context) {
	err, event := models.EventBySlug(c.Params.ByName("event_name"))
	if err != nil {
		do404(c, "Event not found.")
		return
	}

	err, edition := event.EditionBySlug(c.Params.ByName("edition_name"))
	if err != nil {
		do404(c, "Edition not found.")
		return
	}

	edition.LoadAttributes()

	tags := strings.Split(edition.Tags, ",")
	stags := make([]string, len(tags))

	for i, t := range tags {
		stags[i] = "#" + t
	}

	socialStream := strings.Join(stags, " OR ")

	c.HTML(200, "editions/show", struct {
		SEO          SEO
		Categories   []models.Category
		Countries    []string
		Event        models.Event
		Edition      models.Edition
		Tags         []string
		SocialStream string
	}{
		SEO{
			Title:       fmt.Sprintf("%s / %s", event.Title, edition.Title),
			Description: edition.Description,
			Keywords:    edition.Tags,
			Version:     config.APP_VERSION,
		},
		models.Categories(),
		models.Countries(),
		event,
		edition,
		tags,
		socialStream,
	})
}
