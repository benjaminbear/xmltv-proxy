package epgdownload

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/benjaminbear/xmltv-proxy/config"

	"github.com/benjaminbear/xmltv-proxy/epgdata"
	"github.com/benjaminbear/xmltv-proxy/today"
)

const (
	channelListURL = "https://live.tvspielfilm.de/static/content/channel-list/livetv"
	epgURL         = "https://live.tvspielfilm.de/static/broadcast/list"
	maxDays        = 14
)

type EPGDownloader struct {
	TimeToday     *today.Today
	Days          int
	ForceDownload bool
	InsecureTLS   bool
	client        *http.Client
	channels      []*epgdata.Channel
}

func (e *EPGDownloader) DownloadEPG() error {
	err := os.MkdirAll(config.FolderEPGData, os.ModePerm)
	if err != nil {
		return err
	}

	err = e.removeDeprecated()
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: e.InsecureTLS},
	}

	e.client = &http.Client{
		Transport: tr,
	}

	fmt.Println("Downloading channel list...")

	err = e.downloadChannelList()
	if err != nil {
		return err
	}

	fmt.Printf("Trying to download epg for %d days\n", e.Days)
	for i := 0; i < e.Days && i < maxDays; i++ {
		err = e.downloadEPG(e.TimeToday.GetDayPlus(i))
		if err != nil {
			return err
		}

		fmt.Printf("Successfully downloaded epg for day %s\n", e.TimeToday.GetDayPlus(i))
	}

	return nil
}

func (e *EPGDownloader) downloadEPG(date string) error {
	err := os.MkdirAll(filepath.Join(config.FolderEPGData, date), os.ModePerm)
	if err != nil {
		return err
	}

	for _, channel := range e.channels {
		if fileExists(date, channel.ID) {
			if !e.ForceDownload || today.New().GetString() != date {
				continue
			}

			err := os.Remove(filepath.Join(config.FolderEPGData, date, channel.ID+".json"))
			if err != nil {
				return err
			}
		}

		resp, err := http.Get(fmt.Sprintf("%s/%s/%s", epgURL, channel.ID, date))
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Printf("request %s returned response %v\n", resp.Request.URL, resp.StatusCode)
			continue
		}

		// Create an empty file
		file, err := os.Create(filepath.Join(config.FolderEPGData, date, channel.ID+".json"))
		if err != nil {
			return err
		}

		defer file.Close()

		// Write the bytes to the file
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *EPGDownloader) downloadChannelList() error {
	resp, err := http.Get(channelListURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("request returned response %v", resp.StatusCode)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &e.channels)
	if err != nil {
		return err
	}

	err = epgdata.WriteChannelsFile(e.channels)

	return err
}

func (e *EPGDownloader) removeDeprecated() error {
	folderContent, err := ioutil.ReadDir(config.FolderEPGData)
	if err != nil {
		return err
	}

	for _, info := range folderContent {
		fileDate, err := today.Parse(info.Name())
		if err != nil {
			continue
		}

		if fileDate.Before(e.TimeToday.Time) {
			err = os.RemoveAll(filepath.Join(config.FolderEPGData, info.Name()))
			if err != nil {
				return err
			}

			fmt.Println("Successfully removed epg for day", info.Name())
		}
	}

	return nil
}

func fileExists(date string, channel string) bool {
	_, err := os.Stat(filepath.Join(config.FolderEPGData, date, channel+".json"))

	return !errors.Is(err, os.ErrNotExist)
}
