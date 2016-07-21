/*
An implementation of the XDG Base Directory Specification.

For more information, see: http://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html
*/
package xdg

import (
	"os"
	"path/filepath"
)

func IsValid(path string) bool {
	return path != "" && filepath.IsAbs(string(path))
}

type Path string

func PathWithDefault(path string, defaultPath Path) Path {
	if IsValid(path) {
		return Path(path)
	}

	return defaultPath
}

func (p Path) IsValid() bool {
	return IsValid(string(p))
}

func (p Path) Find(elem ...string) string {
	file := p.Join(elem...)
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return ""
	}

	return file
}

func (p Path) Glob(pattern string) ([]string, error) {
	return filepath.Glob(filepath.Join(string(p), pattern))
}

func (p Path) Join(elem ...string) string {
	return filepath.Join(append([]string{string(p)}, elem...)...)
}

type Paths []Path

func PathsWithDefault(paths []string, defaultPaths Paths) Paths {
	var p Paths
	for _, path := range paths {
		if IsValid(path) {
			p = append(p, Path(path))
		}
	}

	if len(p) == 0 {
		for _, path := range defaultPaths {
			p = append(p, path)
		}
	}

	return p
}

func (p Paths) Find(elem ...string) []string {
	var found []string
	for _, path := range p {
		file := path.Find(elem...)
		if file != "" {
			found = append(found, file)
		}
	}
	return found
}

func (p Paths) Glob(pattern string) ([]string, error) {
	var matches []string
	for _, path := range p {
		pathMatches, err := path.Glob(pattern)
		if err != nil {
			return nil, err
		}

		matches = append(matches, pathMatches...)
	}

	return matches, nil
}

func (p Paths) Join(elem ...string) []string {
	pathElem := append([]string{""}, elem...)

	var joined []string
	for _, path := range p {
		pathElem[0] = string(path)
		joined = append(joined, filepath.Join(pathElem...))
	}

	return joined
}
