package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tsironi93/miniHttp/internal/request"
	"github.com/tsironi93/miniHttp/internal/response"
	"github.com/tsironi93/miniHttp/internal/server"
)

const port = 42069

func loadHtml(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	return string(data)
}

func htmlHandler(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.StatusCode = response.StatusBadRequest
		w.WriteString(loadHtml("./internal/htmlTemplates/400.html"))
	case "/myproblem":
		w.StatusCode = response.StatusInternalServerError
		w.WriteString(loadHtml("./internal/htmlTemplates/500.html"))
	default:
		w.StatusCode = response.StatusOK
		w.WriteString(loadHtml("./internal/htmlTemplates/200.html"))
	}
}

func main() {
	server, err := server.Serve(port, htmlHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
