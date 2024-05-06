// TODO: Fix documentation
// Copyright 2024-present Olivere. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v8/uritemplates"
	"net/http"
	"net/url"
	"strings"
)

// SearchTemplateService returns the list of indices plus some additional
// information about them.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/current/search-template-api.html
// for details.
type SearchTemplateService struct {
	client *Client

	pretty     *bool    // pretty format the returned JSON response
	human      *bool    // return human readable values for statistics
	errorTrace *bool    // include the stack trace of returned errors
	filterPath []string // list of filters used to reduce the response
	headers    http.Header

	index  string
	id     string
	params map[string]any
}

// NewCatIndicesService creates a new CatIndicesService.
func NewSearchTemplateService(client *Client) *SearchTemplateService {
	return &SearchTemplateService{
		client: client,
	}
}

// Pretty tells Elasticsearch whether to return a formatted JSON response.
func (s *SearchTemplateService) Pretty(pretty bool) *SearchTemplateService {
	s.pretty = &pretty
	return s
}

// Human specifies whether human readable values should be returned in
// the JSON response, e.g. "7.5mb".
func (s *SearchTemplateService) Human(human bool) *SearchTemplateService {
	s.human = &human
	return s
}

// ErrorTrace specifies whether to include the stack trace of returned errors.
func (s *SearchTemplateService) ErrorTrace(errorTrace bool) *SearchTemplateService {
	s.errorTrace = &errorTrace
	return s
}

// FilterPath specifies a list of filters used to reduce the response.
func (s *SearchTemplateService) FilterPath(filterPath ...string) *SearchTemplateService {
	s.filterPath = filterPath
	return s
}

// Header adds a header to the request.
func (s *SearchTemplateService) Header(name string, value string) *SearchTemplateService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Headers specifies the headers of the request.
func (s *SearchTemplateService) Headers(headers http.Header) *SearchTemplateService {
	s.headers = headers
	return s
}

// Index is the name of the index to list (by default all indices are returned).
func (s *SearchTemplateService) Index(index string) *SearchTemplateService {
	s.index = index
	return s
}

// Index is the name of the index to list (by default all indices are returned).
func (s *SearchTemplateService) Id(id string) *SearchTemplateService {
	s.id = id
	return s
}

// Index is the name of the index to list (by default all indices are returned).
func (s *SearchTemplateService) Params(params map[string]any) *SearchTemplateService {
	s.params = params
	return s
}

// buildURL builds the URL for the operation.
func (s *SearchTemplateService) buildURL() (string, url.Values, error) {
	// Build URL
	var (
		path string
		err  error
	)

	if s.index != "" {
		path, err = uritemplates.Expand("{index}/_search/template", map[string]string{
			"index": s.index,
		})
	} else {
		path = "_search/template"
	}
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{
		"format": []string{"json"}, // always returns as JSON
	}
	if v := s.pretty; v != nil {
		params.Set("pretty", fmt.Sprint(*v))
	}
	if v := s.human; v != nil {
		params.Set("human", fmt.Sprint(*v))
	}
	if v := s.errorTrace; v != nil {
		params.Set("error_trace", fmt.Sprint(*v))
	}
	if len(s.filterPath) > 0 {
		params.Set("filter_path", strings.Join(s.filterPath, ","))
	}
	if s.id != "" {
		params.Set("id", s.id)
	}
	if s.params != nil {
		serializedParams, err := json.Marshal(s.params)
		if err != nil {
			return "", url.Values{}, err
		}

		params.Set("params", string(serializedParams))
	}

	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *SearchTemplateService) Validate() error {
	return nil
}

// Do executes the search and returns a SearchResult.
func (s *SearchTemplateService) Do(ctx context.Context) (*SearchResult, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	// Perform request
	body, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformRequest(ctx, PerformRequestOptions{
		Method:  "POST",
		Path:    path,
		Params:  params,
		Body:    body,
		Headers: s.headers,
		// TODO: Return to this once you figure out the rest
		//MaxResponseSize: s.maxResponseSize,
	})
	if err != nil {
		return nil, err
	}

	// Return search results
	ret := new(SearchResult)
	if err := s.client.decoder.Decode(res.Body, ret); err != nil {
		ret.Header = res.Header
		return nil, err
	}
	ret.Header = res.Header
	return ret, nil
}
