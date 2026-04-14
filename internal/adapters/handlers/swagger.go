package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

func ServeSwaggerUI(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/swagger/index.html" || r.URL.Path == "/swagger/" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
  <title>Transactions Service API Documentation</title>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.0/swagger-ui.min.css">
  <style>
    html {
      box-sizing: border-box;
      overflow: -moz-scrollbars-vertical;
      overflow-y: scroll;
    }
    *, *:before, *:after {
      box-sizing: inherit;
    }
    body {
      margin: 0;
      padding: 0;
    }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.0/swagger-ui.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.0/swagger-ui-bundle.min.js"></script>
  <script>
    window.onload = function() {
      SwaggerUIBundle({
        url: "/swagger.json",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIBundle.SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
      });
    };
  </script>
</body>
</html>
`))
		return
	}

	if r.URL.Path == "/swagger.json" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		swaggerPath := filepath.Join(os.Getenv("PWD"), "docs", "swagger.json")
		if swaggerPath == "" {
			swaggerPath = "./docs/swagger.json"
		}
		data, err := os.ReadFile(swaggerPath)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"openapi": "3.0.0",
				"info": map[string]any{
					"title":   "Transactions Service API",
					"version": "1.0.0",
				},
			})
			return
		}
		w.Write(data)
		return
	}

	http.NotFound(w, r)
}
