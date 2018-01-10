package controllers

import (
	"encoding/json"
	"github.com/ConfWatch/confwatchd/config"
	"github.com/ConfWatch/confwatchd/log"
	"github.com/ConfWatch/confwatchd/models"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"net/http"
	"time"
)

var gitCache = cache.New(15*time.Minute, 30*time.Minute)
var httpClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

type Contributor struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HtmlURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
	Contributions     int    `json:"contributions"`
}

func getContributors(repo string) (err error, contributors []Contributor) {
	cacheKey := repo + "-data-contributors-v2"
	obj, found := gitCache.Get(cacheKey)
	if found == false {
		url := "https://api.github.com/repos/ConfWatch/" + repo + "/contributors"

		log.Infof("Fetching github contributors from %s ...", url)

		err := getJson(url, &contributors)
		if err == nil {
			gitCache.Set(cacheKey, contributors, 15*time.Minute)
		}
	} else {
		contributors = obj.([]Contributor)
	}

	return
}

func AboutPage(c *gin.Context) {
	cats := models.Categories()
	countries := models.Countries()
	_, dataContribs := getContributors("confwatch-data")
	_, codeContribs := getContributors("confwatchd")

	c.HTML(200, "pages/about", struct {
		SEO              SEO
		Categories       []models.Category
		Countries        []string
		DataContributors []Contributor
		CodeContributors []Contributor
		CountEditions    int
		CountEvents      int
		CountCategories  int
		CountCountries   int
	}{
		SEO{
			Title:       "About - ConfWatch.ninja",
			Description: "About the confwatch.ninja project.",
			Version:     config.APP_VERSION,
		},
		cats,
		countries,
		dataContribs,
		codeContribs,
		models.CountEditions(),
		models.CountEvents(),
		len(cats),
		len(countries),
	})
}
