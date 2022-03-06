package epgdata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/benjaminbear/xmltv-proxy/config"
)

const (
	folderPersistence      = "persistence"
	filePersistence        = "epgdata.bin"
	fileEPGIncludeChannels = "channels.json"
)

var (
	pathPersistence = filepath.Join(folderPersistence, filePersistence)
)

type ShowDay struct {
	Date  string
	Shows []*Show
}

type Show struct {
	ID                string          `json:"id"`
	HexID             string          `json:"hexId"`
	AssetID           string          `json:"assetId"`
	SharingID         string          `json:"sharingId"`
	TrackingID        string          `json:"trackingId"`
	BroadcasterID     string          `json:"broadcasterId"`
	BroadcasterName   string          `json:"broadcasterName"`
	Country           string          `json:"country,omitempty"`
	CurrentTopics     string          `json:"currentTopics,omitempty"`
	EpisodeTitle      string          `json:"episodeTitle,omitempty"`
	Genre             string          `json:"genre"`
	Images            Images          `json:"images"`
	IsHDTV            bool            `json:"isHDTV"`
	IsLive            bool            `json:"isLive"`
	IsNew             bool            `json:"isNew"`
	IsPayTv           bool            `json:"isPayTv"`
	IsStereo          bool            `json:"isStereo"`
	IsTipOfTheDay     bool            `json:"isTipOfTheDay"`
	IsTopTip          bool            `json:"isTopTip"`
	LengthNetAndGross string          `json:"lengthNetAndGross"`
	RepeatHint        string          `json:"repeatHint,omitempty"`
	Subline           string          `json:"subline,omitempty"`
	Text              string          `json:"text,omitempty"`
	IsOriginalText    bool            `json:"isOriginalText"`
	Timeend           dateTime        `json:"timeend"`
	Timestart         dateTime        `json:"timestart"`
	Title             string          `json:"title"`
	OriginalTitle     string          `json:"originalTitle,omitempty"`
	SartID            string          `json:"sart_id"`
	Year              int             `json:"year"`
	Actors            []*Actor        `json:"actors,omitempty"`
	Anchorman         *StdListElement `json:"anchorman,omitempty"`
	Director          *StdListElement `json:"director,omitempty"`
	StudioGuests      *StdListElement `json:"studio_guests,omitempty"`
	EpisodeNumber     seriesNumber    `json:"episodeNumber,omitempty"`
	SeasonNumber      seriesNumber    `json:"seasonNumber,omitempty"`
	Videos            []struct {
		Title      string `json:"title"`
		Catchline  string `json:"catchline"`
		BlockAds   int    `json:"blockAds"`
		StillImage string `json:"stillImage"`
		Video      []struct {
			Type     string `json:"type"`
			URL      string `json:"url"`
			Duration int    `json:"duration"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
		} `json:"video"`
	} `json:"videos,omitempty"`
	ThumbID        string `json:"thumbId,omitempty"`
	ThumbIDNumeric int    `json:"thumbIdNumeric,omitempty"`
	Fsk            int    `json:"fsk,omitempty"`
	Preview        string `json:"preview,omitempty"`
}

func (s *Show) HasCredits() bool {
	if len(s.Actors) > 0 {
		return true
	}
	if s.StudioGuests != nil {
		if len(s.StudioGuests.List) > 0 {
			return true
		}
	}

	if s.Anchorman != nil {
		if len(s.Anchorman.List) > 0 {
			return true
		}
	}

	if s.Director != nil {
		if len(s.Director.List) > 0 {
			return true
		}
	}

	return false
}

func NewEPGData() []*Show {
	epgData := make([]*Show, 0)

	return epgData
}

func ReadEPGFile(path string) ([]*Show, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	epgData := NewEPGData()

	err = json.Unmarshal(data, &epgData)
	if err != nil {
		return nil, err
	}

	return epgData, nil
}

func ReadEPGDay(date string) ([]*Show, error) {
	showDay := make([]*Show, 0)

	fileInfos, err := ioutil.ReadDir(filepath.Join(config.FolderEPGData, date))
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range fileInfos {
		shows, err := ReadEPGFile(filepath.Join(config.FolderEPGData, date, fileInfo.Name()))
		if err != nil {
			return nil, err
		}

		showDay = append(showDay, shows...)
	}

	return showDay, nil
}

func Save(data interface{}) (err error) {
	if _, err = os.Stat(pathPersistence); !os.IsNotExist(err) {
		err = os.Remove(pathPersistence)
		if err != nil {
			return err
		}
	}

	if _, err = os.Stat(folderPersistence); os.IsNotExist(err) {
		err = os.MkdirAll(folderPersistence, os.ModePerm)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(pathPersistence)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	_, err = io.Copy(f, bytes.NewReader(b))

	return err
}

func Load(showDays []*ShowDay) (err error) {
	if _, err = os.Stat(pathPersistence); os.IsNotExist(err) {
		return nil
	}

	fmt.Println("Loading persistent data from disk")

	bytes, err := ioutil.ReadFile(pathPersistence)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &showDays)

	return err
}
