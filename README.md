# render
a template engine better than default render of baa.

## Features

- template collector, let you can include an other template in template.
- template cache, cache template file in memory, detect file change and rebuild cache.
- custom template dir, extensions and functions.

## Getting Started

```
package main

import (
    "github.com/go-baa/render"
    "gopkg.in/baa.v1"
)

func main() {
    // new app
    app := baa.New()
    
    // register render
    // render is template DI for baa, must this name.
    app.SetDI("render", render.New(render.Options{
		Baa:        app,
		Root:       "templates/",
		Extensions: []string{".html", ".tmpl"},
	}))
    
    // router
    app.Get("/", func(c *baa.Context) {
        c.HTML(200, "index")
    })
    
    // run app
    app.Run(":1323")
}
```

you should first have a dir named ``templates`` in your application root and a template file named ``index.html`` in template dir.
