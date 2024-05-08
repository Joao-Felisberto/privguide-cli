// Miscelaneous utilities used in the project
package util

import (
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var Pipeline = false

func SetupLogger() {
	if !Pipeline {
		slog.SetDefault(slog.New(NewHumanFriendlyHandler(nil)))
	}
}

// The generic iterator map function
//
// `arr`: the original array
//
// `mapper`: the function used to map each element
//
// returns: the array containing the result of applying the `mapper` function to each element of `arr`
func Map[T1 any, T2 any](arr []T1, mapper func(T1) T2) []T2 {
	new := []T2{}

	for _, e := range arr {
		new = append(new, mapper(e))
	}

	return new
}

// The generic iterator filter function
//
// `arr`: the original array from which to filter
//
// `filter`: The function that decides whether an element should be in the final array
//
// returns: The elements of `arr` for which `filter` returned `true`
func Filter[T any](arr []T, filter func(T) bool) []T {
	new := []T{}

	for _, e := range arr {
		if filter(e) {
			new = append(new, e)
		}
	}

	return new
}

// Cast a generic map into a map of a specific key and value types. Errors will make the function panic
//
// `m`: The map to convert
//
// returns: The converted map
func MapCast[K comparable, V any](m map[interface{}]interface{}) map[K]V {
	newMap := map[K]V{}

	for k, v := range m {
		newMap[k.(K)] = v.(V)
	}

	return newMap
}

// Converts an array into a map through a mapping function that returns a key-value pair
//
// `arr`: The starting array
//
// `mapper`: The function to map an element of the array into a key-value pair
//
// returns: The map of keys to values
func ArrayToMap[T any, K comparable, V any](arr []T, mapper func(T) (K, V)) map[K]V {
	res := map[K]V{}

	for _, e := range arr {
		k, v := mapper(e)
		res[k] = v
	}

	return res
}

// Finds out whether at least one element of an array
//
// `arr`: the array
//
// `condition`: the function at least one element of the array should obbey
//
// returns: Whether at least one element of the array obeys the condition
func Any[T any](arr []T, condition func(T) bool) bool {
	for _, e := range arr {
		if condition(e) {
			return true
		}
	}

	return false
}

// Compares two arrays disregarding order, as if they were sets
//
// `set1`: the first array
//
// `set2`: the second array
//
// returns: whether the sets are equal
func CompareSets[T comparable](set1 []T, set2 []T) bool {
	if len(set1) != len(set2) {
		return false
	}

	for _, e := range set1 {
		if !slices.Contains(set2, e) {
			return false
		}
	}

	return true
}

// Creates a file with the given data as string.
// Also creates parent directories as needed.
// Directories have permissions 0766 and files 0666.
//
// `filePath`: where to store the file
//
// `data`: the string to write to the file
//
// returns: an error if the file could not be written to or the directories could not be created
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

// Deletes the file and each of the parent directories.
//
// `filePath`: The path to the file to delete
func DeleteFileAndParentPath(filePath string) {
	path := strings.Split(filePath, "/")
	for i := len(path); i >= 0; i-- {
		path := filepath.Join(path[:i]...)
		slog.Info("deleting", "full", filePath, "to delete", path)
		os.Remove(path)
	}
}
