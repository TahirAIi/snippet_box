package main

import (
	"html/template"
	"path/filepath"
	"snippet_box/pkg/forms"
	"snippet_box/pkg/models"
)

type templateData struct {
	Form     *forms.Form
	Snippet  *models.Snippets
	Snippets []*models.Snippets
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}
	for _, page := range pages {

		name := filepath.Base(page)

		templateToCache, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		templateToCache, err = templateToCache.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		templateToCache, err = templateToCache.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		cache[name] = templateToCache
	}
	return cache, nil
}
