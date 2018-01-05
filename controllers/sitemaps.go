package controllers

import (
	"fmt"
	"github.com/ConfWatch/confwatchd/models"
	"github.com/gin-gonic/gin"
	"time"
)

const (
	baseUrl    = "https://confwatch.ninja"
	staticDate = "2018-01-05T17:04:54+01:00"
)

func doXML(c *gin.Context, xml string) {
	c.Writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
	c.Writer.Write([]byte(xml))
}

func IndexSitemap(c *gin.Context) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	base := `<sitemap>
      <loc>%s</loc>
      <lastmod>%s</lastmod>
   </sitemap>`

	lastEdition := models.LastEdition()
	lastCatTime := models.LastCategory().UpdatedAt.Format(time.RFC3339)
	lastEdTime := lastEdition.UpdatedAt.Format(time.RFC3339)
	lastEvTime := models.LastEvent().UpdatedAt.Format(time.RFC3339)

	xml += fmt.Sprintf(base, fmt.Sprintf("%s/sitemap-categories.xml", baseUrl), lastCatTime)
	xml += fmt.Sprintf(base, fmt.Sprintf("%s/sitemap-countries.xml", baseUrl), lastEdTime)
	xml += fmt.Sprintf(base, fmt.Sprintf("%s/sitemap-events.xml", baseUrl), lastEvTime)
	xml += fmt.Sprintf(base, fmt.Sprintf("%s/sitemap-editions.xml", baseUrl), lastEdTime)
	xml += fmt.Sprintf(base, fmt.Sprintf("%s/sitemap-pages.xml", baseUrl), staticDate)

	xml += `</sitemapindex>`

	doXML(c, xml)
}

func CategorySitemap(c *gin.Context) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	base := `<url>
<loc>https://confwatch.ninja/cats/%s</loc>
<lastmod>%s</lastmod>
<changefreq>daily</changefreq>
<priority>0.8</priority>
</url>`

	for _, cat := range models.Categories() {
		xml += fmt.Sprintf(base,
			cat.Slug,
			cat.UpdatedAt.Format(time.RFC3339))
	}

	xml += `</urlset>`

	doXML(c, xml)
}

func CountrySitemap(c *gin.Context) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	base := `<url>
<loc>https://confwatch.ninja/c/%s</loc>
<lastmod>%s</lastmod>
<changefreq>daily</changefreq>
<priority>0.8</priority>
</url>`

	lastEdition := models.LastEdition()
	lastEdTime := lastEdition.UpdatedAt.Format(time.RFC3339)

	for _, c := range models.Countries() {
		xml += fmt.Sprintf(base,
			c,
			lastEdTime)
	}

	xml += `</urlset>`

	doXML(c, xml)
}

func EventSitemap(c *gin.Context) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	base := `<url>
<loc>https://confwatch.ninja/events/%s</loc>
<lastmod>%s</lastmod>
<changefreq>daily</changefreq>
<priority>1.0</priority>
</url>`

	for _, e := range models.Events() {
		xml += fmt.Sprintf(base,
			e.Slug,
			e.UpdatedAt.Format(time.RFC3339))
	}

	xml += `</urlset>`

	doXML(c, xml)
}

func EditionSitemap(c *gin.Context) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	base := `<url>
<loc>https://confwatch.ninja/events/%s/editions/%s</loc>
<lastmod>%s</lastmod>
<changefreq>daily</changefreq>
<priority>1.0</priority>
</url>`

	for _, e := range models.Editions() {
		xml += fmt.Sprintf(base,
			e.Event().Slug,
			e.Slug,
			e.UpdatedAt.Format(time.RFC3339))
	}

	xml += `</urlset>`

	doXML(c, xml)
}

func PageSitemap(c *gin.Context) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	base := `<url>
<loc>https://confwatch.ninja/pages/%s</loc>
<lastmod>%s</lastmod>
<changefreq>weekly</changefreq>
<priority>1.0</priority>
</url>`

	xml += fmt.Sprintf(base,
		"about",
		staticDate)

	xml += `</urlset>`

	doXML(c, xml)
}
