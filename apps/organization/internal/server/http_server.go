package server

import (
    "log"
    "net/http"

    folderHttp "github.com/DO-2K23-26/polypass-microservices/organization/internal/ports/http"
    tagHttp "github.com/DO-2K23-26/polypass-microservices/organization/internal/ports/http"
)

type HttpServer struct {
    port          string
    folderHandler *folderHttp.FolderHandler
    tagHandler    *tagHttp.TagHandler
}

func NewHttpServer(port string, folderHandler *folderHttp.FolderHandler, tagHandler *tagHttp.TagHandler) *HttpServer {
    return &HttpServer{
        port:          port,
        folderHandler: folderHandler,
        tagHandler:    tagHandler,
    }
}

func (s *HttpServer) Start() {
    http.HandleFunc("/folders", s.folderHandler.CreateFolder)
    http.HandleFunc("/tags", s.tagHandler.CreateTag)

    log.Printf("Starting HTTP server on %s", s.port)
    log.Fatal(http.ListenAndServe(s.port, nil))
}
