///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package utils

import (
	"gopkg.in/yaml.v2"
)

func ParseYamlFromString(input string, target interface{}) error {
	return yaml.Unmarshal([]byte(input), target)
}

func Convert2Yaml(c interface{}) (string, error) {
	s, err := yaml.Marshal(c)
	return string(s), err
}
