package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFind(t *testing.T) {
	findTests := []struct {
		id        string
		expectedO ClassifyBookResponse
		expectedR bool
		expectedE bool
	}{
		{
			id: "803736779",
			expectedO: ClassifyBookResponse{
				BookData: BookData{
					Title:  "Cloning internet applications with Ruby : make your own TinyURL, Twitter, Flickr, or Facebook using Ruby",
					Author: "Chang, Sau Sheong",
					ID:     "803736779",
				},
				Classification: Classification{},
			},
			expectedR: true,
			expectedE: false,
		},
		{
			id: "qwegbvuhe",
			expectedO: ClassifyBookResponse{
				BookData: BookData{
					Title:  "Cloning internet applications with Ruby : make your own TinyURL, Twitter, Flickr, or Facebook using Ruby",
					Author: "Chang, Sau Sheong",
					ID:     "803736779",
				},
				Classification: Classification{},
			},
			expectedR: false,
			expectedE: false,
		},
	}

	for _, test := range findTests {
		o, e := Find(test.id)
		RunSimpleTest(t, test.expectedO, o, test.expectedE, e, test.expectedR)
	}
}

func TestSearch(t *testing.T) {
	searchTests := []struct {
		query     string
		expectedO []SearchResult
		expectedR bool
		expectedE bool
	}{
		{
			query: "Cloning internet",
			expectedO: []SearchResult{
				{
					Title:  "Cloning internet applications with Ruby : make your own TinyURL, Twitter, Flickr, or Facebook using Ruby",
					Author: "Chang, Sau Sheong",
					Year:   "2010",
					ID:     "803736779",
				},
				{
					Title:  "Cloning Internet Applications with Ruby",
					Author: "",
					Year:   "2010",
					ID:     "918749340",
				},
			},
			expectedR: true,
			expectedE: false,
		},
		{
			query: "Cloning internet",
			expectedO: []SearchResult{
				{
					Title:  "Cloning internet applications with Ruby : make your own TinyURL, Twitter, Flickr, or Facebook using Ruby",
					Author: "Chang, Sau Sheong",
					Year:   "2010",
					ID:     "803736779",
				},
			},
			expectedR: false,
			expectedE: false,
		},
	}

	for _, test := range searchTests {
		o, e := Search(test.query)
		RunSimpleTest(t, test.expectedO, o, test.expectedE, e, test.expectedR)
	}
}

func RunSimpleTest(t *testing.T, expectedO interface{}, o interface{}, expectedE bool, e error, expectedR bool) {
	if !expectedE {
		assert.Nil(t, e)
	} else {
		assert.NotNil(t, e)
	}

	if expectedR {
		assert.Equal(t, expectedO, o)
	} else {
		assert.NotEqual(t, expectedO, o)
	}
}
