package controllers

import (
	"github.com/ConfWatch/confwatchd/models"
	"github.com/gin-gonic/gin"
)

func ListEvents(c *gin.Context) {
	events := models.Events()
	c.HTML(200, "events/list", struct {
		SEO        SEO
		Categories []models.Category
		ListTitle  string
		Events     []models.Event
	}{
		SEO{
			Title:       "confwatch / events",
			Description: "List of events in confwatch database.",
		},
		models.Categories(),
		"Events",
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
			Event      models.Event
			Past       []models.Edition
			Present    []models.Edition
			Future     []models.Edition
		}{
			SEO{
				Title:       event.Title,
				Description: event.Description,
				Keywords:    event.Tags(),
			},
			models.Categories(),
			event,
			event.Past(5),
			event.Present(5),
			event.Future(5),
		})
	}
}
