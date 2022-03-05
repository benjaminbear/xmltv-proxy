package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/benjaminbear/xmltv-proxy/config"
	"github.com/benjaminbear/xmltv-proxy/epgdata"
	"github.com/benjaminbear/xmltv-proxy/epgdownload"
	"github.com/benjaminbear/xmltv-proxy/today"
	"github.com/benjaminbear/xmltv-proxy/xmltv"
)

type RunTime struct {
	Channels  []*xmltv.Channel
	XMLTv     *xmltv.Tv
	EPG       []*epgdata.Show
	Config    *config.Config
	initiated bool
}

func (r *RunTime) EPGCron() error {
	fmt.Println("Starting cron ... ")
	timeToday := today.New()

	// download days
	epgDownloader := &epgdownload.EPGDownloader{
		TimeToday:   timeToday,
		Days:        r.Config.Days,
		InsecureTLS: r.Config.InsecureTLS,
	}

	if r.Config.DailyDownload && r.initiated {
		epgDownloader.ForceDownload = true
	}

	if !r.initiated {
		r.initiated = true
	}

	err := epgDownloader.DownloadEPG()
	if err != nil {
		return err
	}

	err = r.parseChannels()
	if err != nil {
		return err
	}

	// parse downloaded days
	for i := 0; i < r.Config.Days-1; i++ {
		showDay, err := epgdata.ReadEPGDay(timeToday.GetDayPlus(i))
		if err != nil {
			return err
		}

		r.EPG = append(r.EPG, showDay...)
		if err != nil {
			return err
		}
	}

	err = r.Merge()
	if err != nil {
		return err
	}

	fmt.Println("Finished cron!")

	return nil
}

func (r *RunTime) parseChannels() error {
	channels, err := epgdata.ReadChannelsFile()
	if err != nil {
		return err
	}

	for _, channel := range channels {
		r.Channels = append(r.Channels, &xmltv.Channel{
			Id: channel.ID,
			DisplayName: &xmltv.StdLangElement{
				Name:     channel.Name,
				Language: "de",
			},
			Icon: &xmltv.Icon{
				Source: channel.ImageLarge.URL,
			},
		})
	}

	return nil
}

func (r *RunTime) Merge() (err error) {
	// Create XMLTV File
	r.XMLTv = xmltv.NewXMLTVFile()
	r.XMLTv.Channels = r.Channels
	r.XMLTv.GeneratorInfoName = "xmltv-proxy"
	r.XMLTv.GeneratorInfoURL = "https://github.com/benjaminbear/xmltv-proxy"

	for _, show := range r.EPG {
		tvProgram := &xmltv.Program{
			Channel: show.BroadcasterID,
			Start:   show.Timestart.In(r.Config.TimeZone).Format(xmltv.DateTimeFormat),
			Stop:    show.Timeend.In(r.Config.TimeZone).Format(xmltv.DateTimeFormat),
		}

		if show.Title != "" {
			tvProgram.Title = &xmltv.StdLangElement{
				Name:     show.Title,
				Language: "de",
			}
		}

		if show.EpisodeTitle != "" {
			tvProgram.SubTitle = &xmltv.StdLangElement{
				Name:     show.EpisodeTitle,
				Language: "de",
			}
		} else if show.CurrentTopics != "" {
			tvProgram.SubTitle = &xmltv.StdLangElement{
				Name:     show.CurrentTopics,
				Language: "de",
			}
		} else if show.Subline != "" {
			tvProgram.SubTitle = &xmltv.StdLangElement{
				Name:     show.Subline,
				Language: "de",
			}
		}

		if show.Text != "" {
			tvProgram.Desc = &xmltv.StdLangElement{
				Name:     show.Text,
				Language: "de",
			}
		}

		if show.EpisodeNumber != "" || show.SeasonNumber != "" {
			seasonNum := 1
			episodeNum := "0"

			if show.EpisodeNumber != "" {
				episodeNum = show.EpisodeNumber
			}

			if show.SeasonNumber != "" {
				seasonNum, err = strconv.Atoi(show.SeasonNumber)
				if err != nil {
					return err
				}
			}

			tvProgram.EpisodeNum = &xmltv.EpisodeNum{
				Value:  fmt.Sprintf("%d . %s . ", seasonNum-1, episodeNum),
				System: "xmltv_ns",
			}
		}

		if show.Images.Size476 != nil {
			tvProgram.Icon = &xmltv.Icon{
				Source: show.Images.Size476.Source,
				Width:  show.Images.Size476.Width,
				Height: show.Images.Size476.Height,
			}
		}

		if show.Year != 0 {
			tvProgram.Date = strconv.Itoa(show.Year)
		}

		if show.Country != "" {
			tvProgram.Country = &xmltv.StdLangElement{
				Name: show.Country,
			}
		}

		if show.LengthNetAndGross != "" {
			splits := strings.Split(show.LengthNetAndGross, "/")

			if len(splits) == 2 {
				tvProgram.Length = &xmltv.Length{
					Value: splits[1],
					Unit:  xmltv.Minutes,
				}
			}
		}

		if show.Fsk != 0 {
			tvProgram.Rating = append(tvProgram.Rating, &xmltv.Rating{
				Value:  strconv.Itoa(show.Fsk),
				System: "FSK",
			})
		}

		if show.Genre != "" {
			tvProgram.Categories = append(tvProgram.Categories, &xmltv.StdLangElement{
				Name:     show.Genre,
				Language: "de",
			})
		}

		if show.IsNew {
			tvProgram.Premiere = &xmltv.StdLangElement{
				Name:     "Neue Sendung/Folgen!",
				Language: "de",
			}
		}

		if show.HasCredits() {
			tvProgram.Credits = &xmltv.Credits{}

			if show.Anchorman != nil {
				if len(show.Anchorman.List) > 0 {
					tvProgram.Credits.Presenters = show.Anchorman.List
				}
			}

			if show.StudioGuests != nil {
				if len(show.StudioGuests.List) > 0 {
					tvProgram.Credits.Guests = show.StudioGuests.List
				}
			}

			if show.Director != nil {
				if len(show.Director.List) > 0 {
					tvProgram.Credits.Directors = show.Director.List
				}
			}

			if len(show.Actors) > 0 {
				for _, actor := range show.Actors {
					a := &xmltv.Actor{
						Name: actor.Actor,
					}

					if actor.Role != "" {
						a.Role = actor.Role
					}

					tvProgram.Credits.Actors = append(tvProgram.Credits.Actors, a)
				}
			}
		}

		r.XMLTv.Programs = append(r.XMLTv.Programs, tvProgram)
	}

	return nil
}

func (r *RunTime) XMLTvServer(w http.ResponseWriter, req *http.Request) {
	err := r.XMLTv.WriteFile("latest.xml")
	if err != nil {
		fmt.Println(err)
	}

	http.ServeFile(w, req, "latest.xml")

	IPAddress := req.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = req.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = req.RemoteAddr
	}

	fmt.Println("Request from: ", IPAddress)
}
