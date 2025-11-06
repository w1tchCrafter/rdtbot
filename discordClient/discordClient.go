package discordclient

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type DiscordClient struct {
	token    string
	channels []string
}

func New(token string, channels ...string) *DiscordClient {
	return &DiscordClient{
		token:    token,
		channels: channels,
	}
}

// note: change the inner loop to use goroutines later
func (dc *DiscordClient) SendMsg(content string) error {
	errs := make([]error, 0)

	for _, i := range dc.channels {
		api := fmt.Sprintf("https://discord.com/api/v9/channels/%v/messages", i)
		body := strings.NewReader(fmt.Sprintf(`{"content" : "%v"}`, content))
		req, err := http.NewRequest(http.MethodPost, api, body)

		if err != nil {
			return err
		}

		req.Header.Set("authorization", dc.token)
		req.Header.Set("content-type", "application/json")
		client := &http.Client{}

		_, err = client.Do(req)
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
