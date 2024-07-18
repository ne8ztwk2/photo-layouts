package photoLayouts_test

import (
	"testing"

	photoLayouts "github.com/leptobo/photo-layouts"
)

func TestLayout(t *testing.T) {
	photos := []photoLayouts.Photo{
		{File: "example/James.jpg", Dpi: 300, PhotoWidth: 25, PhotoHeight: 35, ContainerWidth: 127, ContainerHeight: 89, Color: "#000000"},    // 1 5
		{File: "example/Michael.jpg", Dpi: 300, PhotoWidth: 38, PhotoHeight: 51, ContainerWidth: 127, ContainerHeight: 89, Color: "#FFFFFF"},  // 2 5
		{File: "example/William.jpg", Dpi: 300, PhotoWidth: 25, PhotoHeight: 35, ContainerWidth: 152, ContainerHeight: 102, Color: "#808080"}, // 1 6
		{File: "example/John.jpg", Dpi: 300, PhotoWidth: 38, PhotoHeight: 51, ContainerWidth: 152, ContainerHeight: 102, Color: "#808080"},    // 2 6
		// add more...
	}
	for _, photo := range photos {
		if _, err := photoLayouts.Layout(&photo); err != nil {
			t.Error(err)
		}
	}
}
