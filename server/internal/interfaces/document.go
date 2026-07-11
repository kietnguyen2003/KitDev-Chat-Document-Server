package interfaces

import (
	"fmt"
	"net/http"
	"server/internal/application/document"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	docService document.DocumentService
}

func NewDocumentHandler(docService *document.DocumentService) *DocumentHandler {
	return &DocumentHandler{docService: *docService}
}

func (dh *DocumentHandler) UploadFile(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("looix ddaay nef")
		errorResponse(c, http.StatusBadGateway, err.Error())
		return
	}
	cateName := c.PostForm("cateName")
	username := c.GetHeader("KIT-DEV-USERNAME")

	src, _ := file.Open()
	defer src.Close()

	objectName, err := dh.docService.UploadFile(c.Request.Context(), document.UploadDocDTO{
		Username:    username,
		Category:    cateName,
		Bucket:      "documents",
		DocName:     file.Filename,
		Reader:      src,
		Size:        file.Size,
		ContentType: file.Header.Get("Content-Type"),
	})

	if err != nil {
		println(err.Error())
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	successResponse(c, http.StatusAccepted, "Uploadfile Success", UploadFileRes{ObjectName: objectName})

}

func (dh *DocumentHandler) GetListFileByUserID(c *gin.Context) {
	cateName := c.Param("cateName")
	username := c.GetHeader("KIT-DEV-USERNAME")

	res, err := dh.docService.GetDocumentList(c.Request.Context(), username, cateName)

	if err != nil {
		errorResponse(c, http.StatusBadGateway, err.Error())
	}

	var data []GetDocumentListRes
	for _, doc := range res {
		data = append(data, GetDocumentListRes{
			Doc_ID:     doc.ID,
			UserID:     doc.UserID,
			Name:       doc.Name,
			ObjectName: doc.ObjectName,
			Size:       doc.Size,
			Status:     doc.Status,
		})
	}

	successResponse(c, http.StatusAccepted, "Get Document List sucess", data)
}

type DeleteDocumentReq struct {
	Object string `json:"object_name"`
	Size   uint64 `json:"size"`
}

type UpdateStatus struct {
	Status string `json:"status"`
}

func (dh *DocumentHandler) UpdateStatus(c *gin.Context) {
	var req UpdateStatus

	docID := c.Param("id")

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Cant not bind request")
		return
	}

	username := c.GetHeader("KIT-DEV-USERNAME")

	err := dh.docService.UpdateStatus(c.Request.Context(), username, utils.StringToUint(docID), req.Status)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	successResponse(c, http.StatusCreated, "Upadte Status success", nil)
}

func (dh *DocumentHandler) DeleteDocument(c *gin.Context) {
	var req DeleteDocumentReq

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Cant not bind request")
		return
	}

	username := c.GetHeader("KIT-DEV-USERNAME")

	deleteDocument, err := dh.docService.DeleteDocument(c.Request.Context(), username, req.Object, req.Size)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if deleteDocument == false {
		errorResponse(c, http.StatusBadRequest, "Delete document fail")
		return
	}
	successResponse(c, http.StatusCreated, "Delete document success", nil)
}
