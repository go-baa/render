package render

import (
	"os"

	"github.com/fsnotify/fsnotify"
)

const (
	Create fsnotify.Op = 1 << iota
	Write
	Remove
	Rename
	Chmod
)

type notifyItem struct {
	event fsnotify.Op
	path  string
}

func (r *Render) notify() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		r.Error(err)
		return
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					r.fileChanges <- notifyItem{fsnotify.Write, event.Name}
				} else if event.Op&fsnotify.Create == fsnotify.Create {
					r.fileChanges <- notifyItem{fsnotify.Create, event.Name}
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					r.fileChanges <- notifyItem{fsnotify.Remove, event.Name}
				}
			case err = <-watcher.Errors:
				r.Error(err)
			}
		}
	}()

	var l []string
	l = append(l, r.Root)
	err = recursiveDir(r.Root, &l)
	if err != nil {
		r.Error(err)
		return
	}
	for _, d := range l {
		err = watcher.Add(d)
		if err != nil {
			r.Error(err)
			return
		}
	}

	<-done
}

func recursiveDir(dir string, l *[]string) error {
	dl, err := readDir(dir)
	if err != nil {
		return err
	}
	for _, d := range dl {
		if d.IsDir() {
			_dir := dir + "/" + d.Name()
			*l = append(*l, _dir)
			err = recursiveDir(_dir, l)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func readDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	return list, nil
}
