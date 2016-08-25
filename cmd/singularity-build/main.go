package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/brandur/singularity"
	"github.com/brandur/singularity/templatehelpers"
	"github.com/brandur/sorg/markdown"
	"github.com/brandur/sorg/pool"
	"github.com/joeshaw/envdecode"
	"github.com/yosssi/ace"
)

// Conf contains configuration information for the command.
type Conf struct {
	// Concurrency is how main background Goroutines will be used to build all
	// site resources (e.g. articles, pages, etc.).
	Concurrency int `env:"CONCURRENCY,default=10"`

	// GoogleAnalyticsID is the account identifier for Google Analytics to use.
	GoogleAnalyticsID string `env:"GOOGLE_ANALYTICS_ID"`

	// LocalFonts starts using locally downloaded versions of Google Fonts.
	// This is not ideal for real deployment because you won't be able to
	// leverage Google's CDN and the caching that goes with it, and may not get
	// the font format for requesting browsers, but good for airplane rides
	// where you otherwise wouldn't have the fonts.
	LocalFonts bool `env:"LOCAL_FONTS,default=false"`

	// Verbose is whether the program will print debug output as it's running.
	Verbose bool `env:"VERBOSE,default=false"`
}

//
// Variables
//

// Left as a global for now for the sake of convenience, but it's not used in
// very many places and can probably be refactored as a local if desired.
var conf Conf

// pagesVars contains meta information for static pages that are part of the
// site. This mostly titles, but can also be body classes for custom styling.
//
// This isn't the best system, but was the cheapest way to accomplish what I
// needed for the time being. It could probably use an overhaul to something
// better at some point.
var pagesVars = map[string]map[string]interface{}{
	"about": {
		"Title": "About",
	},
	"index": {
		"Title": "The Self Hosting Singularity",
	},
}

//
// Main
//

func main() {
	start := time.Now()
	defer func() {
		log.Infof("Built site in %v.", time.Now().Sub(start))
	}()

	err := envdecode.Decode(&conf)
	if err != nil {
		log.Fatal(err)
	}

	singularity.InitLog(conf.Verbose)

	err = singularity.CreateOutputDirs(singularity.TargetDir)
	if err != nil {
		log.Fatal(err)
	}

	var tasks []*pool.Task

	tasks = append(tasks, pool.NewTask(func() error {
		return linkFontAssets()
	}))

	tasks = append(tasks, pool.NewTask(func() error {
		return linkImageAssets()
	}))

	articleTasks, err := tasksForArticles()
	if err != nil {
		log.Fatal(err)
	}
	tasks = append(tasks, articleTasks...)

	pageTasks, err := tasksForPages()
	if err != nil {
		log.Fatal(err)
	}
	tasks = append(tasks, pageTasks...)

	if !runTasks(tasks) {
		os.Exit(1)
	}
}

//
// Compilation functions
//
// These functions perform the heavy-lifting in compiling the site's resources.
// They are normally run concurrently.
//

func linkFontAssets() error {
	start := time.Now()
	defer func() {
		log.Debugf("Linked font assets in %v.", time.Now().Sub(start))
	}()

	source, err := filepath.Abs(singularity.FontsDir)
	if err != nil {
		return err
	}

	dest, err := filepath.Abs(singularity.TargetDir + "/assets/fonts/")
	if err != nil {
		return err
	}

	return ensureSymlink(source, dest)
}

func linkImageAssets() error {
	start := time.Now()
	defer func() {
		log.Debugf("Linked image assets in %v.", time.Now().Sub(start))
	}()

	err := os.RemoveAll(singularity.TargetDir + path.Clean(singularity.AssetsDir))
	if err != nil {
		return err
	}

	// we use absolute paths for source and destination because not doing so
	// can result in some weird symbolic link inception
	source, err := filepath.Abs(singularity.AssetsDir)
	if err != nil {
		return err
	}

	dest, err := filepath.Abs(singularity.TargetDir + singularity.AssetsDir)
	if err != nil {
		return err
	}

	err = ensureSymlink(source, dest)
	if err != nil {
		return err
	}

	return nil
}

func compileArticle(articleFile string) error {
	log.Debugf("Rendered article '%v'", articleFile)

	source, err := ioutil.ReadFile(singularity.ArticlesDir + articleFile)
	if err != nil {
		return err
	}
	rendered := markdown.Render(string(source))

	template, err := ace.Load(singularity.LayoutsDir+"main", singularity.LayoutsDir+"article", nil)
	if err != nil {
		return err
	}

	file, err := os.Create(singularity.TargetDir + trimExtension(articleFile))
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	err = template.Execute(writer, map[string]string{"Content": string(rendered)})
	if err != nil {
		return err
	}

	return nil
}

