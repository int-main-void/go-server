/*
Package config has a method to read a json config file into key/value map
 */
package config

import (
	"encoding/json"
	"log"
	"os"
)

const stageKey="STAGE"
const configFilenameKey="CONFIG_FILENAME"

func readConfig(configFileName string, configName string, version string, stage string) (map[string]string, error) {

	config := map[string]string{}

	configFile, err := os.Open(configFileName)
	if(err != nil) { return config, err }

	//  "configName"-"version"-"stage"-"key"-"value"
	allConfigs := map[string]map[string]map[string]map[string]string{}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&allConfigs)
	if(err != nil) { return config, err }

	config = allConfigs[configName][version][stage]

	return config, nil
}

/*
SetupConfig reads a json config file into a key/value map.
*/
func SetupConfig(configName string, runtimeVersion string) (map[string]string, error) {
	runtimeStage := os.Getenv(stageKey)
	if(runtimeStage != "dev" && runtimeStage != "integration" && runtimeStage != "staging" && runtimeStage != "live") {
		log.Println("invalid runtime stage: ", runtimeStage)
		os.Exit(1)
	}

	configFileName := os.Getenv(configFilenameKey)
	config, err := readConfig(configFileName, configName, runtimeVersion, runtimeStage)
	if(err != nil) {
		log.Println(err)
		os.Exit(1)
	}
	return config, err
}
