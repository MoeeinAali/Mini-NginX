package handlers

import (
	"bufio"
	"net"
)

func ReverseProxy(client net.Conn, reader *bufio.Reader, target string, firstLine string) {}
