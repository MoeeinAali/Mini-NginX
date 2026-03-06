package router

import (
	"Mini-NginX/internal/config"
	"Mini-NginX/internal/handlers"
	"Mini-NginX/internal/logger"
	"bufio"
	"fmt"
	"net"
	"strings"
)

func Handle(conn net.Conn, cfg *config.Config) {
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	parts := strings.Split(line, " ")
	method := parts[0]
	path := parts[1]
	logger.Info(fmt.Sprintf("request method: %s, path: %s", method, path))
	for route, rule := range cfg.Paths {
		if strings.HasPrefix(path, route) {
			logger.Info("route matched", "route", route, "type", rule.Type)

			switch rule.Type {
			case "redirect":
				handlers.Redirect(conn, rule.Target)
			case "staticfile":
				handlers.StaticFile(conn, rule.Target)
			case "servedir":
				handlers.ServeDir(conn, rule.Target, route, path)
			case "reverseproxy":
				handlers.ReverseProxy(conn, reader, rule.Target, line)
			default:
				write404(conn)
			}
		}
	}
}

func write404(conn net.Conn) {
	body := "404 not found"
	resp := fmt.Sprintf("HTTP/1.1 404 Not Found\r\nContent-Length: %d\r\n\r\n", len(body), body)
	conn.Write([]byte(resp))
}
