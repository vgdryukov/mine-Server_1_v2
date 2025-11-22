package server_1

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var cafeList = map[string][]string{
	"moscow": {"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент", "Ложка и вилка"},
	"tula":   {"Пир и мир", "Красиво есть не запретишь", "Поздний завтрак"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	var err error

	count := 25
	countStr := req.FormValue("count")
	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			http.Error(w, "incorrect count", http.StatusBadRequest)
			return
		}
		if count < 0 {
			http.Error(w, "count cannot be negative", http.StatusBadRequest)
			return
		}
	}

	city := req.FormValue("city")
	cafe, ok := cafeList[city]
	if !ok {
		http.Error(w, "unknown city", http.StatusBadRequest)
		return
	}

	if search := req.FormValue("search"); search != "" {
		var found []string
		for _, v := range cafe {
			if strings.Contains(strings.ToLower(v), strings.ToLower(search)) {
				found = append(found, v)
			}
		}
		cafe = found
	}

	count = min(count, len(cafe))
	answer := strings.Join(cafe[:count], ",")
	io.WriteString(w, answer)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	http.HandleFunc("/cafe", mainHandle)

	port := getPort()

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "8080"
}
