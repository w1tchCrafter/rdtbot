package redditclient

import "encoding/json"

type data struct {
	MediaUrl string `json:"url_overridden_by_dest"`
}

type children struct {
	Data data `json:"data"`
}

type listingData struct {
	Children []children `json:"children"`
}

type Listing struct {
	Data listingData `json:"data"`
}

func NewListing(apiResponse []byte) (*Listing, error) {
	listing := &Listing{}
	err := json.Unmarshal(apiResponse, listing)

	return listing, err
}

func (l *Listing) GetLinks() []string {
	links := make([]string, 0)

	for _, i := range l.Data.Children {
		links = append(links, i.Data.MediaUrl)
	}

	return links
}
