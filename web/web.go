package web

import (
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/MaxBosse/MakaBot/cache"
	"github.com/MaxBosse/MakaBot/utils"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type MakaWeb struct {
	iDB     *utils.InfluxDB
	db      *sql.DB
	cache   *cache.Cache
	router  *gin.Engine
	tickers map[string]*time.Ticker
	store   sessions.CookieStore
}

func NewMakaWeb(bind, secret string, metrics *utils.InfluxDB, db *sql.DB, cache *cache.Cache) *MakaWeb {
	//var err error

	web := new(MakaWeb)
	web.iDB = metrics
	web.tickers = make(map[string]*time.Ticker)
	web.db = db
	web.cache = cache

	web.store = sessions.NewCookieStore([]byte(secret))

	web.router = gin.New()
	web.router.Use(web.globalRecover)
	web.router.Use(sessions.Sessions("session", web.store))
	web.router.Use(gzip.Gzip(gzip.DefaultCompression))

	web.router.HTMLRender = web.loadTemplates("./templates")
	web.router.StaticFS("/assets", http.Dir("./public"))
	web.router.GET("/", web.main)

	pprof.Register(web.router, nil)
	go web.router.Run(bind)

	return web
}

func (web *MakaWeb) loadTemplates(templatesDir string) multitemplate.Render {
	r := multitemplate.New()

	pages, err := filepath.Glob(templatesDir + "/**/*")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/includes/*.tmpl")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our pages/ and includes/ directories
	for _, page := range pages {
		files := append([]string{page}, includes...)
		templateName, _ := filepath.Rel(templatesDir, page)
		templateName = strings.Replace(templateName, "\\", "/", -1)
		n := strings.LastIndexByte(templateName, '.')
		if n > 0 {
			templateName = templateName[:n]
		}
		r.Add(templateName, template.Must(template.ParseFiles(files...)))
	}
	return r
}

func (web *MakaWeb) globalRecover(c *gin.Context) {
	defer func(c *gin.Context) {
		if rec := recover(); rec != nil {
			// that recovery also handle XHR's
			// you need handle it
			if web.XHR(c) {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": rec,
				})
			} else {
				c.HTML(http.StatusOK, "500", gin.H{})
			}
		}
	}(c)
	c.Next()
}

func (web *MakaWeb) XHR(c *gin.Context) bool {
	return strings.ToLower(c.Request.Header.Get("X-Requested-With")) == "xmlhttprequest"
}
