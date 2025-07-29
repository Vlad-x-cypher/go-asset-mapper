# Go Asset Mapper

go-asset-mapper is a Go package to serve versioned static files in Go HTML Templates. Main feature: map and version all static files in specific directory.

## Instalation

Install via Go modules:
```bash
go get github.com/Vlad-x-cypher/go-asset-mapper
```

### Example 

```go
 
package main

import (
	"html/template"
	"log"
    "os"

	"github.com/Vlad-x-cypher/go-asset-mapper"
)

func main() {
	t := template.New("")

	assetMapper := asset.NewAssetMapper()

    // Scan specific directory to map all files
	err := assetMapper.ScanDir("assets")
	if err != nil {
		log.Fatalf("assets scandir err: %v", err)
	}

    // Add functions to use inside templates
	t.Funcs(template.FuncMap{
		"scriptTag": assetMapper.ScriptTag,
		"linkTag":   assetMapper.LinkTag,
		"asset":     assetMapper.Get,
	})

    const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
        {{ linkTag "assets/style.css" }}
        {{ scriptTag "assets/main.js" }}
	</head>
	<body>
        <img src="{{ asset "assets/example.png" }}" />
	</body>
</html>`

	templates, err := t.Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}

    err = templates.Execute(os.Stdout, nil)
    if err != nil {
        log.Fatal(err)
    }
}
```
Output:
```html
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
        <link href="assets/style.css" rel="stylesheet"/>
        <script src="assets/main.js"></script>
	</head>
	<body>
        <img src="assets/example.png" />
	</body>
</html>
```

For more complete example see [example dir](./example)

