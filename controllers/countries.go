package controllers

import (
	"fmt"
	"github.com/ConfWatch/confwatchd/models"
	"github.com/gin-gonic/gin"
)

func ShowCountry(c *gin.Context) {
	country := c.Params.ByName("country_name")
	events, err := models.EventsByCountry(country, 10)
	if err != nil || len(events) == 0 {
		do404(c, "Country not found.")
		return
	}

	c.HTML(200, "events/list", struct {
		SEO        SEO
		Categories []models.Category
		Countries  []string
		ListTitle  string
		Events     []models.Event
	}{
		SEO{
			Title:       fmt.Sprintf("confwatch / in %s", country),
			Description: fmt.Sprintf("List of events in the %s country.", country),
		},
		models.Categories(),
		models.Countries(),
		country,
		events,
	})
}
