package values

import (
	"fmt"
	"net/url"
	"strings"
)

type Handle struct {
	value string
}

func TryNewHandleFromUrl(urlString string) (*Handle, error) {
	parsedUrl, err := url.ParseRequestURI(urlString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL %s: %w", urlString, err)
	}

	if parsedUrl.Host != "www.youtube.com" {
		return nil, fmt.Errorf("%s is not Youtube URL", urlString)
	}

	pathElements := strings.Split(parsedUrl.Path, "/")
	handleValue := pathElements[1]

	return TryNewHandle(handleValue)
}

func TryNewHandle(value string) (*Handle, error) {
	if !strings.HasPrefix(value, "@") {
		return nil, fmt.Errorf("handle %s does not start with @", value)
	}

	trimmedValue := strings.TrimSpace(value)
	if len(trimmedValue) < 2 {
		return nil, fmt.Errorf("handle %s is too short", trimmedValue)
	}

	return &Handle{
		value: trimmedValue,
	}, nil
}

func (handle *Handle) String() string {
	return handle.value
}
