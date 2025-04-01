package config

type FieldConfig struct {
	Field     string `json:"field,omitempty"`
	DataField string `json:"dataField,omitempty"`
	Index     string `json:"index,omitempty"`
	Label     string `json:"label"`
	Type      string `json:"type,omitempty"`
}

type FilterTab struct {
	Title        string                 `json:"title,omitempty"`
	Fields       []string               `json:"fields"`
	FieldsConfig map[string]FieldConfig `json:"fieldsConfig,omitempty"`
}

type FiltersConfig struct {
	Tabs []FilterTab `json:"tabs"`
}

type TableConfig struct {
	Enabled       bool                          `json:"enabled"`
	Fields        []string                      `json:"fields"`
	Columns       map[string]TableColumnsConfig `json:"columns,omitempty"`
	DetailsConfig TableDetailsConfig            `json:"detailsConfig,omitempty"`
}

type TableColumnsConfig struct {
	Field string `json:"field"`
	Title string `json:"title"`
}

type TableDetailsConfig struct {
	Panel       string            `json:"panel,omitempty"`
	Mode        string            `json:"mode,omitempty"`
	IDField     string            `json:"idField,omitempty"`
	FilterField string            `json:"filterField,omitempty"`
	Title       string            `json:"title,omitempty"`
	NodeType    string            `json:"nodeType,omitempty"`
	NodeFields  map[string]string `json:"nodeFields,omitempty"`
}

type GuppyConfig struct {
	DataType                  string              `json:"dataType"`
	NodeCountTitle            string              `json:"nodeCountTitle"`
	FieldMapping              []GuppyFieldMapping `json:"fieldMapping,omitempty"`
	AccessibleFieldCheckList  []string            `json:"accessibleFieldCheckList,omitempty"`
	AccessibleValidationField string              `json:"accessibleValidationField,omitempty"`
	ManifestMapping           ManifestMapping     `json:"manifestMapping,omitempty"`
}

type GuppyFieldMapping struct {
	Field string `json:"field,omitempty"`
	Name  string `json:"name,omitempty"`
}

type ManifestMapping struct {
	ResourceIndexType               string `json:"resourceIndexType,omitempty"`
	ResourceIdField                 string `json:"resourceIdField,omitempty"`
	ReferenceIdFieldInResourceIndex string `json:"referenceIdFieldInResourceIndex,omitempty"`
	ReferenceIdFieldInDataIndex     string `json:"referenceIdFieldInDataIndex,omitempty"`
}

type Chart struct {
	ChartType string `json:"chartType"`
	Title     string `json:"title"`
}

type ButtonConfig struct {
	Enabled    bool             `json:"enabled,omitempty"`
	Type       string           `json:"type,omitempty"`
	Action     string           `json:"action,omitempty"`
	Title      string           `json:"title,omitempty"`
	LeftIcon   string           `json:"leftIcon,omitempty"`
	RightIcon  string           `json:"rightIcon,omitempty"`
	FileName   string           `json:"fileName,omitempty"`
	ActionArgs ButtonActionArgs `json:"actionArgs,omitempty"`
}

type ButtonActionArgs struct {
	ResourceIndexType               string   `json:"resourceIndexType,omitempty"`
	ResourceIdField                 string   `json:"resourceIdField,omitempty"`
	ReferenceIdFieldInDataIndex     string   `json:"referenceIdFieldInDataIndex,omitempty"`
	ReferenceIdFieldInResourceIndex string   `json:"referenceIdFieldInResourceIndex,omitempty"`
	FileFields                      []string `json:"fileFields,omitempty"`
}

type ConfigItem struct {
	TabTitle         string           `json:"tabTitle"`
	GuppyConfig      GuppyConfig      `json:"guppyConfig"`
	Charts           map[string]Chart `json:"charts,omitempty"`
	Filters          FiltersConfig    `json:"filters"`
	Table            TableConfig      `json:"table"`
	Dropdowns        map[string]any   `json:"dropdowns,omitempty"`
	Buttons          []ButtonConfig   `json:"buttons,omitempty"`
	LoginForDownload bool             `json:"loginForDownload,omitempty"`
}
