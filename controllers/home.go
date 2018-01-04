package controllers

import (
	"math/rand"
	"time"

	"github.com/ConfWatch/confwatchd/config"
	"github.com/ConfWatch/confwatchd/models"
	"github.com/gin-gonic/gin"
)

func removeCategory(cats []models.Category, i int) []models.Category {
	cats[i] = cats[len(cats)-1]
	return cats[:len(cats)-1]
}

func removeCountry(countries []string, i int) []string {
	countries[i] = countries[len(countries)-1]
	return countries[:len(countries)-1]
}

type CloudElement struct {
	Type  string
	Size  float64
	Slug  string
	Title string
}

func ShowHome(c *gin.Context) {
	cats := models.Categories()
	countries := models.Countries()
	minSize := 4.0
	maxCloudElements := 50
	maxCountryCount := 0
	minCountryCount := 9999999
	maxCatCount := 0
	minCatCount := 9999999
	cloud := make([]CloudElement, 0)

	for _, c := range cats {
		count := c.CountEvents()
		if count > maxCatCount {
			maxCatCount = count
		} else if count < minCatCount {
			minCatCount = count
		}
	}

	for _, c := range countries {
		count := models.CountByCountry(c)
		if count > maxCountryCount {
			maxCountryCount = count
		} else if count < minCountryCount {
			minCountryCount = count
		}
	}

	avail_cats := make([]models.Category, len(cats))
	copy(avail_cats, cats)

	avail_countries := make([]string, len(countries))
	copy(avail_countries, countries)

	rsrc := rand.NewSource(time.Now().UnixNano())

	for {
		n := rsrc.Int63()
		if n%2 == 0 && len(avail_cats) > 0 {
			cat := avail_cats[0]
			ratio := float64(cat.CountEvents()) / float64(maxCatCount)
			size := minSize * ratio
			avail_cats = removeCategory(avail_cats, 0)

			cloud = append(cloud, CloudElement{
				Type:  "a",
				Slug:  cat.Slug,
				Title: cat.Title,
				Size:  size,
			})
		} else if len(avail_countries) > 0 {
			country := avail_countries[0]
			ratio := float64(models.CountByCountry(country)) / float64(maxCountryCount)
			size := minSize * ratio
			avail_countries = removeCountry(avail_countries, 0)

			cloud = append(cloud, CloudElement{
				Type:  "img",
				Slug:  country,
				Title: country,
				Size:  size,
			})
		}

		if len(cloud) >= maxCloudElements {
			break
		} else if len(avail_cats) == 0 && len(avail_countries) == 0 {
			break
		}
	}

	c.HTML(200, "home/index", struct {
		SEO             SEO
		Categories      []models.Category
		Countries       []string
		CountEditions   int
		CountEvents     int
		CountCategories int
		CountCountries  int
		Cloud           []CloudElement
		Next            []models.Event
	}{
		SEO{
			Title:       "Home - ConfWatch.ninja",
			Description: "ConfWatch homepage.",
			Version:     config.APP_VERSION,
		},
		cats,
		countries,
		models.CountEditions(),
		models.CountEvents(),
		len(cats),
		len(countries),
		cloud,
		models.NextEvents(25),
	})
}
