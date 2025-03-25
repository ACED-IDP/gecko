package config

type FieldConfig struct {
	Field     string `json:"field"`
	DataField string `json:"dataField"`
	Index     string `json:"index"`
	Label     string `json:"label"`
	Type      string `json:"type"`
}

type FilterTab struct {
	Title        string                 `json:"title"`
	Fields       []string               `json:"fields"`
	FieldsConfig map[string]FieldConfig `json:"fieldsConfig"`
}

type FiltersConfig struct {
	Tabs []FilterTab `json:"tabs"`
}

type TableConfig struct {
	Enabled bool                          `json:"enabled"`
	Fields  []string                      `json:"fields"`
	Columns map[string]TableColumnsConfig `json:"columns"`
}

type TableColumnsConfig struct {
	Field string `json:"field"`
	Title string `json:"title"`
}

type GuppyConfig struct {
	DataType       string   `json:"dataType"`
	NodeCountTitle string   `json:"nodeCountTitle"`
	FieldMapping   []string `json:"fieldMapping"`
}

type Chart struct {
	ChartType string `json:"chartType"`
	Title     string `json:"title"`
}

type ConfigItem struct {
	TabTitle         string           `json:"tabTitle"`
	GuppyConfig      GuppyConfig      `json:"guppyConfig"`
	Charts           map[string]Chart `json:"charts"`
	Filters          FiltersConfig    `json:"filters"` // Updated type
	Table            TableConfig      `json:"table"`   // Updated type
	Dropdowns        map[string]any   `json:"dropdowns"`
	Buttons          []any            `json:"buttons"`
	LoginForDownload bool             `json:"loginForDownload"`
}
