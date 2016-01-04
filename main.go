package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday"
	"github.com/yosssi/ace"
)

const (
	ArticlesDir = "./articles/"
	AssetsDir   = "./assets/"
	LayoutDir   = "./layout/"
	PagesDir    = "./pages/"
	TargetDir   = "./public/"
)

func main() {
	// create an output directory
	err := os.MkdirAll(TargetDir, 0755)
	if err != nil {
		panic(err)
	}

	err = linkAssets()
	if err != nil {
		panic(err)
	}

	err = renderArticles()
	if err != nil {
		panic(err)
	}

	err = renderPages()
	if err != nil {
		panic(err)
	}
}

func linkAssets() error {
	err := os.RemoveAll(TargetDir + "assets")
	if err != nil {
		return err
	}

	source, err := filepath.Abs("./assets")
	if err != nil {
		return err
	}

	dest, err := filepath.Abs(TargetDir + "assets")
	if err != nil {
		return err
	}

	err = os.Symlink(source, dest)
	if err != nil {
		return err
	}

	return nil
}

func renderArticles() error {
	files, err := ioutil.ReadDir(ArticlesDir)
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		articleFile := fileInfo.Name()

		source, err := ioutil.ReadFile(ArticlesDir + articleFile)
		if err != nil {
			return err
		}
		rendered := renderMarkdown(source)

		template, err := ace.Load(LayoutDir+"main", LayoutDir+"article", nil)
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

func renderPages() error {
	files, err := ioutil.ReadDir(PagesDir)
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		// ace.Load doesn't like to take .ace extensions anyway
		pageFile := trimExtension(fileInfo.Name())

		template, err := ace.Load(LayoutDir+"main", PagesDir+pageFile, nil)
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
	}

	return nil
}

func trimExtension(file string) string {
	return strings.TrimSuffix(file, filepath.Ext(file))
}
