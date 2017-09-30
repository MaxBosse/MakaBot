package web

import (
	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

func (web *MakaWeb) main(c *gin.Context) {
	c.HTML(200, "pages/index", gin.H{
		"title": "Welcome!",
	})
}
