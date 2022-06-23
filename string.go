package myGin

import (
	"fmt"
	"net/http"
	"unsafe"
)

type String struct {
	Format string
	Data   []interface{}
}

var plainContentType = []string{"text/plain; charset=utf-8"}


func (r String) Render(w http.ResponseWriter) error {
	return WriteString(w, r.Format, r.Data)
}

func (r String) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, plainContentType)
}

func WriteString(w http.ResponseWriter, format string, data []interface{}) (err error) {
	writeContentType(w, plainContentType)
	if len(data) > 0 {
		_, err = fmt.Fprintf(w, format, data...)
		return
	}
	_, err = w.Write(StringToBytes(format))
	return
}

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
