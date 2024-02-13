package action

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"strings"

	"github.com/essentialkaos/ek/v12/req"

	"github.com/buger/jsonparser"

	"github.com/essentialkaos/bibop/recipe"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	PROP_HTTP_REQUEST_HEADERS = "HTTP_REQUEST_HEADERS"
	PROP_HTTP_AUTH_USERNAME   = "HTTP_AUTH_USERNAME"
	PROP_HTTP_AUTH_PASSWORD   = "HTTP_AUTH_PASSWORD"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// HTTPStatus is action processor for "http-status"
func HTTPStatus(action *recipe.Action) error {
	var payload string

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

	if action.Has(3) {
		payload, _ = action.GetS(3)
	}

	err = checkRequestData(method, payload)

	if err != nil {
		return err
	}

	resp, err := makeHTTPRequest(action, method, url, payload).Do()

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
	var payload string

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

	if action.Has(4) {
		payload, _ = action.GetS(4)
	}

	err = checkRequestData(method, payload)

	if err != nil {
		return err
	}

	resp, err := makeHTTPRequest(action, method, url, payload).Do()

	if err != nil {
		return fmt.Errorf("Can't send HTTP request %s %s", method, url)
	}

	isHeaderPresent := resp.Header.Get(headerName) == headerValue

	switch {
	case !action.Negative && !isHeaderPresent:
		return fmt.Errorf(
			"HTTP request returns different header (%s ≠ %s)",
			fmtValue(resp.Header.Get(headerName)), headerValue,
		)
	case action.Negative && isHeaderPresent:
		return fmt.Errorf("HTTP request return invalid header (%s)", headerValue)
	}

	return nil
}

// HTTPContains is action processor for "http-contains"
func HTTPContains(action *recipe.Action) error {
	var payload string

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

	if action.Has(3) {
		payload, _ = action.GetS(3)
	}

	err = checkRequestData(method, payload)

	if err != nil {
		return err
	}

	resp, err := makeHTTPRequest(action, method, url, payload).Do()

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

// HTTPJSON is action processor for "http-json"
func HTTPJSON(action *recipe.Action) error {
	method, err := action.GetS(0)

	if err != nil {
		return err
	}

	url, err := action.GetS(1)

	if err != nil {
		return err
	}

	query, err := action.GetS(2)

	if err != nil {
		return err
	}

	value, err := action.GetS(3)

	if err != nil {
		return err
	}

	err = checkRequestData(method, "")

	if err != nil {
		return err
	}

	resp, err := makeHTTPRequest(action, method, url, "").Do()

	if err != nil {
		return fmt.Errorf("Can't send HTTP request %s %s", method, url)
	}

	querySlice := parseJSONQuery(query)
	jsonValue, _, _, err := jsonparser.Get(resp.Bytes(), querySlice...)

	if err != nil {
		return fmt.Errorf("Can't get JSON data: %v", err)
	}

	containsValue := string(jsonValue) == value

	switch {
	case !action.Negative && !containsValue:
		return fmt.Errorf("JSON response doesn't contain given value")
	case action.Negative && containsValue:
		return fmt.Errorf("JSON response contains given value")
	}

	return nil
}

// HTTPSetAuth is action processor for "http-set-auth"
func HTTPSetAuth(action *recipe.Action) error {
	command := action.Command

	username, err := action.GetS(0)

	if err != nil {
		return err
	}

	password, err := action.GetS(1)

	if err != nil {
		return err
	}

	command.Data.Set(PROP_HTTP_AUTH_USERNAME, username)
	command.Data.Set(PROP_HTTP_AUTH_PASSWORD, password)

	return nil
}

// HTTPSetHeader is action processor for "http-set-header"
func HTTPSetHeader(action *recipe.Action) error {
	command := action.Command

	headerName, err := action.GetS(0)

	if err != nil {
		return err
	}

	headerValue, err := action.GetS(1)

	if err != nil {
		return err
	}

	var headers req.Headers

	if !command.Data.Has(PROP_HTTP_REQUEST_HEADERS) {
		headers = req.Headers{}
	} else {
		headers = command.Data.Get(PROP_HTTP_REQUEST_HEADERS).(req.Headers)
	}

	headers[headerName] = headerValue

	command.Data.Set(PROP_HTTP_REQUEST_HEADERS, headers)

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// checkRequestData checks request data
func checkRequestData(method, payload string) error {
	switch method {
	case req.GET, req.POST, req.DELETE, req.PUT, req.PATCH, req.HEAD:
		// NOOP
	default:
		return fmt.Errorf("Method %s is not supported", method)
	}

	switch method {
	case req.GET, req.DELETE, req.HEAD:
		if payload != "" {
			return fmt.Errorf("Method %s does not support payload", method)
		}
	}

	return nil
}

// makeHTTPRequest creates request struct
func makeHTTPRequest(action *recipe.Action, method, url, payload string) *req.Request {
	command := action.Command
	request := &req.Request{
		Method:         method,
		URL:            url,
		AutoDiscard:    true,
		FollowRedirect: true,
	}

	if payload != "" {
		request.Body = payload
	}

	if command.Data.Has(PROP_HTTP_AUTH_USERNAME) && command.Data.Has(PROP_HTTP_AUTH_PASSWORD) {
		request.BasicAuthUsername = command.Data.Get(PROP_HTTP_AUTH_USERNAME).(string)
		request.BasicAuthPassword = command.Data.Get(PROP_HTTP_AUTH_PASSWORD).(string)
	}

	if command.Data.Has(PROP_HTTP_REQUEST_HEADERS) {
		request.Headers = command.Data.Get(PROP_HTTP_REQUEST_HEADERS).(req.Headers)
	}

	return request
}

// parseJSONQuery converts json query to slice
func parseJSONQuery(q string) []string {
	q = strings.ReplaceAll(q, "[", ".[")
	return strings.Split(q, ".")
}
