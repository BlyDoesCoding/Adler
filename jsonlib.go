package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type JSONData struct {
	ID   int      `json:"id"`
	Data []string `json:"data"`
}

// Function to read data for a specific ID from a JSON file
func readJSONFromFile(filename string, id int) ([]string, error) {
	var jsonData []JSONData

	// Read the existing data from the file
	existingData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(existingData, &jsonData); err != nil {
		return nil, err
	}

	// Search for an entry with the given ID
	for _, entry := range jsonData {
		if entry.ID == id {
			return entry.Data, nil
		}
	}

	return nil, fmt.Errorf("ID %d not found in %s", id, filename)
}

// Function to append ".json" to a given filename
func appendJSONExtension(filename string) string {
	if !strings.HasSuffix(filename, ".json") {
		filename += ".json"
	}
	return filename
}

// Function to write JSON data to a file
func writeJSONToFile(filename string, id int, data []string) error {
	var jsonData []JSONData

	// Read the existing data from the file, if it exists
	existingData, err := os.ReadFile(filename)
	if err == nil {
		if err := json.Unmarshal(existingData, &jsonData); err != nil {
			return err
		}
	}

	// Search for an existing entry with the same ID and update its data
	updated := false
	for i, entry := range jsonData {
		if entry.ID == id {
			jsonData[i].Data = data
			updated = true
			break
		}
	}

	// If no existing entry was found, add a new one
	if !updated {
		jsonData = append(jsonData, JSONData{ID: id, Data: data})
	}

	// Marshal the updated data into JSON format
	jsonDataBytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}

	// Write the updated JSON data to the file
	err = os.WriteFile(filename, jsonDataBytes, 0644)
	return err
}

// Function to write or update JSON data to a file
func writeOrUpdateJSONToFile(filename string, id int, data []string) error {
	var jsonData []JSONData

	// Read the existing data from the file, if it exists
	existingData, err := os.ReadFile(filename)
	if err == nil {
		if err := json.Unmarshal(existingData, &jsonData); err != nil {
			return err
		}
	}

	// Search for an existing entry with the same ID and update its data
	updated := false
	for i, entry := range jsonData {
		if entry.ID == id {
			jsonData[i].Data = append(jsonData[i].Data, data...)
			updated = true
			break
		}
	}

	// If no existing entry was found, add a new one
	if !updated {
		jsonData = append(jsonData, JSONData{ID: id, Data: data})
	}

	// Marshal the updated data into JSON format
	jsonDataBytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}

	// Write the updated JSON data to the file
	err = os.WriteFile(filename, jsonDataBytes, 0644)
	return err
}

func readAllKeysAndData(filename string) (map[int][]string, error) {
	var jsonData []JSONData

	// Read the existing data from the file
	existingData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(existingData, &jsonData); err != nil {
		return nil, err
	}

	keyDataMap := make(map[int][]string)

	// Iterate through the entries and populate the map
	for _, entry := range jsonData {
		keyDataMap[entry.ID] = entry.Data
	}

	return keyDataMap, nil
}

func createEmptyJSONFile(filename string) error {
	// Create an empty map for the JSON structure
	data := make(map[string]interface{})

	// Create a new JSON file with the given filename
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the empty map as JSON and write it to the file
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	fmt.Printf("Empty JSON file '%s' created successfully.\n", filename)
	return nil
}

func getJSONFileNames() []string {
	// Get the current working directory
	dirPath, err := os.Getwd()
	if err != nil {
		return nil
	}

	jsonFileNames := []string{}

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") && info.Name() != "default.json" {
			// Remove the .json extension
			nameWithoutExtension := strings.TrimSuffix(info.Name(), ".json")
			jsonFileNames = append(jsonFileNames, nameWithoutExtension)
		}
		return nil
	})

	if err != nil {
		return nil
	}

	return jsonFileNames
}
