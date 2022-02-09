package cmd

type config struct {
	// Groups
	EnabledProviders []string `json:"enabled_providers"`
	Groups           []group  `json:"groups"`
}

type group struct {
	Name string       `json:"name"`
	Refs []canteenRef `json:"refs"`
}

type canteenRef struct {
	Provider string `json:"provider"`
	Id       string `json:"id"`
}
