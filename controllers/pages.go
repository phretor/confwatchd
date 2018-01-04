package controllers

import (
	"github.com/ConfWatch/confwatchd/config"
	"github.com/ConfWatch/confwatchd/models"
	"github.com/gin-gonic/gin"
)

func AboutPage(c *gin.Context) {
	cats := models.Categories()
	countries := models.Countries()

	c.HTML(200, "pages/about", struct {
		SEO             SEO
		Categories      []models.Category
		Countries       []string
		CountEditions   int
		CountEvents     int
		CountCategories int
		CountCountries  int
	}{
		SEO{
			Title:       "About - ConfWatch.ninja",
			Description: "About the confwatch.ninja project.",
			Version:     config.APP_VERSION,
		},
		cats,
		countries,
		models.CountEditions(),
		models.CountEvents(),
		len(cats),
		len(countries),
	})
}
