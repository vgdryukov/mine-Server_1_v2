package server_1

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
		{"/cafe?city=moscow&count=-1", http.StatusBadRequest, "count cannot be negative"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		url       string
		wantCount int
		status    int // добавьте статус для проверки
	}{
		{"/cafe?city=moscow&count=1", 1, http.StatusOK},
		{"/cafe?city=moscow&count=3", 3, http.StatusOK},
		{"/cafe?city=moscow&count=10", 5, http.StatusOK}, // всего 5 кафе в Москве
		{"/cafe?city=tula&count=2", 2, http.StatusOK},
		{"/cafe?city=tula&count=5", 3, http.StatusOK},            // всего 3 кафе в Туле
		{"/cafe?city=moscow&count=-1", 0, http.StatusBadRequest}, // если count < 0 ожидаем ошибку
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.url, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)

		if v.status == http.StatusOK {
			answer := response.Body.String()
			if answer == "" {
				assert.Equal(t, 0, v.wantCount, "Empty response but expected %d cafes", v.wantCount)
				continue
			}

			cafes := strings.Split(answer, ",")
			assert.Equal(t, v.wantCount, len(cafes), "URL: %s", v.url)
		}
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string // передаваемое значение search
		wantCount int    // ожидаемое количество кафе в ответе
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		url := "/cafe?city=moscow&search=" + v.search
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)

		answer := response.Body.String()

		// Проверяем количество кафе
		if v.wantCount == 0 {
			assert.Equal(t, "", answer, "Expected empty response for search '%s'", v.search)
		} else {
			cafes := strings.Split(answer, ",")
			assert.Equal(t, v.wantCount, len(cafes), "Search: %s", v.search)

			// Проверяем, что каждое кафе содержит искомую строку (без учета регистра)
			searchLower := strings.ToLower(v.search)
			for _, cafe := range cafes {
				cafeLower := strings.ToLower(cafe)
				assert.True(t, strings.Contains(cafeLower, searchLower),
					"Cafe '%s' should contain '%s'", cafe, v.search)
			}
		}
	}
}

// тест для отрицательного count
func TestCafeNegativeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	response := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/cafe?city=moscow&count=-1", nil)
	handler.ServeHTTP(response, req)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "count cannot be negative", strings.TrimSpace(response.Body.String()))
}
