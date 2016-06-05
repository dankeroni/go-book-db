package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFind(t *testing.T) {
	findTests := []struct {
		id        string
		expectedC ClassifyBookResponse
		expectedR bool
		expectedE bool
	}{
		{
			id: "803736779",
			expectedC: ClassifyBookResponse{
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
			expectedC: ClassifyBookResponse{
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
		c, err := Find(test.id)

		if !test.expectedE {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		if test.expectedR {
			assert.Equal(t, test.expectedC, c)
		} else {
			assert.NotEqual(t, test.expectedC, c)
		}
	}
}
