///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package utils

import (
	"net/url"
	"strings"
)

func AbsoluteURL(parent, subUrl string) (string, error) {
	if strings.HasPrefix(subUrl, "#") {
		return "", nil
	}
	base, err := url.Parse(parent)
	if err != nil {
		return "", err
	}

	absURL, err := base.Parse(subUrl)
	if err != nil {
		return "", err
	}
	absURL.Fragment = ""
	if absURL.Scheme == "//" {
		absURL.Scheme = base.Scheme
	}
	return absURL.String(), nil
}

func Parse(u string) (*url.URL, error) {
	if !strings.HasPrefix(u, "https://") && !strings.HasPrefix(u, "http://") {
		u = "http://" + u
	}
	return url.Parse(u)
}
