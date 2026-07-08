package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// CustomSwaggerHandler serves a custom Swagger UI HTML that properly uses current window location
func (h *Handler) CustomSwaggerHandler(c echo.Context) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin:0; padding:0; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
    <script>
    window.onload = function() {
        const specUrl = window.location.protocol + '//' + window.location.host + '/swagger/doc.json';
        
        // Fetch and modify the spec to remove hardcoded host
        fetch(specUrl)
            .then(response => response.json())
            .then(spec => {
                // Remove host and schemes to force using current location
                delete spec.host;
                delete spec.schemes;
                
                window.ui = SwaggerUIBundle({
                    spec: spec,
                    dom_id: '#swagger-ui',
                    deepLinking: true,
                    presets: [
                        SwaggerUIBundle.presets.apis,
                        SwaggerUIStandalonePreset
                    ],
                    plugins: [
                        SwaggerUIBundle.plugins.DownloadUrl
                    ],
                    layout: "StandaloneLayout"
                });
            });
    };
    </script>
</body>
</html>`

	return c.HTML(http.StatusOK, html)
}
