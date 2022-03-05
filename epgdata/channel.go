package epgdata

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/benjaminbear/xmltv-proxy/config"
)

type Channel struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	NameShort  string       `json:"name_short"`
	IsPremium  bool         `json:"is_premium"`
	ImageSmall ChannelImage `json:"image_small"`
	ImageLarge ChannelImage `json:"image_large"`
}

type ChannelImage struct {
	Hash string `json:"hash"`
	URL  string `json:"url"`
}

func ReadChannelsFile() ([]*Channel, error) {
	data, err := ioutil.ReadFile(filepath.Join(config.FolderEPGData, fileEPGIncludeChannels))
	if err != nil {
		return nil, err
	}

	var channels []*Channel

	err = json.Unmarshal(data, &channels)
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func WriteChannelsFile(channels []*Channel) error {
	bytes, err := json.Marshal(channels)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(config.FolderEPGData, fileEPGIncludeChannels), bytes, 0644)
}