func compilePage(dir, name string) error {
	// Remove the "./pages" directory, but keep the rest of the path.
	//
	// Looks something like "about".
	pagePath := strings.TrimPrefix(dir, singularity.PagesDir) + name

	// Looks something like "./public/about".
	target := singularity.TargetDir + "/" + pagePath

	locals, ok := pagesVars[pagePath]
	if !ok {
		log.Errorf("No page meta information: %v", pagePath)
	}

	locals = getLocals("Page", locals)

	err := os.MkdirAll(singularity.TargetDir+"/"+dir, 0755)
	if err != nil {
		return err
	}

	err = renderView(singularity.MainLayout, dir+"/"+name,
		target, locals)
	if err != nil {
		return err
	}

	return nil
}

//
// Task generation functions
//
// These functions are the main entry points for compiling the site's
// resources.
//

func tasksForArticles() ([]*pool.Task, error) {
	files, err := ioutil.ReadDir(singularity.ArticlesDir)
	if err != nil {
		return nil, err
	}

	var tasks []*pool.Task
	for _, fileInfo := range files {
		// be careful with closures in loops
		name := fileInfo.Name()

		if isHidden(name) {
			continue
		}

		tasks = append(tasks, pool.NewTask(func() error {
			return compileArticle(name)
		}))
	}

	return tasks, nil
}

func tasksForPages() ([]*pool.Task, error) {
	return tasksForPagesDir(singularity.PagesDir)
}

func tasksForPagesDir(dir string) ([]*pool.Task, error) {
	log.Debugf("Descending into pages directory: %v", dir)

	files, err := ioutil.ReadDir(singularity.PagesDir)
	if err != nil {
		return nil, err
	}

	var tasks []*pool.Task
	for _, fileInfo := range files {
		// be careful with closures in loops
		name := fileInfo.Name()

		if fileInfo.IsDir() {
			subtasks, err := tasksForPagesDir(dir + name)
			if err != nil {
				return nil, err
			}
			tasks = append(tasks, subtasks...)
		} else {
			tasks = append(tasks, pool.NewTask(func() error {
				return compilePage(dir, trimExtension(name))
			}))
		}
	}

	return tasks, nil
}

//
// Other functions
//
// Any other functions. Try to keep them alphabetized.
//

func ensureSymlink(source, dest string) error {
	log.Debugf("Checking symbolic link (%v): %v -> %v",
		path.Base(source), source, dest)

	var actual string

	_, err := os.Stat(dest)

	// Note that if a symlink file does exist, but points to a non-existent
	// location, we still get an "does not exist" error back, so we fall down
	// to the general create path so that the symlink file can be removed.
	//
	// The call to RemoveAll does not affect the other path of the symlink file
	// not being present because it doesn't care whether or not the file it's
	// trying remove is actually there.
	if os.IsNotExist(err) {
		log.Debugf("Destination link does not exist. Creating.")
		goto create
	}
	if err != nil {
		return err
	}

	actual, err = os.Readlink(dest)
	if err != nil {
		return err
	}

	if actual == source {
		log.Debugf("Link exists.")
		return nil
	}

	log.Debugf("Destination links to wrong source. Creating.")

create:
	err = os.RemoveAll(dest)
	if err != nil {
		return err
	}

	return os.Symlink(source, dest)
}

// Gets a map of local values for use while rendering a template and includes
// a few "special" values that are globally relevant to all templates.
func getLocals(title string, locals map[string]interface{}) map[string]interface{} {
	defaults := map[string]interface{}{
		"GoogleAnalyticsID": conf.GoogleAnalyticsID,
		"LocalFonts":        conf.LocalFonts,
		"Release":           singularity.Release,
		"Title":             title,
		"ViewportWidth":     "device-width",
	}

	for k, v := range locals {
		defaults[k] = v
	}

	return defaults
}

func isHidden(file string) bool {
	return strings.HasPrefix(file, ".")
}

func renderView(layout, view, target string, locals map[string]interface{}) error {
	log.Debugf("Rendering: %v", target)

	template, err := ace.Load(layout, view, &ace.Options{FuncMap: templatehelpers.FuncMap})
	if err != nil {
		return err
	}

	file, err := os.Create(target)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	err = template.Execute(writer, locals)
	if err != nil {
		return err
	}

	return nil
}

// Runs the given tasks in a pool.
//
// After the run, if any errors occurred, it prints the first 10. Returns true
// if all tasks succeeded. If a false is returned, the caller should consider
// exiting with non-zero status.
func runTasks(tasks []*pool.Task) bool {
	p := pool.NewPool(tasks, conf.Concurrency)
	p.Run()

	var numErrors int
	for _, task := range p.Tasks {
		if task.Err != nil {
			log.Error(task.Err)
			numErrors++
		}
		if numErrors >= 10 {
			log.Error("Too many errors.")
			break
		}
	}

	return !p.HasErrors()
}

func trimExtension(file string) string {
	return strings.TrimSuffix(file, filepath.Ext(file))
}
