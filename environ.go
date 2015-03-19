package main

import (
	"fmt"
	"strings"
)

func ParseVar(envVar string) (string, string) {
	parts := make([]string, 2)
	parts = append(parts[:0], strings.SplitN(envVar, "=", 2)...)
	return parts[0], parts[1]
}

func ParseEnviron(envs []string) map[string]string {
	envsMap := make(map[string]string)
	for _, env := range envs {
		key, value := ParseVar(env)
		envsMap[key] = value
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
