package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EMus88/Market/internal/models"
	"github.com/EMus88/Market/internal/repository"
	"github.com/EMus88/Market/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	"github.com/pashagolub/pgxmock"
	"github.com/sirupsen/logrus"
)

func Test_AddCategory(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name     string
		category models.Category
		want     want
	}{
		{
			name: "Bad request",
			category: models.Category{
				Name: "",
			},
			want: want{statusCode: 400},
		},
		{
			name: "Internal error",
			category: models.Category{
				Name: "books",
			},
			want: want{statusCode: 500},
		},
		{
			name: "Ok",
			category: models.Category{
				Name: "food",
			},
			want: want{statusCode: 200},
		},
	}
	//init logger
	logger := logrus.New()

	//init mock db connection
	mock, err := pgxmock.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	defer mock.Close(context.Background())

	//init main components
	r := repository.NewRepository(mock, logger)
	s := service.NewService(r, logger)
	h := NewHandler(s, logger)

	//set mock
	Rows := mock.NewRows([]string{"id"}).
		AddRow("123")

	mock.ExpectQuery("INSERT INTO categories").
		WithArgs("food").
		WillReturnRows(Rows)

	//run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category, _ := json.Marshal(tt.category)
			req := httptest.NewRequest(http.MethodPost, "/catalog/category", bytes.NewBuffer(category))
			w := httptest.NewRecorder()

			//init router
			gin.SetMode(gin.ReleaseMode)
			router := gin.Default()
			router.POST("/catalog/category", h.AddCategory)

			//sent request
			router.ServeHTTP(w, req)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, result.StatusCode, tt.want.statusCode)
		})
	}

}
