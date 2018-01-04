package controllers

import (
	"fmt"
	"github.com/ConfWatch/confwatchd/models"
	"github.com/gin-gonic/gin"
	"github.com/pariz/gountries"
)

var (
	cQuery = gountries.New()
)

func ShowCountry(c *gin.Context) {
	country := c.Params.ByName("country_name")
	events, err := models.EventsByCountry(country, 10)
	if err != nil || len(events) == 0 {
		do404(c, "Country not found.")
		return
	}

	cName := country
	cData, err := cQuery.FindCountryByAlpha(country)
	if err == nil {
		cName = cData.Name.Common
	}

	c.HTML(200, "events/list", struct {
		SEO        SEO
		Categories []models.Category
		Countries  []string
		ListTitle  string
		Events     []models.Event
	}{
		SEO{
			Title:       fmt.Sprintf("Events in %s - ConfWatch.ninja", cName),
			Description: fmt.Sprintf("List of events in %s.", cName),
		},
		models.Categories(),
		models.Countries(),
		fmt.Sprintf("Events in %s", cName),
		events,
	})
}
