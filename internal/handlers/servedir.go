package handlers

import (
	"fmt"
	"html/template"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func ServeDir(conn net.Conn, dir string) {
	files := []string{}
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(dir, path)
		if rel != "." {
			if info.IsDir() {
				rel += "/"
			}
			files = append(files, rel)
		}
		return nil
	})

	htmlTemplate := `
					<html>
						<head>
							<title>Index of {{.Dir}}</title>
						</head>
						<body>
							<h1>Index of {{.Dir}}</h1>
							<ul>
								{{range .Files}}
									<li>
										<a href="{{.}}">{{.}}</a>
									</li>
								{{end}}
							</ul>
						</body>
					</html>`

	templ := template.Must(template.New("dirlist").Parse(htmlTemplate))
	var body strings.Builder
	templ.Execute(&body, map[string]interface{}{"Dir": dir, "Files": files})

	resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: text/html\r\n\r\n%s", len(body.String()), body.String())
	conn.Write([]byte(resp))

}
