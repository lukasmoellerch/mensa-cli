package base

import "context"

type MealPrices struct {
	Student int64
	Staff   int64
	Extern  int64
}

type Meal struct {
	Label       string
	Description []string
	Prices      MealPrices
}

type CanteenMenu struct {
	Canteen string
	Meals   []Meal
}

type CanteenMetadata struct {
	// A unique identifier for the canteen in the given provider
	ID string

	// The label of the canteen.
	Label string

	// Meta is a map of arbitrary key-value pairs that can be used to store
	// additional information about the canteen.
	Meta map[string]string
}

type CanteenRef struct {
	// Meta is a map of arbitrary key-value pairs that can be used to store
	// additional information about the canteen.
	Meta map[string]string

	// The ID of the canteen
	ID string
}

type Provider interface {
	// A unique identifier for the provider.
	Id() string

	// A human readable name for the provider.
	Label() string

	// List all menus of the given canteens for the given date and daytime.
	FetchMenus(ctx context.Context, canteens []CanteenRef, date string, daytime string, lang string) ([]CanteenMenu, error)

	// List all canteens fo this provider
	FetchCanteens(ctx context.Context, lang string) ([]CanteenMetadata, error)
}
