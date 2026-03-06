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

	breadcrumb := "/" + strings.Trim(requestPath, "/")

	var entriesHTML strings.Builder
	if len(entries) == 0 {
		entriesHTML.WriteString(`<div class="empty-state"><p>📭 This directory is empty</p></div>`)
	} else {
		entriesHTML.WriteString(`<ul class="entries-list">`)
		for _, e := range entries {
			name := e.Name()
			link := requestPath + name
			var entryType, icon string

			if e.IsDir() {
				link = link + "/"
				entryType = "Folder"
				icon = "📁"
			} else {
				entryType = "File"
				icon = "📄"
			}

			entriesHTML.WriteString(fmt.Sprintf(`
					<li><a href="%s" class="entry-item">
                    <span class="entry-icon">%s</span>
                    <span class="entry-name">%s</span>
                    <span class="entry-type">%s</span></a></li>`, link, icon, escapeHTML(name), entryType))
		}
		entriesHTML.WriteString(`</ul>`)
	}

	templatePath := "static/dirlist.html"
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		body := "error loading template"
		resp := fmt.Sprintf("HTTP/1.1 500 Internal Server Error\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
		conn.Write([]byte(resp))
		return
	}

	html := string(templateData)
	html = strings.ReplaceAll(html, "{{BREADCRUMB}}", escapeHTML(breadcrumb))
	html = strings.ReplaceAll(html, "{{ENTRIES}}", entriesHTML.String())

	resp := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: text/html; charset=utf-8\r\n\r\n%s", len(html), html)
	conn.Write([]byte(resp))
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

func write404(conn net.Conn) {
	body := "404 not found"
	resp := fmt.Sprintf("HTTP/1.1 404 Not Found\r\nContent-Length: %d\r\n\r\n", len(body), body)
	conn.Write([]byte(resp))
}
