package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

func Get() map[string]interface{} {
	_, filename, _, _ := runtime.Caller(0)
	configFilePath := filepath.Join(filepath.Dir(filename), "../../../config/wrapper_config.yaml")

	yamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}

	var config map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return config
}
