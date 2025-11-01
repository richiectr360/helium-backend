package main

// Helper function to get component keys
func getComponentKeys() []string {
	keys := make([]string, 0, len(componentTemplates))
	for key := range componentTemplates {
		keys = append(keys, key)
	}
	return keys
}

// Helper function to get language keys
func getLanguageKeys() []string {
	keys := make([]string, 0, len(localizationDB))
	for key := range localizationDB {
		keys = append(keys, key)
	}
	return keys
}

