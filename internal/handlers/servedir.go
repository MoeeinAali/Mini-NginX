package handlers

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func ServeDir(conn net.Conn, dir string, route string, requestPath string) {
	rel := strings.TrimPrefix(requestPath, route)
	fsPath := filepath.Join(dir, rel)

	info, err := os.Stat(fsPath)
	if err != nil {
		write404(conn)
		return
	}

	if !info.IsDir() {
		StaticFile(conn, fsPath)
		return
	}

	entries, err := os.ReadDir(fsPath)

	if err != nil {
		body := "cannot read directory"
		resp := fmt.Sprintf("HTTP/1.1 500 Internal Server Erorr\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
		conn.Write([]byte(resp))
		return
	}

	var body strings.Builder

	body.WriteString("<html><body><h1>Directory listing</h1><ul>")

	for _, e := range entries {
		name := e.Name()
		link := requestPath + name

		if e.IsDir() {
			link = link + "/"
			name = name + "/"
		}

		body.WriteString(fmt.Sprintf(`<li><a href="%s">%s</a></li>`, link, name))
	}
	body.WriteString("</ul></body></html>")

	resp := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: text/html\r\n\r\n%s", len(body.String()), body.String())
	conn.Write([]byte(resp))
}

func write404(conn net.Conn) {
	body := "404 not found"
	resp := fmt.Sprintf("HTTP/1.1 404 Not Found\r\nContent-Length: %d\r\n\r\n", len(body), body)
	conn.Write([]byte(resp))
}
