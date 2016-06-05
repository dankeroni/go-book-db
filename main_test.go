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
		o, err := Find(test.id)

		if !test.expectedE {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		if test.expectedR {
			assert.Equal(t, test.expectedO, o)
		} else {
			assert.NotEqual(t, test.expectedO, o)
		}
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
		o, err := Search(test.query)

		if !test.expectedE {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		if test.expectedR {
			assert.Equal(t, test.expectedO, o)
		} else {
			assert.NotEqual(t, test.expectedO, o)
		}
	}
}
