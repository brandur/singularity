package singularity

import (
	"os"
	"path"
	"time"

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

	// FontsDir is the location of the site web fonts.
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

// A list of all directories that are in the built static site.
var outputDirs = []string{
	".",
	"assets",
	"assets/" + Release,
	"fonts",
}

// CreateOutputDirs creates a target directory for the static site and all
// other necessary directories for the build if they don't already exist.
func CreateOutputDirs(targetDir string) error {
	start := time.Now()
	defer func() {
		log.Debugf("Created target directories in %v.", time.Now().Sub(start))
	}()

	for _, dir := range outputDirs {
		dir = path.Join(targetDir, dir)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
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
