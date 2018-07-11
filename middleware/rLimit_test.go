package mid

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLimit(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(PublicLimit())

	r.POST("/c", func(c *gin.Context) {
		time.Sleep(3 * time.Second)
		c.String(200, "fuck")
	})
	r.POST("/u", func(c *gin.Context) {
		c.String(200, "fuck")
	})

	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/c", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
			assert.Equal(t, "fuck", w.Body.String())
			wg.Done()
		}()
	}
	time.Sleep(1 * time.Second)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/u", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 429, w.Code)

	wg.Wait()
}
