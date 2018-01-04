package controllers

import (
	"github.com/ConfWatch/confwatchd/config"
	"github.com/ConfWatch/confwatchd/models"
	"github.com/gin-gonic/gin"
)

func ListEvents(c *gin.Context) {
	events := models.Events()
	c.HTML(200, "events/list", struct {
		SEO        SEO
		Categories []models.Category
		Countries  []string
		ListTitle  string
		Events     []models.Event
	}{
		SEO{
			Title:       "All Events - ConfWatch.ninja",
			Description: "List of events in confwatch database.",
			Version:     config.APP_VERSION,
		},
		models.Categories(),
		models.Countries(),
		"All Events",
		events,
	})
}

func ShowEvent(c *gin.Context) {
	err, event := models.EventBySlug(c.Params.ByName("event_name"))
	if err != nil {
		do404(c, "Event not found.")
	} else {
		c.HTML(200, "events/show", struct {
			SEO        SEO
			Categories []models.Category
			Countries  []string
			Event      models.Event
			Past       []models.Edition
			Present    []models.Edition
			Future     []models.Edition
		}{
			SEO{
				Title:       event.Title,
				Description: event.Description,
				Keywords:    event.Tags(),
				Version:     config.APP_VERSION,
			},
			models.Categories(),
			models.Countries(),
			event,
			event.Past(5),
			event.Present(5),
			event.Future(5),
		})
	}
}
