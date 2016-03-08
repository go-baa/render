// Package render providers a template engine for baa.
package render

import (
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-baa/baa"
)

// Render a powerful template engine than default render of baa
type Render struct {
	Options
	template    *template.Template // template handle
	fileChanges chan notifyItem    // notify file changes
}

// Options render options
type Options struct {
	Baa        *baa.Baa         // baa
	Root       string           // template root dir
	Extensions []string         // template file extensions
	FuncMap    template.FuncMap // template functions
}

// New create a template engine
func New(o Options) *Render {
	r := new(Render)
	r.Baa = o.Baa
	r.Root = o.Root
	r.Extensions = o.Extensions
	r.FuncMap = o.FuncMap

	// check template dir
	if r.Root == "" {
		panic("Render template dir is empty!")
	}
	r.Root, _ = filepath.Abs(r.Root)
	if r.Root[len(r.Root)-1] != '/' {
		r.Root += "/" // add right slash
	}
	if f, err := os.Stat(r.Root); err != nil {
		panic("Render template dir[" + r.Root + "] open error: " + err.Error())
	} else {
		if !f.IsDir() {
			panic("Render template dir[" + r.Root + "] is not s directory!")
		}
	}

	// check extension
	if r.Extensions == nil {
		r.Extensions = []string{".html"}
	}

	// set template
	r.template = template.New("_DEFAULT_")
	r.template.Funcs(r.FuncMap)

	// load templates
	r.loadTpls()

	// notify
	r.fileChanges = make(chan notifyItem, 32)
	go r.notify()
	go func() {
		for item := range r.fileChanges {
			if r.Baa != nil && r.Baa.Debug() {
				r.Error("filechanges Receive -> " + item.path)
			}
			if item.event == Create || item.event == Write {
				r.parseFile(item.path)
			}
		}
	}()

	return r
}

// Render template
func (r *Render) Render(w io.Writer, tpl string, data interface{}) error {
	return r.template.ExecuteTemplate(w, r.tplName(tpl), data)
}

// loadTpls load all template files
func (r *Render) loadTpls() {
	paths, err := r.readDir(r.Root)
	if err != nil {
		r.Error(err)
		return
	}
	for _, path := range paths {
		err = r.parseFile(path)
		if err != nil {
			r.Error(err)
		}

	}
}

// readDir scan dir load all template files
func (r *Render) readDir(path string) ([]string, error) {
	var paths []string
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fs, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var p string
	for _, f := range fs {
		p = filepath.Clean(path + "/" + f.Name())
		if f.IsDir() {
			fs, err := r.readDir(p)
			if err != nil {
				continue
			}
			for _, f := range fs {
				paths = append(paths, f)
			}
		} else {
			if r.checkExt(p) {
				paths = append(paths, p)
			}
		}
	}
	return paths, nil
}

// tplName get template alias from a template file path
func (r *Render) tplName(path string) string {
	if len(path) > len(r.Root) && path[:len(r.Root)] == r.Root {
		path = path[len(r.Root):]
	}
	ext := filepath.Ext(path)
	return path[:len(path)-len(ext)]
}

// checkExt check path extension allow use
func (r *Render) checkExt(path string) bool {
	ext := filepath.Ext(path)
	if ext == "" {
		return false
	}
	for i := range r.Extensions {
		if r.Extensions[i] == ext {
			return true
		}
	}
	return false
}

// parseFile load file and parse to template
func (r *Render) parseFile(path string) error {
	if r.Baa != nil && r.Baa.Debug() {
		r.Error("loadTpl -> " + path)
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	s := string(b)
	t := r.template.New(r.tplName(path))
	_, err = t.Parse(s)
	if err != nil {
		return err
	}
	return nil
}

// Error log error
func (r *Render) Error(v interface{}) {
	if r.Baa != nil {
		r.Baa.Logger().Println(v)
	}
}
