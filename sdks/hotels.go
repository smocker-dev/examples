package sdks

import (
	"context"
	"encoding/json"
	"net/http"

	"example/types"
)

type Hotels interface {
	GetHotels(ctx context.Context) ([]types.Hotel, error)
	GetHotelByName(ctx context.Context, name string) (types.Hotel, error)
}

type hotelsClient struct {
	url string
}

func NewHotelsClient(url string) Hotels {
	return &hotelsClient{url: url}
}

func (c *hotelsClient) GetHotels(ctx context.Context) ([]types.Hotel, error) {
	var hotels []types.Hotel
	req, err := http.NewRequestWithContext(ctx, "GET", c.url+"/hotels", nil)
	if err != nil {
		return hotels, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return hotels, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&hotels)
	return hotels, err
}

func (c *hotelsClient) GetHotelByName(ctx context.Context, name string) (types.Hotel, error) {
	var hotel types.Hotel
	req, err := http.NewRequestWithContext(ctx, "GET", c.url+"/hotels/"+name, nil)
	if err != nil {
		return hotel, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return hotel, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&hotel)
	return hotel, err
}
