package discordkit

import (
	"fmt"
	"strings"
)

type componentRoute struct {
	pattern     string
	segments    []string
	staticCount int
	handler     ComponentHandler
}

func newComponentRoute(
	pattern string,
	handler ComponentHandler,
) (componentRoute, error) {
	pattern = normalizeComponentPath(pattern)

	if pattern == "/" {
		return componentRoute{}, fmt.Errorf("component route cannot be root")
	}

	segments := splitComponentPath(pattern)

	seenParams := make(map[string]struct{})

	staticCount := 0

	for _, segment := range segments {
		if !strings.HasPrefix(segment, ":") {
			staticCount++
			continue
		}

		name := strings.TrimPrefix(segment, ":")

		if name == "" {
			return componentRoute{}, fmt.Errorf("route %q contains an unnamed parameter", pattern)
		}

		if _, exists := seenParams[name]; exists {
			return componentRoute{}, fmt.Errorf(
				"route %q contains duplicate parameter %q",
				pattern,
				name,
			)
		}

		seenParams[name] = struct{}{}
	}

	return componentRoute{
		pattern:     pattern,
		segments:    segments,
		staticCount: staticCount,
		handler:     handler,
	}, nil
}

func (r componentRoute) match(path string) (map[string]string, bool) {
	path = normalizeComponentPath(path)

	segments := splitComponentPath(path)

	if len(segments) != len(r.segments) {
		return nil, false
	}

	params := make(map[string]string)

	for index, expected := range r.segments {
		actual := segments[index]

		if strings.HasPrefix(expected, ":") {
			name := strings.TrimPrefix(expected, ":")

			params[name] = actual

			continue
		}

		if expected != actual {
			return nil, false
		}
	}

	return params, true
}

func (r componentRoute) signature() string {
	parts := make([]string, len(r.segments))

	for index, segment := range r.segments {
		if strings.HasPrefix(segment, ":") {
			parts[index] = ":"
			continue
		}

		parts[index] = segment
	}

	return strings.Join(parts, "/")
}

func normalizeComponentPath(path string) string {
	path = strings.TrimSpace(path)

	if path == "" {
		return "/"
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if len(path) > 1 {
		path = strings.TrimSuffix(path, "/")
	}

	return path
}

func splitComponentPath(path string) []string {
	path = strings.Trim(path, "/")

	if path == "" {
		return nil
	}

	return strings.Split(path, "/")
}
