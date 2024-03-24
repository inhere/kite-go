package httpcmd

import (
	"strings"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/rux"
)

// https://github.com/swagger-api/swagger-ui
// html from https://swagger.io/docs/open-source-tools/swagger-ui/usage/installation/
var oapiUISwagger = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta
      name="description"
      content="SwaggerUI"
    />
    <title>Open API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui.css" />
  </head>
  <body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-bundle.js" crossorigin></script>
  <script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-standalone-preset.js" crossorigin></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: '{{docUrl}}',
        dom_id: '#swagger-ui',
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        layout: "StandaloneLayout",
      });
    };
  </script>
  </body>
</html>
`

// github: https://github.com/Redocly/redoc
// from https://redocly.com/docs/redoc/quickstart/
var oapiUiRedoc = `
<!DOCTYPE html>
<html>
  <head>
    <title>Open API Docs</title>
    <!-- needed for adaptive design -->
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet" />

    <!-- Redoc doesn't change outer page styles -->
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <!-- Redoc element with link to your OpenAPI definition -->
    <redoc spec-url="{{docUrl}}"></redoc>
    <!-- Link to Redoc JavaScript on CDN for rendering standalone element -->
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
  </body>
</html>
`

// use https://github.com/stoplightio/elements
var oapiUiElements = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>Open API Documents</title>
    <!-- Embed elements Elements via Web Component -->
    <script src="https://unpkg.com/@stoplight/elements/web-components.min.js"></script>
    <link rel="stylesheet" href="https://unpkg.com/@stoplight/elements/styles.min.css">
    <style>
        body {
            display: flex;
            flex-direction: column;
            height: 100vh;
        }
        .api-container {
            flex: 1 0 0;
            overflow: hidden;
        }
        .sl-elements {
            font-size: 14px;
        }
        .sl-text-base {
            font-size: 14px;
        }
        div.sl-py-16 {
            max-width: none !important;
        }
        div.sl-relative.sl-ml-16 {
            max-width: none !important;
        }
        .sl-elements .TextRequestBody {
            max-height: 45%;
        }
    </style>
</head>
<body>
<header>
    <div class="sl-inverted sl-flex sl-flex-shrink-0 sl-items-center sl-h-2xl sl-px-5 sl-bg-canvas-pure">
        <div class="sl-stack sl-stack--horizontal sl-stack--4 sl-flex sl-flex-row sl-items-center sl-w-1/3">
            <span class="sl-text-lg sl-leading-none sl-font-semibold">OpenAPI Docs Browser</span></div>
        <div class="sl-flex sl-justify-center sl-w-1/3">
            <div class="sl-stack sl-stack--horizontal sl-stack--2 sl-flex sl-flex-1 sl-flex-row sl-items-center">
                <div class="sl-input sl-flex-1 sl-relative sl-bg-canvas-100 sl-rounded">
                    <label for="doc-url-input" style="display: none">Doc URL</label>
                    <input placeholder="URL to an OpenAPI document..." type="text" id="doc-url-input"
                           class="sl-relative sl-w-full sl-h-md sl-text-base sl-pr-2.5 sl-pl-2.5 sl-rounded sl-border-transparent hover:sl-border-input focus:sl-border-primary sl-border"
                           value="{{docUrl}}"></div>
                <button type="button" id="try-it-button"
                        class="sl-button sl-h-md sl-text-base sl-font-medium sl-px-2.5 sl-bg-canvas hover:sl-bg-canvas-50 active:sl-bg-canvas-100 sl-rounded sl-border-button sl-border disabled:sl-opacity-60">
                    Try It!
                </button>
            </div>
        </div>
    </div>
</header>
<div class="api-container">
    <elements-api id="api-docs"
                  apiDescriptionUrl="{{docUrl}}"
                  router="hash"
                  layout="sidebar"
    />
</div>
</body>
<script>
    (async () => {
        const docsBox = document.getElementById('api-docs');
        const tryItBtn = document.getElementById('try-it-button');
        const urlInput = document.getElementById('doc-url-input');

        tryItBtn.addEventListener('click', async () => {
            let newUrl = urlInput.value.trim();
            if (newUrl !== '' && newUrl !== docsBox.apiDescriptionUrl) {
                // set url will cause a blank flash. 会导致空白闪烁
                // docsBox.apiDescriptionUrl = newUrl;
                docsBox.apiDescriptionDocument = await fetch(newUrl).then(res => res.text());
            }
        });
    })();
</script>
</html>
`

// NewOAPIServeCmd instance
func NewOAPIServeCmd() *gcli.Command {
	var esOpts = struct {
		port   uint
		style  string
		uiHtml string
	}{}

	return &gcli.Command{
		Name:    "oapi-serve",
		Desc:    "start an open api doc UI http server",
		Aliases: []string{"swag-ui", "oapi-ui", "swagger"},
		Config: func(c *gcli.Command) {
			c.UintOpt(&esOpts.port, "port", "P", 0, "custom the echo server port, default will use random `port`")
			c.StrOpt2(&esOpts.style, "style,s", "the openapi doc UI style, support: redoc, elements, swagger")
			c.StrOpt2(&esOpts.uiHtml, "ui-html,ui", "the custom UI html file for render openapi doc")

			c.AddArg("docPath", "the open api doc file path, support: local file path, remote url", true)
		},
		Help: `
Use custom UI html:
  {$fullCmd} --ui-html static/oapi-ui-custom.html static/swagger.json

Some Example APIs:
- https://petstore.swagger.io/v2/swagger.json
- https://petstore3.swagger.io/api/v3/openapi.json
`,
		Func: func(c *gcli.Command, args []string) error {
			if esOpts.port < 1 {
				esOpts.port = uint(mathutil.RandInt(6000, 19999))
			}

			docPath := c.Arg("docPath").String()

			srv := rux.New()
			if fsutil.IsFile(docPath) {
				filePath := docPath
				docPath = fsutil.Name(docPath)
				srv.StaticFile(docPath, filePath)
			}

			srv.GET("/", func(c *rux.Context) {
				var htmlStr string
				switch esOpts.style {
				case "redoc":
					htmlStr = oapiUiRedoc
				case "elem", "element", "elements":
					htmlStr = oapiUiElements
				default:
					if esOpts.uiHtml != "" {
						htmlStr = fsutil.ReadString(esOpts.uiHtml)
					} else { // swagger ui
						htmlStr = oapiUISwagger
					}
				}

				c.HTMLString(200, strings.Replace(htmlStr, "{{docUrl}}", docPath, -1))
			})

			srv.Listen("127.0.0.1", mathutil.String(esOpts.port))
			return srv.Err()
		},
	}
}
