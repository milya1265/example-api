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

type MockAuthHandler struct {
	AuthorizeFn func(c *gin.Context)
}

func (m *MockAuthHandler) Authorize(c *gin.Context) {
	m.AuthorizeFn(c)
	return
}

func TestProductHandler_Register(t *testing.T) {
	type mockProductBehaviour func(s mock_service.MockProductService, product *DTO.ReqReserveProduct)
	type mockAuthBehaviour func(c *gin.Context)

	testTable := &[]struct {
		name                 string
		requestBody          string
		mockProductBehaviour mockProductBehaviour
		mockAuthBehavior     mockAuthBehaviour
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
			mockProductBehaviour: func(s mock_service.MockProductService, product *DTO.ReqReserveProduct) {
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
			mockAuthBehavior: func(c *gin.Context) {
				c.Set("id", "lol")
				c.Set("login", "lol")
				c.Set("role", int32(2))
			},
			expectedResponseBody: "{\"successful\":[{\"id\":6,\"unique_codes\":\"olkiuj\"},{\"id\":7,\"unique_codes\":\"tghyuj\"}]}",
			expectedStatusCode:   200,
		},

		{
			name: "Bad request - one bad code and one not enough",
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
			mockProductBehaviour: func(s mock_service.MockProductService, product *DTO.ReqReserveProduct) {
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
			mockAuthBehavior: func(c *gin.Context) {
				c.Set("id", "lol")
				c.Set("login", "lol")
				c.Set("role", int32(2))
			},
			expectedResponseBody: "{\"errors\":[\"invalid unique code\",\"not enough product\"],\"unsuccessful\":[\"olkij\",\"tghyuj\"]}",
			expectedStatusCode:   400,
		},

		{
			name: "Bad request - unavailable warehouse",
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
			mockProductBehaviour: func(s mock_service.MockProductService, product *DTO.ReqReserveProduct) {
				s.EXPECT().Reserve(product).Return(
					nil,
					errors.New("warehouse is unavailable"),
				)
			},
			mockAuthBehavior: func(c *gin.Context) {
				c.Set("id", "lol")
				c.Set("login", "lol")
				c.Set("role", int32(2))
			},
			expectedResponseBody: "{\"error\":\"warehouse is unavailable\"}",
			expectedStatusCode:   400,
		},

		{
			name: "Multi - Status",
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
			mockProductBehaviour: func(s mock_service.MockProductService, product *DTO.ReqReserveProduct) {
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
			mockAuthBehavior: func(c *gin.Context) {
				c.Set("id", "lol")
				c.Set("login", "lol")
				c.Set("role", int32(2))
			},
			expectedResponseBody: "{\"successful\":[{\"id\":8,\"unique_codes\":\"tghyuj\"}],\"unsuccessful\":[\"olkij\",\"tghyuj\"],\"errors\":[\"invalid unique code\",\"not enough product\"]}",
			expectedStatusCode:   207,
		},

		{
			name:        "Lack of data",
			requestBody: "{\n    \"warehouse_id\": 2,\n    \"counts\": [\n        1000,\n        10000,\n        1000\n    ]\n}",
			reqDTO: DTO.ReqReserveProduct{
				WarehouseID: 2,
				Counts: []int{
					1000,
					10000,
					1000,
				},
			},
			mockProductBehaviour: func(s mock_service.MockProductService, product *DTO.ReqReserveProduct) {
			},
			mockAuthBehavior: func(c *gin.Context) {
				c.Set("id", "lol")
				c.Set("login", "lol")
				c.Set("role", int32(2))
			},
			expectedResponseBody: "{\"error\":\"lack of data\"}",
			expectedStatusCode:   400,
		},

		{
			name:        "Invalid Body",
			requestBody: "\t{\n   \"warehouse_id\": 2,\n   \"counts\": [\n       1000,\n       \"10000\",\n       1000\n   ]\n}",
			reqDTO:      DTO.ReqReserveProduct{},
			mockProductBehaviour: func(s mock_service.MockProductService, product *DTO.ReqReserveProduct) {
			},
			mockAuthBehavior: func(c *gin.Context) {
				c.Set("id", "lol")
				c.Set("login", "lol")
				c.Set("role", int32(2))
			},
			expectedResponseBody: "{\"error\":\"invalid body\"}",
			expectedStatusCode:   400,
		},

		{
			name:        "Invalid Body",
			requestBody: "\t{\n   \"warehouse_id\": 2,\n   \"counts\": [\n       1000,\n       \"10000\",\n       1000\n   ]\n}",
			reqDTO:      DTO.ReqReserveProduct{},
			mockProductBehaviour: func(s mock_service.MockProductService, product *DTO.ReqReserveProduct) {
			},
			mockAuthBehavior: func(c *gin.Context) {
				c.Set("id", "lol")
				c.Set("login", "lol")
				c.Set("role", int32(2))
			},
			expectedResponseBody: "{\"error\":\"invalid body\"}",
			expectedStatusCode:   400,
		},

		{
			name:        "Invalid warehouse",
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
			mockProductBehaviour: func(s mock_service.MockProductService, product *DTO.ReqReserveProduct) {
				s.EXPECT().Reserve(product).Return(
					nil,
					service.ErrInvalidWarehouse,
				)
			},
			mockAuthBehavior: func(c *gin.Context) {
				c.Set("id", "lol")
				c.Set("login", "lol")
				c.Set("role", int32(2))
			},
			expectedResponseBody: "{\"error\":\"invalid warehouse id\"}",
			expectedStatusCode:   400,
		},

		{
			name:        "Forbidden",
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
			mockProductBehaviour: func(s mock_service.MockProductService, product *DTO.ReqReserveProduct) {
			},
			mockAuthBehavior: func(c *gin.Context) {
				c.Set("id", "lol")
				c.Set("login", "lol")
				c.Set("role", int32(1))
			},
			expectedResponseBody: "{\"error\":\"forbidden\"}",
			expectedStatusCode:   403,
		},
	}

	for _, test := range *testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := gin.New()

			middleware := &MockAuthHandler{
				AuthorizeFn: test.mockAuthBehavior,
			}

			productService := mock_service.NewMockProductService(ctrl)
			test.mockProductBehaviour(*productService, &test.reqDTO)
			handler := NewProductHandler(productService, middleware)
			handler.Register(r)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/ReserveProduct", bytes.NewBufferString(test.requestBody))
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
