package controllers

import (
	"github.com/ConfWatch/confwatchd/models"
	"github.com/gin-gonic/gin"
)

func ShowHome(c *gin.Context) {
	cats := models.Categories()
	countries := models.Countries()

	c.HTML(200, "home/index", struct {
		SEO             SEO
		Categories      []models.Category
		Countries       []string
		CountEditions   int
		CountEvents     int
		CountCategories int
		CountCountries  int
		Next            []models.Event
	}{
		SEO{
			Title:       "Home - ConfWatch.ninja",
			Description: "ConfWatch homepage.",
		},
		cats,
		countries,
		models.CountEditions(),
		models.CountEvents(),
		len(cats),
		len(countries),
		models.NextEvents(25),
	})
}
