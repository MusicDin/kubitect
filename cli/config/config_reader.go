package config

import (
	"cli/config/modelconfig"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)


func ReadConfig(path string) modelconfig.Config {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	config := modelconfig.Config{}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return config
}

func ConfigToString(config *modelconfig.Config) (string, error) {
	s, err := yaml.Marshal(config)
	return string(s), err
}
