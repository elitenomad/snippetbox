package main

import (
	"github.com/elitenomad/snippetbox/pkg/models"
	"html/template"
	"path/filepath"
)

type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	/*
		Initialize a new map to act as a cache
	 */
	cache := map[string]*template.Template{}

	/*
		Use the filePath.Glob function to get a slice of all filepaths with the
		extension *.page.tmpl
	 */
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	/*
		Loop through the pages one by one
	 */
	for _, page := range pages {
		/*
			Extract the full file name and assign it to a name variable
		 */
		name := filepath.Base(page)

		/*
			Parse the file to the templateSet
		 */
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		/*
			Add Layout templates to Template set
		 */

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		/*
			Add partial templates to Template set
		 */
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		/*
			Add Template set to the cache
		 */

		cache[name] = ts
	}

	return cache, nil
}