package eth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lukasmoellerch/mensa-cli/internal/base"
	"github.com/mailru/easyjson"
	"golang.org/x/sync/errgroup"
)

var _ base.Provider = (*Provider)(nil)

type Provider struct{}

func (p *Provider) Id() string {
	return "eth"
}

func (p *Provider) Label() string {
	return "ETH"
}

func (p *Provider) FetchMenus(ctx context.Context, canteens []base.CanteenRef, date string, daytime string, lang string) ([]base.CanteenMenu, error) {
	ids := make([]int, len(canteens))
	for i, canteen := range canteens {
		id, err := strconv.Atoi(canteen.ID)
		if err != nil {
			return []base.CanteenMenu{}, err
		}
		ids[i] = id
	}

	eg, ctx := errgroup.WithContext(ctx)
	menus := make([]base.CanteenMenu, len(ids))
	for i, canteenID := range ids {
		i, canteenID := i, canteenID
		eg.Go(func() error {
			menu, err := fetchMenu(ctx, canteenID, lang, date, daytime)
			if err != nil {
				return err
			}
			menus[i].Canteen = strconv.Itoa(menu.ID)
			meals := make([]base.Meal, len(menu.Menu.Meals))
			for j, meal := range menu.Menu.Meals {
				meals[j].Label = meal.Label
				meals[j].Description = meal.Description
				studentPrice, err := base.ParsePrice(meal.Prices.Student)
				if err != nil {
					return err
				}
				meals[j].Prices.Student = studentPrice
				staffPrice, err := base.ParsePrice(meal.Prices.Staff)
				if err != nil {
					return err
				}
				meals[j].Prices.Staff = staffPrice
				externPrice, err := base.ParsePrice(meal.Prices.Extern)
				if err != nil {
					return err
				}
				meals[j].Prices.Extern = externPrice
			}
			menus[i].Meals = meals
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return []base.CanteenMenu{}, err
	}
	return menus, nil
}

func (p *Provider) FetchCanteens(ctx context.Context, lang string) ([]base.CanteenMetadata, error) {
	facilities, err := fetchFacilities(ctx)
	if err != nil {
		return []base.CanteenMetadata{}, err
	}
	metadatas := make([]base.CanteenMetadata, len(facilities.Facilites))
	for i, facility := range facilities.Facilites {
		metadatas[i].ID = strconv.Itoa(facility.ID)
		if lang == "de" {
			metadatas[i].Label = facility.Label
		} else {
			metadatas[i].Label = facility.LabelEn
		}
	}
	return metadatas, nil
}

func fetchMenu(ctx context.Context, canteenID int, lang string, date string, daytime string) (MensaMenuResponse, error) {
	url := fmt.Sprintf("https://www.webservices.ethz.ch/gastro/v1/RVRI/Q1E1/mensas/%d/%s/menus/daily/%s/%s", canteenID, lang, date, daytime)
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
	err = easyjson.UnmarshalFromReader(resp.Body, &menu)
	if err != nil {
		return MensaMenuResponse{}, err
	}
	return menu, nil
}

func fetchFacilities(ctx context.Context) (FacilitiesListResponse, error) {
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
	err = easyjson.UnmarshalFromReader(resp.Body, &facilities)
	if err != nil {
		return FacilitiesListResponse{}, err
	}
	return facilities, nil
}
