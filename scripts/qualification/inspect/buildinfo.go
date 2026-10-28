package main

import "debug/buildinfo"

func readInfo(file string) (map[string]string, error) {
	info, err := buildinfo.ReadFile(file)
	if err != nil {
		return nil, err
	}
	result := map[string]string{"go_version": info.GoVersion}
	for _, setting := range info.Settings {
		result[setting.Key] = setting.Value
	}
	return result, nil
}
