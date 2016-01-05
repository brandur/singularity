package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/russross/blackfriday"
	"github.com/yosssi/ace"
)

const (
	ArticlesDir = "./articles/"
	AssetsDir   = "./assets/"
	LayoutsDir  = "./layouts/"
	PagesDir    = "./pages/"
	TargetDir   = "./public/"
)

var (
	verbose = false
)

func main() {
	start := time.Now()
	defer func() {
		fmt.Printf("Site built in %v\n", time.Now().Sub(start))
	}()

	if os.Getenv("VERBOSE") == "true" {
		verbose = true
	}

	// create an output directory
	err := os.MkdirAll(TargetDir, 0755)
	if err != nil {
		panic(err)
	}

	articleJobs, err := generateArticleJobs()
	if err != nil {
		panic(err)
	}

	pageJobs, err := generatePageJobs()
	if err != nil {
		panic(err)
	}

	// we build jobs for everything and then just run it all in parallel
	var jobs []func() error
	jobs = append(jobs, func() error {
		return linkAssets()
	})
	jobs = append(jobs, articleJobs...)
	jobs = append(jobs, pageJobs...)

	errors := make([]error, len(jobs))

	var wg sync.WaitGroup
	for i, job := range jobs {
		wg.Add(1)

		// be careful with closures in loops
		localJob := job

		go func() {
			defer wg.Done()

			// avoiding an append() keeps this safe between goroutines
			errors[i] = localJob()
		}()
	}
	wg.Wait()

	// should probably have a more complete approach to error handling here
	for _, err := range errors {
		if err != nil {
			panic(err)
		}
	}
}

func generateArticleJobs() ([]func() error, error) {
	files, err := ioutil.ReadDir(ArticlesDir)
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
	files, err := ioutil.ReadDir(PagesDir)
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
	if verbose {
		fmt.Printf("Linking assets directory\n")
	}

	err := os.RemoveAll(TargetDir + AssetsDir)
	if err != nil {
		return err
	}

	// we use absolute paths for source and destination because not doing so
	// can result in some weird symbolic link inception
	source, err := filepath.Abs(AssetsDir)
	if err != nil {
		return err
	}

	dest, err := filepath.Abs(TargetDir + AssetsDir)
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
	if verbose {
		fmt.Printf("Rendered article '%v'\n", articleFile)
	}

	source, err := ioutil.ReadFile(ArticlesDir + articleFile)
	if err != nil {
		return err
	}
	rendered := renderMarkdown(source)

	template, err := ace.Load(LayoutsDir+"main", LayoutsDir+"article", nil)
	if err != nil {
		return err
	}

	file, err := os.Create(TargetDir + trimExtension(articleFile))
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
	if verbose {
		fmt.Printf("Rendered page '%v'\n", pageFile)
	}

	template, err := ace.Load(LayoutsDir+"main", PagesDir+pageFile, nil)
	if err != nil {
		return err
	}

	file, err := os.Create(TargetDir + pageFile)
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
