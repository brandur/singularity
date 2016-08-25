package singularity

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

const (
	// Release is the asset version of the site. Bump when any assets are
	// updated to blow away any browser caches.
	Release = "10"
)

const (
	// ArticlesDir is the location of site articles.
	ArticlesDir = "./articles/"

	// AssetsDir is the location of site static assets (e.g. CSS, images,
	// etc.).
	AssetsDir = "./assets/"

	// AssetsDir is the location of the site web fonts.
	FontsDir = "./fonts/"

	// LayoutsDir is the location of site layouts.
	LayoutsDir = "./layouts/"

	// MainLayout is the site's main layout.
	MainLayout = LayoutsDir + "/main"

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

// InitLog initializes logging for singularity programs.
func InitLog(verbose bool) {
	log.SetFormatter(&plainFormatter{})

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

// plainFormatter is a logrus formatter that displays text in a much more
// simple fashion that's more suitable as CLI output.
type plainFormatter struct {
}

// Format takes a logrus.Entry and returns bytes that are suitable for log
// output.
func (f *plainFormatter) Format(entry *log.Entry) ([]byte, error) {
	bytes := []byte(entry.Message + "\n")

	if entry.Level == log.DebugLevel {
		bytes = append([]byte("DEBUG: "), bytes...)
	}

	return bytes, nil
}
