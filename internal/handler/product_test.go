package handler

import (
	"bytes"
	"errors"
	"example1/internal/DTO"
	"example1/internal/service"
	mock_service "example1/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProductHandler_Register(t *testing.T) {
	type mockBehaviour func(s mock_service.MockProductService, product *DTO.ReqReserveProduct)

	testTable := &[]struct {
		name                 string
		requestBody          string
		mockBehaviour        mockBehaviour
		reqDTO               DTO.ReqReserveProduct
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			requestBody: "{\n    \"warehouse_id\": 2,\n    \"unique_codes\": [\n        \"olkiuj\",\n      " +
				"  \"tghyuj\"\n    ],\n    \"counts\": [\n        1000,\n        1200\n    ]\n}",
			reqDTO: DTO.ReqReserveProduct{
				WarehouseID: 2,
				UniqueCodes: []string{
					"olkiuj",
					"tghyuj",
				},
				Counts: []int{
					1000,
					1200,
				},
			},
			mockBehaviour: func(s mock_service.MockProductService, product *DTO.ReqReserveProduct) {
				s.EXPECT().Reserve(product).Return(
					&DTO.ResReserveProduct{
						Successful: []DTO.Successful{
							{
								ID:         6,
								UniqueCode: "olkiuj",
							},
							{
								ID:         7,
								UniqueCode: "tghyuj",
							},
						},
						Unsuccessful: []string{},
						Errors:       []string{},
					},
					nil,
				)
			},
			expectedResponseBody: "{\"successful\":[{\"id\":6,\"unique_codes\":\"olkiuj\"},{\"id\":7,\"unique_codes\":\"tghyuj\"}]}",
			expectedStatusCode:   200,
		},

		{
			name: "IM TEAPOT All one bad code and one not enough",
			requestBody: "{\n    \"warehouse_id\": 2,\n    \"unique_codes\": [\n        \"olkij\",\n        \"tghyuj\"\n   " +
				" ],\n    \"counts\": [\n        1000,\n        10000\n    ]\n}",
			reqDTO: DTO.ReqReserveProduct{
				WarehouseID: 2,
				UniqueCodes: []string{
					"olkij",
					"tghyuj",
				},
				Counts: []int{
					1000,
					10000,
				},
			},
			mockBehaviour: func(s mock_service.MockProductService, product *DTO.ReqReserveProduct) {
				s.EXPECT().Reserve(product).Return(
					&DTO.ResReserveProduct{
						Successful: []DTO.Successful{},
						Unsuccessful: []string{
							"olkij",
							"tghyuj",
						},
						Errors: []string{
							"invalid unique code",
							"not enough product",
						},
					},
					nil,
				)
			},
			expectedResponseBody: "{\"errors\":[\"invalid unique code\",\"not enough product\"],\"unsuccessful\":[\"olkij\",\"tghyuj\"]}",
			expectedStatusCode:   400,
		},

		{
			name: "IM TEAPOT unavailable warehouse",
			requestBody: "{\n    \"warehouse_id\": 2,\n    \"unique_codes\": [\n        \"olkij\",\n        \"tghyuj\"\n   " +
				" ],\n    \"counts\": [\n        1000,\n        10000\n    ]\n}",
			reqDTO: DTO.ReqReserveProduct{
				WarehouseID: 2,
				UniqueCodes: []string{
					"olkij",
					"tghyuj",
				},
				Counts: []int{
					1000,
					10000,
				},
			},
			mockBehaviour: func(s mock_service.MockProductService, product *DTO.ReqReserveProduct) {
				s.EXPECT().Reserve(product).Return(
					nil,
					errors.New("warehouse is unavailable"),
				)
			},
			expectedResponseBody: "{\"error\":\"warehouse is unavailable\"}",
			expectedStatusCode:   400,
		},

		{
			name: "MULTI STATUS",
			requestBody: "{\n    \"warehouse_id\": 2,\n    \"unique_codes\": [\n        \"olkij\",\n        " +
				"\"tghyuj\",\n        \"tghyuj\"\n    ],\n    \"counts\": [\n        1000,\n        10000,\n     " +
				"   1000\n    ]\n}",
			reqDTO: DTO.ReqReserveProduct{
				WarehouseID: 2,
				UniqueCodes: []string{
					"olkij",
					"tghyuj",
					"tghyuj",
				},
				Counts: []int{
					1000,
					10000,
					1000,
				},
			},
			mockBehaviour: func(s mock_service.MockProductService, product *DTO.ReqReserveProduct) {
				s.EXPECT().Reserve(product).Return(
					&DTO.ResReserveProduct{
						Successful: []DTO.Successful{
							{
								ID:         8,
								UniqueCode: "tghyuj",
							},
						},
						Unsuccessful: []string{
							"olkij",
							"tghyuj",
						},
						Errors: []string{
							"invalid unique code",
							"not enough product",
						},
					},
					nil,
				)
			},
			expectedResponseBody: "{\"successful\":[{\"id\":8,\"unique_codes\":\"tghyuj\"}],\"unsuccessful\":[\"olkij\",\"tghyuj\"],\"errors\":[\"invalid unique code\",\"not enough product\"]}",
			expectedStatusCode:   207,
		},

		{
			name:        "LACK OF DATA",
			requestBody: "{\n    \"warehouse_id\": 2,\n    \"counts\": [\n        1000,\n        10000,\n        1000\n    ]\n}",
			reqDTO: DTO.ReqReserveProduct{
				WarehouseID: 2,
				Counts: []int{
					1000,
					10000,
					1000,
				},
			},
			mockBehaviour: func(s mock_service.MockProductService, product *DTO.ReqReserveProduct) {
			},
			expectedResponseBody: "{\"error\":\"lack of data\"}",
			expectedStatusCode:   400,
		},

		{
			name:        "INVALID BODY",
			requestBody: "\t{\n   \"warehouse_id\": 2,\n   \"counts\": [\n       1000,\n       \"10000\",\n       1000\n   ]\n}",
			reqDTO:      DTO.ReqReserveProduct{},
			mockBehaviour: func(s mock_service.MockProductService, product *DTO.ReqReserveProduct) {
			},
			expectedResponseBody: "{\"error\":\"invalid body\"}",
			expectedStatusCode:   400,
		},

		{
			name:        "INVALID BODY",
			requestBody: "\t{\n   \"warehouse_id\": 2,\n   \"counts\": [\n       1000,\n       \"10000\",\n       1000\n   ]\n}",
			reqDTO:      DTO.ReqReserveProduct{},
			mockBehaviour: func(s mock_service.MockProductService, product *DTO.ReqReserveProduct) {
			},
			expectedResponseBody: "{\"error\":\"invalid body\"}",
			expectedStatusCode:   400,
		},

		{
			name:        "INVALID WAREHOUSE",
			requestBody: "{\n    \"warehouse_id\": 12,\n    \"unique_codes\": [\n        \"olkiuj\",\n        \"tghyuj\"\n    ],\n    \"counts\": [\n        1000,\n        1200\n    ]\n}",
			reqDTO: DTO.ReqReserveProduct{
				WarehouseID: 12,
				UniqueCodes: []string{
					"olkiuj",
					"tghyuj",
				},
				Counts: []int{
					1000,
					1200,
				},
			},
			mockBehaviour: func(s mock_service.MockProductService, product *DTO.ReqReserveProduct) {
				s.EXPECT().Reserve(product).Return(
					nil,
					errors.New(service.ErrInvalidWarehouse),
				)
			},
			expectedResponseBody: "{\"error\":\"invalid warehouse id\"}",
			expectedStatusCode:   400,
		},
	}

	for _, test := range *testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := gin.New()
			productService := mock_service.NewMockProductService(ctrl)
			test.mockBehaviour(*productService, &test.reqDTO)
			handler := NewProductHandler(productService)
			handler.Register(r)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/ReserveProduct", bytes.NewBufferString(test.requestBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
