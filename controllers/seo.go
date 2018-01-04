package controllers

import (
	"github.com/ConfWatch/confwatchd/models"
	"github.com/gin-gonic/gin"
)

type SEO struct {
	Title       string
	Description string
	Keywords    string
}

func defSEO() SEO {
	return SEO{
		"confwatch",
		"Discover hacking conferences around the world.",
		"hacking, hacker, conference, conf",
	}
}

func do404(c *gin.Context, message string) {
	c.HTML(404, "misc/404", struct {
		SEO        SEO
		Categories []models.Category
		Countries  []string
		Message    string
	}{
		defSEO(),
		models.Categories(),
		models.Countries(),
		message,
	})
}
