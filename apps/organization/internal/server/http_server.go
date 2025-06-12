package server

import (
	"log"
	"net/http"

	credHttp "github.com/DO-2K23-26/polypass-microservices/organization/internal/transport/http"
	folderHttp "github.com/DO-2K23-26/polypass-microservices/organization/internal/transport/http"
	tagHttp "github.com/DO-2K23-26/polypass-microservices/organization/internal/transport/http"
	"github.com/gorilla/mux"
)

type HttpServer struct {
	port                    string
	folderHandler           *folderHttp.FolderHandler
	tagHandler              *tagHttp.TagHandler
	folderCredentialHandler *credHttp.FolderCredentialHandler
}

func NewHttpServer(port string, folderHandler *folderHttp.FolderHandler, tagHandler *tagHttp.TagHandler, credHandler *credHttp.FolderCredentialHandler) *HttpServer {
	return &HttpServer{
		port:                    port,
		folderHandler:           folderHandler,
		tagHandler:              tagHandler,
		folderCredentialHandler: credHandler,
	}
}

func (s *HttpServer) Start() {
	r := mux.NewRouter()

	// Folders
	r.HandleFunc("/folders", s.folderHandler.CreateFolder).Methods("POST")
	r.HandleFunc("/folders", s.folderHandler.ListFolders).Methods("GET")
	r.HandleFunc("/folders/{id}", s.folderHandler.GetFolder).Methods("GET")
	r.HandleFunc("/folders/{id}", s.folderHandler.UpdateFolder).Methods("PUT")
	r.HandleFunc("/folders/{id}", s.folderHandler.DeleteFolder).Methods("DELETE")
	r.HandleFunc("/folders/{id}/users", s.folderHandler.ListUsersInFolder).Methods("GET")

	// Folder credentials
	r.HandleFunc("/folders/{folderId}/credentials/{type}", s.folderCredentialHandler.ListCredentials).Methods("GET")
	r.HandleFunc("/folders/{folderId}/credentials/{type}", s.folderCredentialHandler.CreateCredential).Methods("POST")
	r.HandleFunc("/folders/{folderId}/credentials/{type}/{credentialId}", s.folderCredentialHandler.UpdateCredential).Methods("PUT")
	r.HandleFunc("/folders/{folderId}/credentials/{type}", s.folderCredentialHandler.DeleteCredentials).Methods("DELETE")

	r.HandleFunc("/users/credentials", s.folderCredentialHandler.ListUserCredentials).Methods("GET")

	// Tags
	r.HandleFunc("/tags", s.tagHandler.CreateTag).Methods("POST")
	r.HandleFunc("/tags", s.tagHandler.ListTags).Methods("GET")
	r.HandleFunc("/tags/{id}", s.tagHandler.GetTag).Methods("GET")
	r.HandleFunc("/tags/{id}", s.tagHandler.UpdateTag).Methods("PUT")
	r.HandleFunc("/tags/{id}", s.tagHandler.DeleteTag).Methods("DELETE")

	log.Printf("Starting HTTP server on %s", s.port)
	log.Fatal(http.ListenAndServe(s.port, r))
}
