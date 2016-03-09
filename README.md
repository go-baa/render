# render
a powerful template engine than default render of baa.

## Features

- template collector, let you can include an other template in template.
- template cache, cache template file in memory, detect file change and rebuild cache.
- custom template dir, extensions and functions.

## Getting Started

```
package main

import (
    "github.com/go-baa/baa"
    "github.com/go-baa/render"
)

func main() {
    app := baa.New()
    app.SetDI("render", render.New(render.Options{
		Baa:        app,
		Root:       "templates/",
		Extensions: []string{".html", ".tmpl"},
	}))
    app.Get("/", func(c *baa.Context) {
        c.HTML(200, "index")
    })
    app.Run(":1323")
}
```

you should first have a dir named ``templates`` in your application root and a template file named ``index.html`` in template dir.
