package swagger

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/tbe-team/raybot/api/openapi"
)

// Register registers the swagger handler on the given router
func Register(r chi.Router) {
	specPath := "/docs/openapi.json"
	template := getTemplate(specPath)

	r.Get("/docs", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		//nolint:errcheck
		w.Write([]byte(template))
	})

	specBytes := openapi.GetSpecBytes()
	r.Get(specPath, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//nolint:errcheck
		w.Write(specBytes)
	})
}

// getTemplate returns the HTML template for Swagger UI
func getTemplate(specPath string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <meta name="description" content="SwaggerUI" />
  <title>SwaggerUI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
  <link rel="icon" type="image/png" href="https://static1.smartbear.co/swagger/media/assets/swagger_fav.png" sizes="32x32" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
<script>
  window.onload = () => {
    window.ui = SwaggerUIBundle({
      url: '%s',
      dom_id: '#swagger-ui',
      deepLinking: true,
	  showExtensions: true,
	  showCommonExtensions: true,
    });
  };
</script>
</body>
</html>
`, specPath)
}
