package eth

type MensaMenuResponse struct {
	ID      int    `json:"id"`
	Mensa   string `json:"mensa"`
	Daytime string `json:"daytime"`
	Hours   Hours  `json:"hours"`
	Menu    Menu   `json:"menu"`
}
type Opening struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
}
type Mealtime struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
}
type Hours struct {
	Opening  []Opening  `json:"opening"`
	Mealtime []Mealtime `json:"mealtime"`
}
type Mealtype struct {
	MealtypeID int    `json:"mealtype_id"`
	Label      string `json:"label"`
}
type Price struct {
	Student string `json:"student"`
	Staff   string `json:"staff"`
	Extern  string `json:"extern"`
}
type Allergen struct {
	AllergenID int    `json:"allergen_id"`
	Label      string `json:"label"`
}
type Origin struct {
	OriginID int    `json:"origin_id"`
	Label    string `json:"label"`
}
type Meal struct {
	ID          int        `json:"id"`
	Mealtypes   []Mealtype `json:"mealtypes"`
	Label       string     `json:"label"`
	Description []string   `json:"description"`
	Position    int        `json:"position"`
	Prices      Price      `json:"prices"`
	Allergens   []Allergen `json:"allergens"`
	Origins     []Origin   `json:"origins"`
}
type Menu struct {
	Date  string `json:"date"`
	Day   string `json:"day"`
	Meals []Meal `json:"meals"`
}

type FacilitiesListResponse struct {
	Locations []Location `json:"locations"`
	Facilites []Facility `json:"facilites"`
}
type Location struct {
	ID      int    `json:"id"`
	Label   string `json:"label"`
	LabelEn string `json:"label_en"`
}
type Facility struct {
	ID         int    `json:"id"`
	Label      string `json:"label"`
	LabelEn    string `json:"label_en"`
	Type       string `json:"type"`
	LocationID int    `json:"location_id"`
}
