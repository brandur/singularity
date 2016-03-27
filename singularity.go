package singularity

import (
	"os"
)

const (
	// ArticlesDir is the location of site articles.
	ArticlesDir = "./articles/"

	// AssetsDir is the location of site static assets (e.g. CSS, images,
	// etc.).
	AssetsDir = "./assets/"

	// LayoutsDir is the location of site layouts.
	LayoutsDir = "./layouts/"

	// PagesDir is the location of site static pages (i.e. mostly any page
	// that isn't an article).
	PagesDir = "./pages/"

	// TargetDir is the location where the site will be built to.
	TargetDir = "./public/"
)

// CreateTargetDir creates TargetDir if it doesn't already exist.
func CreateTargetDir() error {
	return os.MkdirAll(TargetDir, 0755)
}
