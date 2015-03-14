package main

import (
	"fmt"
	"strings"
)

func ParseEnviron(envs []string) map[string]string {
	envsMap := make(map[string]string)
	for _, env := range envs {
		parts := strings.SplitN(env, "=", 2)
		envsMap[parts[0]] = parts[1]
	}
	return envsMap
}

func CreateEnviron(envsMap map[string]string) []string {
	envs := make([]string, 0, len(envsMap))
	for key, val := range envsMap {
		envs = append(envs, fmt.Sprintf("%s=%s", key, val))
	}
	return envs
}
