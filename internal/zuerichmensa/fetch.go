package zuerichmensa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchMenuEth(ctx context.Context, mensaID int, lang string, date string, daytime string) (MensaMenuResponse, error) {
	url := fmt.Sprintf("https://www.webservices.ethz.ch/gastro/v1/RVRI/Q1E1/mensas/%d/%s/menus/daily/%s/%s", mensaID, lang, date, daytime)
	// Make http request and decode response as json
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return MensaMenuResponse{}, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return MensaMenuResponse{}, err
	}
	defer resp.Body.Close()
	var menu MensaMenuResponse
	err = json.NewDecoder(resp.Body).Decode(&menu)
	if err != nil {
		return MensaMenuResponse{}, err
	}
	return menu, nil
}

func FetchFacilitiesEth(ctx context.Context) (FacilitiesListResponse, error) {
	url := "https://www.webservices.ethz.ch/gastro/v1/RVRI/Q1E1/facilities"
	// Make http request and decode response as json
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return FacilitiesListResponse{}, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return FacilitiesListResponse{}, err
	}
	defer resp.Body.Close()
	var facilities FacilitiesListResponse
	err = json.NewDecoder(resp.Body).Decode(&facilities)
	if err != nil {
		return FacilitiesListResponse{}, err
	}
	return facilities, nil
}
