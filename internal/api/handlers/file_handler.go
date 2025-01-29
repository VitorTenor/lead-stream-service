package handlers

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/vitortenor/lead-stream-service/internal/domain"
	"github.com/vitortenor/lead-stream-service/internal/services"
)

func InitFileRoutes(humaApi huma.API, fileHandler *FileHandler) {
	huma.Register(humaApi, huma.Operation{
		Path:          "/schema/{schemaId}/file",
		OperationID:   "upload-file",
		Method:        http.MethodPost,
		DefaultStatus: http.StatusOK,
		Summary:       "Upload a file",
		Description:   "Upload a file to the given schema",
	}, fileHandler.Upload)
}

type FileHandler struct {
	service *services.FileService
}

func NewFileHandler(service *services.FileService) *FileHandler {
	return &FileHandler{
		service: service,
	}
}

func (fh *FileHandler) Upload(ctx context.Context, fr *FileRequest) (*FileResponse, error) {
	err := fh.service.ProcessAndSave(&ctx, fr.toDomain())
	if err != nil {
		return nil, handleError(err)
	}

	response := &FileResponse{}
	response.Body.Message = "File uploaded successfully"
	return response, nil
}

type FileRequest struct {
	SchemaId string `path:"schemaId" required:"true"`
	RawBody  multipart.Form
}

func (fr *FileRequest) toDomain() *domain.File {
	return &domain.File{
		SchemaId: fr.SchemaId,
		File:     fr.RawBody.File["file"][0],
	}
}

type FileResponse struct {
	Body struct {
		Message string `json:"message" description:"The message of the response"`
	}
}
