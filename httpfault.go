package httpfault
import (
	"net/http"
	"fmt"
	"log"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

type HttpFault struct {
	statusCode int
	err        error
}

func New(statusCode int, err error) error {
	return &HttpFault{
		statusCode: statusCode,
		err: err,
	}
}

func NewWithReason(statusCode int, format string, args ...interface{}) error {
	return &HttpFault{statusCode: statusCode, err: fmt.Errorf(format, args...)}
}

func(h HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	err := h(rw, r)
	if err != nil {
		HandleError(rw, err)
		return
	}
}

func (self *HttpFault) Error() string {
	return fmt.Sprintf("%d: reason: '%s'", self.statusCode, self.err.Error())
}

// HandleError exported so other packages can wrap around it.
func HandleError(rw http.ResponseWriter, err error) {
	if fault, ok := err.(*HttpFault); ok {
		log.Printf("Error %d: %v", fault.statusCode, err.Error())
		rw.WriteHeader(fault.statusCode)
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
	}
	rw.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", err.Error())))
}

func FaultyFunc(h HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := h(rw, r)
		if err != nil {
			HandleError(rw, err)
		}
	}
}
