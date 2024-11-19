package not_found

import (
	// "io"
	"fmt"
	"net/http"
	// "github.com/Chayakorn2002/pms-classroom-backend/domain/exceptions"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	// requestBody, _ := io.ReadAll(r.Body)

	// errNotFound := exceptions.NewGlobalErrors().ErrNotFound
	fmt.Println("Not found")

	w.WriteHeader(http.StatusNotFound)
}
