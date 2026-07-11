package infracstructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	docuemntindexer "server/internal/application/docuemnt_indexer"
)

type PythonIndexClientRepository struct {
	baseUrl string
}

func NewPythonIndexClientRepository(url string) *PythonIndexClientRepository {
	return &PythonIndexClientRepository{
		baseUrl: url,
	}
}

func (pr *PythonIndexClientRepository) IndexDocument(ctx context.Context, req docuemntindexer.IndexDocumentRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		pr.baseUrl+"/api/documents",
		bytes.NewReader(body), // chuyen bytes thanh dang reader
	)
	if err != nil {
		return err
	}

	// Set header
	httpReq.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return fmt.Errorf("Request to rag server failed")
	}

	return nil
}

func (pr *PythonIndexClientRepository) DeleteDocument(ctx context.Context, req docuemntindexer.DeleteDocumentRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		pr.baseUrl+"/api/documents",
		bytes.NewReader(body), // chuyen bytes thanh dang reader
	)
	if err != nil {
		return err
	}

	// Set header
	httpReq.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {

		return err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return fmt.Errorf("Request to rag server failed")
	}

	return nil
}
