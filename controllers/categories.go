package controllers

import (
	"fmt"
	"github.com/ConfWatch/confwatchd/log"
	"github.com/ConfWatch/confwatchd/models"
	"github.com/gin-gonic/gin"
)

func ShowCategory(c *gin.Context) {
	err, cat := models.CategoryBySlug(c.Params.ByName("cat_name"))
	if err != nil {
		do404(c, "Category not found.")
		return
	}

	err = cat.LoadEvents(10)
	if err != nil {
		log.Errorf("Could not load category %s events: %s.", cat.Slug, err)
		do404(c, "WTF?!")
		return
	}

	c.HTML(200, "events/list", struct {
		SEO       SEO
		ListTitle string
		Events    []models.Event
	}{
		SEO{
			Title:       fmt.Sprintf("confwatch / %s", cat.Title),
			Description: fmt.Sprintf("List of events in the %s category.", cat.Title),
		},
		fmt.Sprintf("%s category", cat.Title),
		cat.Events,
	})
}
