package handler

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahernandez9/rockets/internal/models"
	"github.com/ahernandez9/rockets/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//go:embed testdata/rocket/*.json
var expectedFiles embed.FS

func TestGetRocket(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validUUID := "193270a9-c9cf-404a-8f83-838e71d9ae67"
	invalidUUID := "not-a-uuid"

	expectedRocket := &models.Rocket{
		ID:      validUUID,
		Type:    "Falcon-9",
		Speed:   5000,
		Mission: "ARTEMIS",
		Status:  models.StatusActive,
	}

	tests := []struct {
		name           string
		rocketID       string
		mockSetup      func(*mocks.MockRocketService)
		expectedStatus int
		expectedFile   string
	}{
		{
			name:     "successful retrieval",
			rocketID: validUUID,
			mockSetup: func(m *mocks.MockRocketService) {
				m.EXPECT().
					GetRocket(gomock.Any(), validUUID).
					Return(expectedRocket, nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedFile:   "successful_retrieval.json",
		},
		{
			name:           "invalid UUID format",
			rocketID:       invalidUUID,
			mockSetup:      func(m *mocks.MockRocketService) {},
			expectedStatus: http.StatusBadRequest,
			expectedFile:   "invalid_uuid.json",
		},
		{
			name:     "rocket not found",
			rocketID: validUUID,
			mockSetup: func(m *mocks.MockRocketService) {
				m.EXPECT().
					GetRocket(gomock.Any(), validUUID).
					Return(nil, fmt.Errorf("not found")).
					Times(1)
			},
			expectedStatus: http.StatusNotFound,
			expectedFile:   "not_found.json",
		},
		{
			name:           "malformed UUID",
			rocketID:       "12345-abcde-67890",
			mockSetup:      func(m *mocks.MockRocketService) {},
			expectedStatus: http.StatusBadRequest,
			expectedFile:   "invalid_uuid.json",
		},
		{
			name:     "exploded rocket status",
			rocketID: validUUID,
			mockSetup: func(m *mocks.MockRocketService) {
				explodedRocket := &models.Rocket{
					ID:              validUUID,
					Type:            "Falcon-9",
					Speed:           0,
					Mission:         "FAILED_MISSION",
					Status:          models.StatusExploded,
					ExplosionReason: "PRESSURE_VESSEL_FAILURE",
				}
				m.EXPECT().
					GetRocket(gomock.Any(), validUUID).
					Return(explodedRocket, nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedFile:   "exploded_rocket.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockRocketService(ctrl)

			tt.mockSetup(mockService)

			router := gin.New()
			router.GET("/rockets/:id", GetRocket(mockService))

			req, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodGet,
				fmt.Sprintf("/rockets/%s", tt.rocketID),
				http.NoBody,
			)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code, "unexpected status code")

			if tt.expectedFile != "" {
				expectedJSON, err := expectedFiles.ReadFile("testdata/rocket/" + tt.expectedFile)
				assert.NoError(t, err, fmt.Sprintf("failed to read file: %s", tt.expectedFile))

				actualJSON := w.Body.String()
				assert.JSONEq(t, string(expectedJSON), actualJSON, "response body mismatch")
			}
		})
	}
}
