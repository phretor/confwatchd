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
		CountEditions   int
		CountEvents     int
		CountCategories int
	}{
		SEO{
			Title:       "confwatch / home",
			Description: "ConfWatch homepage.",
		},
		cats,
		models.CountEditions(),
		models.CountEvents(),
		len(cats),
	})
}
