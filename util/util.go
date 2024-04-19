package util

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func Map[T1 any, T2 any](arr []T1, mapper func(T1) T2) []T2 {
	new := []T2{}

	for _, e := range arr {
		new = append(new, mapper(e))
	}

	return new
}

func Filter[T any](arr []T, filterFn func(T) bool) []T {
	new := []T{}

	for _, e := range arr {
		if filterFn(e) {
			new = append(new, e)
		}
	}

	return new
}

func MapCast[K comparable, V any](m map[interface{}]interface{}) map[K]V {
	newMap := map[K]V{}

	for k, v := range m {
		newMap[k.(K)] = v.(V)
	}

	return newMap
}

func MapToMap[T any, K comparable, V any](list []T, mapper func(T) (K, V)) map[K]V {
	res := map[K]V{}

	for _, e := range list {
		k, v := mapper(e)
		res[k] = v
	}

	return res
}

func Any[T any](list []T, condition func(T) bool) bool {
	for _, e := range list {
		if condition(e) {
			return true
		}
	}

	return false
}

func CreateFileWithData(filePath string, data string) error {
	path := strings.Split(filePath, "/")
	if len(path) > 1 {
		dirs := path[:len(path)-1]

		err := os.MkdirAll(filepath.Join(dirs...), os.ModePerm)
		if err != nil {
			return err
		}
	}

	err := os.WriteFile(filePath, []byte(data), 0666)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFileAndParentPath(filePath string) {
	path := strings.Split(filePath, "/")
	for i := len(path); i >= 0; i-- {
		path := filepath.Join(path[:i]...)
		slog.Info("deleting", "full", filePath, "to delete", path)
		os.Remove(path)
	}
}
