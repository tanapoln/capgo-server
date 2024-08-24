package spa

import (
	"io/fs"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

func Middleware(router *gin.Engine, prefix string, rootDir string, fallbackFile string) {
	fsw := &fsWrapper{fs: gin.Dir(rootDir, false)}

	hnd := func(c *gin.Context) {
		name := c.Param("filepath")
		if !fsw.Exist(name) {
			http.ServeFileFS(c.Writer, c.Request, fsw, "index.html")
			return
		}

		http.ServeFileFS(c.Writer, c.Request, fsw, name)
	}

	joinedPrefix := path.Join(prefix, "/*filepath")
	router.GET(joinedPrefix, hnd)
	router.HEAD(joinedPrefix, hnd)
}

var _ fs.FS = &fsWrapper{}

type fsWrapper struct {
	fs http.FileSystem
}

func (f *fsWrapper) Open(name string) (fs.File, error) {
	return f.fs.Open(name)
}

func (f *fsWrapper) Exist(name string) bool {
	h, err := f.fs.Open(name)
	if err != nil {
		return false
	}
	h.Close()
	return true
}
