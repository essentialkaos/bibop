package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"strings"

	"pkg.re/essentialkaos/ek.v10/req"
	"pkg.re/essentialkaos/ek.v10/strutil"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// HTTPStatus is action processor for "http-status"
func HTTPStatus(action *recipe.Action) error {
	method, err := action.GetS(0)

	if err != nil {
		return err
	}

	url, err := action.GetS(1)

	if err != nil {
		return err
	}

	code, err := action.GetI(2)

	if err != nil {
		return err
	}

	if !isHTTPMethodSupported(method) {
		return fmt.Errorf("Method %s is not supported", method)
	}

	resp, err := makeHTTPRequest(method, url).Do()

	if err != nil {
		return fmt.Errorf("Can't send HTTP request %s %s", method, url)
	}

	switch {
	case !action.Negative && resp.StatusCode != code:
		return fmt.Errorf("HTTP request returns different status code (%d ≠ %d)", resp.StatusCode, code)
	case action.Negative && resp.StatusCode == code:
		return fmt.Errorf("HTTP request return invalid status code (%d)", code)
	}

	return nil
}

// HTTPHeader is action processor for "http-header"
func HTTPHeader(action *recipe.Action) error {
	method, err := action.GetS(0)

	if err != nil {
		return err
	}

	url, err := action.GetS(1)

	if err != nil {
		return err
	}

	headerName, err := action.GetS(2)

	if err != nil {
		return err
	}

	headerValue, err := action.GetS(3)

	if err != nil {
		return err
	}

	if !isHTTPMethodSupported(method) {
		return fmt.Errorf("Method %s is not supported", method)
	}

	resp, err := makeHTTPRequest(method, url).Do()

	if err != nil {
		return fmt.Errorf("Can't send HTTP request %s %s", method, url)
	}

	isHeaderPresent := resp.Header.Get(headerName) == headerValue

	switch {
	case !action.Negative && !isHeaderPresent:
		return fmt.Errorf(
			"HTTP request returns different header (%d ≠ %d)",
			resp.Header.Get(headerName), headerValue,
		)
	case action.Negative && isHeaderPresent:
		return fmt.Errorf("HTTP request return invalid header (%d)", headerValue)
	}

	return nil
}

// HTTPContains is action processor for "http-contains"
func HTTPContains(action *recipe.Action) error {
	method, err := action.GetS(0)

	if err != nil {
		return err
	}

	url, err := action.GetS(1)

	if err != nil {
		return err
	}

	substr, err := action.GetS(2)

	if err != nil {
		return err
	}

	if !isHTTPMethodSupported(method) {
		return fmt.Errorf("Method %s is not supported", method)
	}

	resp, err := makeHTTPRequest(method, url).Do()

	if err != nil {
		return fmt.Errorf("Can't send HTTP request %s %s", method, url)
	}

	containsSubstr := strings.Contains(resp.String(), substr)

	switch {
	case !action.Negative && !containsSubstr:
		return fmt.Errorf("HTTP request response doesn't contain given substring")
	case action.Negative && containsSubstr:
		return fmt.Errorf("HTTP request response contains given substring")
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// isHTTPMethodSupported returns true if HTTP method is supported
func isHTTPMethodSupported(method string) bool {
	switch method {
	case req.GET, req.POST, req.DELETE,
		req.PUT, req.PATCH, req.HEAD:
		return true
	}

	return false
}

// makeHTTPRequest creates request struct
func makeHTTPRequest(method, url string) *req.Request {
	if !strings.Contains(url, "@") {
		return &req.Request{Method: method, URL: url, AutoDiscard: true, FollowRedirect: true}
	}

	auth := strutil.ReadField(url, 0, false, "@")
	auth = strings.Replace(auth, "http://", "", -1)
	auth = strings.Replace(auth, "https://", "", -1)

	url = strings.Replace(url, auth+"@", "", -1)

	login := strutil.ReadField(auth, 0, false, ":")
	pass := strutil.ReadField(auth, 1, false, ":")

	return &req.Request{
		Method:            method,
		URL:               url,
		BasicAuthUsername: login,
		BasicAuthPassword: pass,
		AutoDiscard:       true,
		FollowRedirect:    true,
	}
}
