package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/aivanovanynines/ping/internal/models"
)

type fetcherUseCase struct {
	httpClient *http.Client
}

func NewFetcherUsecase() *fetcherUseCase {
	return &fetcherUseCase{
		httpClient: &http.Client{},
	}
}

func (fu *fetcherUseCase) Fetch(
	ctx context.Context,
	targetAPI *url.URL,
	modifyRules models.ModifyRules,
) (map[string]any, error) {
	response, err := fu.httpClient.Get(targetAPI.String())
	if err != nil {
		return nil, fmt.Errorf("call target api: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	responseData := make(map[string]any)
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	modifyRules.ApplyTargetAPIData(responseData)

	return responseData, nil
}
