package utils

import "net/url"

func ParseIconURL(basePath, iconURL string) string {
	if iconURL == "" {
		return iconURL
	}

	parsedIconURL, err := url.Parse(iconURL)
	if err != nil {
		return ""
	}
	if parsedIconURL.Host == "" {
		return basePath + "/dist/custom/" + iconURL
	}
	return iconURL
}
