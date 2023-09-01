package controller

import "github.com/gookit/rux"

// JSONPage json page controller
type JSONPage struct {
}

// AddRoutes to rux.Router
func (p *JSONPage) AddRoutes(r *rux.Router) {
	r.GET("/", p.View)
	r.GET("/view", p.View)
}

var jsonHtml = `<!DOCTYPE HTML>
<html lang="en">
<head>
    <!-- when using the mode "code", it's important to specify charset utf-8 -->
    <meta charset="utf-8">
  	<title>JSONEditor | View,Query and more</title>
    <link href="https://cdn.jsdelivr.net/npm/jsoneditor@9.10.2/dist/jsoneditor.min.css" rel="stylesheet">
    <script src="https://cdn.jsdelivr.net/npm/jsoneditor@9.10.2/dist/jsoneditor.min.js"></script>
  <style type="text/css">
    body {
      color: #4d4d4d;
      line-height: 150%;
      display: flex;
      flex-direction: column;
      height: 100vh;
      margin: 0;
      padding: 0;
    }
    code {
      background-color: #f5f5f5;
    }
    #jsoneditor {
      width: 100%;
      /* height: 95%; */
      flex: 1 0 0;
    }
	.jsoneditor-menu {
      padding-left: 42%;
      height: 2.6em;
      padding-top: 5px;
	}
  </style>
</head>
<body>
    <div id="jsoneditor"></div>
    <script>
      // create the editor
      const container = document.getElementById("jsoneditor")
	  const options = {
		mode: 'preview',
		modes: ['code', 'form', 'text', 'tree', 'view', 'preview'], // allowed modes
		onModeChange: function (newMode, oldMode) {
		  console.log('Mode switched from', oldMode, 'to', newMode)
		}
	  }
      const editor = new JSONEditor(container, options)

        // set json
        const initialJson = {
            "Array": [1, 2, 3],
            "Boolean": true,
            "Null": null,
            "Number": 123,
            "Object": {"a": "b", "c": "d"},
            "String": "Hello World"
        }
        editor.set(initialJson)

        // get json
        const updatedJson = editor.get()
    </script>
</body>
</html>`

// View page
func (p *JSONPage) View(c *rux.Context) {
	c.HTMLString(200, jsonHtml)
}
