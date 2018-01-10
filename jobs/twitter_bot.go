package jobs

import (
	"fmt"
	"time"

	"github.com/ConfWatch/confwatchd/config"
	"github.com/ConfWatch/confwatchd/log"
	"github.com/ConfWatch/confwatchd/models"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/dustin/go-humanize"
)

func StartTwitterBot() error {
	c := oauth1.NewConfig(config.Conf.Twitter.ConsumerKey, config.Conf.Twitter.ConsumerSecret)
	token := oauth1.NewToken(config.Conf.Twitter.AccessToken, config.Conf.Twitter.AccessSecret)
	httpClient := c.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	feed, _, err := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{Count: 1})
	if err != nil {
		return err
	}

	go func() {
		log.Infof("Twitter bot started, last tweet: %s", log.Dim(feed[0].Text))

		for {
			log.Debugf("Twitter bot loop.")
			found, edition := models.FirstEditionToShare(1)
			if found == true {
				log.Infof("Sharing %s ...", edition.Title)

				now := time.Now()
				msg := ""
				toStart := edition.Starts.Sub(now)
				will := edition.Starts.After(now)
				t := humanize.RelTime(edition.Starts, now, "ago", "from now")

				// Never been shared before.
				if edition.SharedAt.IsZero() == true {
					msg = fmt.Sprintf("Hey, check this out, %s has been added to the database, it will start in %s! %s",
						edition.Title,
						t,
						fmt.Sprintf("https://confwatch.ninja/events/%s/editions/%s",
							edition.Event().Slug,
							edition.Slug,
						),
					)
				} else if will && toStart.Hours() < 24 {
					msg = fmt.Sprintf("Yo, %s is starting today! %s",
						edition.Title,
						fmt.Sprintf("https://confwatch.ninja/events/%s/editions/%s",
							edition.Event().Slug,
							edition.Slug,
						),
					)
				} else if will && toStart.Hours() < 48 {
					msg = fmt.Sprintf("Hey, %s is starting tomorrow :D %s",
						edition.Title,
						fmt.Sprintf("https://confwatch.ninja/events/%s/editions/%s",
							edition.Event().Slug,
							edition.Slug,
						),
					)
				}

				if msg != "" {
					log.Infof("Tweeting: %s", log.Dim(msg))
					_, _, err := client.Statuses.Update(msg, nil)
					if err != nil {
						log.Errorf("Error while sending tweet: %s", err)
					}
				}

				edition.SharedAt = time.Now()
				models.Save(&edition)
			}

			time.Sleep(time.Hour * 1)
		}
	}()

	return nil
}
