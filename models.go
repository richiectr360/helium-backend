package main

// ComponentTemplate represents a React component template
type ComponentTemplate struct {
	ComponentName string   `json:"component_name"`
	ComponentType string   `json:"component_type"`
	Template      string   `json:"template"`
	RequiredKeys  []string `json:"required_keys"`
}

// ComponentMetadata represents metadata for a component
type ComponentMetadata struct {
	ComponentID  string   `json:"component_id"`
	LastUpdated  string   `json:"last_updated"`
	RequiredKeys []string `json:"required_keys"`
}

// LocalizedComponent represents the response structure
type LocalizedComponent struct {
	ComponentName string            `json:"component_name"`
	ComponentType string            `json:"component_type"`
	Language      string            `json:"language"`
	Template      string            `json:"template"`
	LocalizedData map[string]string `json:"localized_data"`
	Metadata      ComponentMetadata `json:"metadata"`
	Cached        bool              `json:"cached,omitempty"`
}

