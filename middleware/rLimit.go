package mid

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi.io/internal/token"
	"github.com/keller0/yxi.io/internal"
	"strconv"
)

var (
	rq               counter
	maxRequestNumber = internal.GetEnv("MAX_PUBLIC_RUN", "100")
)

type counter struct {
	l   sync.Mutex
	cap int
	val int
}

func (c *counter) full() bool {
	c.l.Lock()
	defer c.l.Unlock()
	if c.val >= c.cap {
		return true
	}
	return false
}
func (c *counter) add() int {
	c.l.Lock()
	defer c.l.Unlock()
	c.val++
	return c.val
}
func (c *counter) done() {
	c.l.Lock()
	c.val--
	c.l.Unlock()
}

func init() {
	max ,_ := strconv.Atoi(maxRequestNumber)
	rq = counter{cap: max, val: 0}
}

// PublicLimit limit unlogined users
func PublicLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := token.JwtGetUserID(c.Request)
		if err != nil {
			if rq.full() {
				c.AbortWithStatusJSON(http.StatusTooManyRequests,
					gin.H{"errNumber": errTooManyRunRequest})
			} else {
				v := rq.add()
				c.Set("cc", v)
			}
		} else {
			c.Set("uid", id)
		}
		c.Next()
		rq.done()
	}
}
