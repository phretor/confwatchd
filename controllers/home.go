package controllers

import (
	"github.com/ConfWatch/confwatchd/models"
	"github.com/gin-gonic/gin"
)

func ShowHome(c *gin.Context) {
	cats := models.Categories()

	c.HTML(200, "home/index", struct {
		SEO             SEO
		Categories      []models.Category
		Countries       []string
		CountEditions   int
		CountEvents     int
		CountCategories int
		Next            []models.Event
	}{
		SEO{
			Title:       "Home - ConfWatch.ninja",
			Description: "ConfWatch homepage.",
		},
		cats,
		models.Countries(),
		models.CountEditions(),
		models.CountEvents(),
		len(cats),
		models.NextEvents(25),
	})
}
