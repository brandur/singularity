package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/brandur/singularity"
	"github.com/russross/blackfriday"
	"github.com/yosssi/ace"
)

var (
	concurrency = 10
	errors      = make(chan error)
)

func main() {
	// We should probably have a more complete approach to error handling here,
	// but for now just error on the first problem.
	go func() {
		for err := range errors {
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	err := singularity.CreateTargetDir()
	errors <- err

	build()
}

func build() {
	var verbose bool
	if os.Getenv("VERBOSE") == "true" {
		verbose = true
	}

	singularity.InitLog(verbose)

	if os.Getenv("CONCURRENCY") != "" {
		c, err := strconv.Atoi(os.Getenv("CONCURRENCY"))
		errors <- err
		if c < 1 {
			errors <- fmt.Errorf("CONCURRENCY must be >= 1")
		}
		concurrency = c
	}

	start := time.Now()
	defer func() {
		log.Infof("Site built in %v", time.Now().Sub(start))
	}()

	log.Debugf("Starting build with concurrency %v", concurrency)

	var wg sync.WaitGroup

	// note that if this buffered channel fills, the producers might block, but
	// that's not a big deal
	jobs := make(chan func() error, 1000)

	for i := 0; i < concurrency; i++ {
		go func() {
			for job := range jobs {
				errors <- job()
				wg.Done()
			}
		}()
	}

	// we build jobs for everything and then just run it all in parallel
	wg.Add(1)
	jobs <- func() error {
		return linkAssets()
	}

	articleJobs, err := generateArticleJobs()
	errors <- err

	pageJobs, err := generatePageJobs()
	errors <- err

	wg.Add(len(articleJobs) + len(pageJobs))
	for _, job := range append(articleJobs, pageJobs...) {
		jobs <- job
	}

	wg.Wait()
}

func generateArticleJobs() ([]func() error, error) {
	files, err := ioutil.ReadDir(singularity.ArticlesDir)
	if err != nil {
		return nil, err
	}

	jobs := make([]func() error, len(files))

	for i, fileInfo := range files {
		// be careful with closures in loops
		localFileInfo := fileInfo

		jobs[i] = func() error {
			return renderArticle(localFileInfo.Name())
		}
	}

	return jobs, nil
}

func generatePageJobs() ([]func() error, error) {
	files, err := ioutil.ReadDir(singularity.PagesDir)
	if err != nil {
		return nil, err
	}

	jobs := make([]func() error, len(files))

	for i, fileInfo := range files {
		// be careful with closures in loops
		localFileInfo := fileInfo

		jobs[i] = func() error {
			// ace.Load doesn't like to take .ace extensions anyway
			return renderPage(trimExtension(localFileInfo.Name()))
		}
	}

	return jobs, nil
}

func linkAssets() error {
	log.Debugf("Linking assets directory")

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

	err = os.Symlink(source, dest)
	if err != nil {
		return err
	}

	return nil
}

func renderArticle(articleFile string) error {
	log.Debugf("Rendered article '%v'", articleFile)

	source, err := ioutil.ReadFile(singularity.ArticlesDir + articleFile)
	if err != nil {
		return err
	}
	rendered := renderMarkdown(source)

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

func renderMarkdown(source []byte) []byte {
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_DASHES
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_USE_XHTML

	extensions := 0
	extensions |= blackfriday.EXTENSION_AUTO_HEADER_IDS
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_HEADER_IDS
	extensions |= blackfriday.EXTENSION_LAX_HTML_BLOCKS
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH

	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")
	return blackfriday.Markdown(source, renderer, extensions)
}

func renderPage(pageFile string) error {
	log.Debugf("Rendered page '%v'", pageFile)

	template, err := ace.Load(singularity.LayoutsDir+"main", singularity.PagesDir+pageFile, nil)
	if err != nil {
		return err
	}

	file, err := os.Create(singularity.TargetDir + pageFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	err = template.Execute(writer, map[string]string{})
	if err != nil {
		return err
	}

	return nil
}

func trimExtension(file string) string {
	return strings.TrimSuffix(file, filepath.Ext(file))
}
