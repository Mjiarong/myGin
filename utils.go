package myGin

import "path"

func assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	return finalPath
}

type H map[string]interface{}
