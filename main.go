package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/ConfWatch/confwatchd/config"
	"github.com/ConfWatch/confwatchd/controllers"
	"github.com/ConfWatch/confwatchd/jobs"
	"github.com/ConfWatch/confwatchd/log"
	"github.com/ConfWatch/confwatchd/middlewares"
	"github.com/ConfWatch/confwatchd/models"

	"github.com/gin-gonic/gin"
	"github.com/michelloworld/ez-gin-template"
	"github.com/pariz/gountries"
)

var (
	signals    = make(chan os.Signal, 1)
	confFile   = ""
	debug      = false
	logfile    = ""
	noColors   = false
	seedFolder = ""
	router     = (*gin.Engine)(nil)
	cQuery     = gountries.New()
)

func init() {
	flag.StringVar(&confFile, "config", "config.json", "JSON configuration file.")
	flag.StringVar(&seedFolder, "seed", seedFolder, "Seed the database with the data from this folder.")
	flag.BoolVar(&debug, "log-debug", debug, "Enable debug logs.")
	flag.StringVar(&logfile, "log-file", logfile, "Log messages to this file instead of standard error.")
	flag.BoolVar(&noColors, "log-colors-off", noColors, "Disable colored output.")
}

func setupSignals() {
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	s := <-signals
	log.Raw("\n")
	log.Importantf("RECEIVED SIGNAL: %s", s)
	os.Exit(1)
}

func setupLogging() {
	var err error

	log.WithColors = !noColors

	if logfile != "" {
		log.Output, err = os.Create(logfile)
		if err != nil {
			log.Fatal(err)
		}

		defer log.Output.Close()
	}

	if debug == true {
		log.MinLevel = log.DEBUG
	} else {
		log.MinLevel = log.INFO
	}
}

func main() {
	flag.Parse()

	go setupSignals()

	setupLogging()

	if confFile != "" {
		if err := config.Load(confFile); err != nil {
			log.Fatal(err)
		}
	}

	if err := models.Setup(config.Conf.Database); err != nil {
		log.Fatal(err)
	}

	if seedFolder != "" {
		log.Infof("Seeding database from %s ...", log.Bold(seedFolder))
		if err := models.Seed(seedFolder); err != nil {
			log.Fatal(err)
		}
		return
	}

	if config.Conf.Dev {
		log.Infof("Running in dev mode.")
	} else {
		log.Infof("Running in prod mode.")
	}
	gin.SetMode(gin.ReleaseMode)

	render := eztemplate.New()
	render.TemplatesDir = "views/"
	render.Layout = "layouts/base"
	render.Ext = ".html"
	render.Debug = false

	render.TemplateFuncMap = template.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"now":   time.Now,
		"timeDiff": func(a time.Time, b time.Time) string {
			return b.Sub(a).String()
		},
		"CountByCountry": func(c string) int {
			return models.CountByCountry(c)
		},
		"toDate": func(t time.Time) string {
			return fmt.Sprintf("%02d/%02d/%d", t.Month(), t.Day(), t.Year())
		},
		"toDateLong": func(t time.Time) string {
			return fmt.Sprintf("%02d %s %d", t.Day(), t.Format("Jan"), t.Year())
		},
		"isPast": func(t time.Time) bool {
			return t.Before(time.Now())
		},
		"countryName": func(c string) string {
			cData, err := cQuery.FindCountryByAlpha(c)
			if err == nil {
				return cData.Name.Common
			}
			return c
		},
		"limit": func(s string, n int) string {
			len := len(s)
			if len > n {
				s = s[0:n-3] + "..."
			}
			return s
		},
	}

	router = gin.New()

	router.HTMLRender = render.Init()
	router.Use(middlewares.Security())
	router.Use(middlewares.ServeStatic("/", "static", "index.html"))
	// make sure connection is closed after every request is routed
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Connection", "close")
		c.Next()
	})

	router.GET("/sitemap-index.xml", controllers.IndexSitemap)
	router.GET("/sitemap-categories.xml", controllers.CategorySitemap)
	router.GET("/sitemap-countries.xml", controllers.CountrySitemap)
	router.GET("/sitemap-events.xml", controllers.EventSitemap)
	router.GET("/sitemap-editions.xml", controllers.EditionSitemap)
	router.GET("/sitemap-pages.xml", controllers.PageSitemap)

	router.GET("/", controllers.ShowHome)

	router.GET("/pages/about", controllers.AboutPage)

	router.GET("/cats/:cat_name", controllers.ShowCategory)
	router.GET("/c/:country_name", controllers.ShowCountry)

	router.GET("/events", controllers.ListEvents)
	router.GET("/events/:event_name", controllers.ShowEvent)
	router.GET("/events/:event_name/editions/:edition_name", controllers.ShowEdition)

	address := fmt.Sprintf("%s:%d", config.Conf.Address, config.Conf.Port)
	if address[0] == ':' {
		address = "0.0.0.0" + address
	}

	if config.Conf.Twitter.Enabled == true {
		if err := jobs.StartTwitterBot(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Infof("Twitter bot disabled.")
	}

	log.Infof("%s v%s is running on %s ...", config.APP_NAME, config.APP_VERSION, log.Bold(config.Conf.Hosts[0]))
	s := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.SetKeepAlivesEnabled(false)
	log.Fatal(s.ListenAndServe())
}
