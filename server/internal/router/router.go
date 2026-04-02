package router

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"bokeui/internal/handler"
	"bokeui/internal/middleware"
	"bokeui/internal/model"
	"bokeui/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:9088", "http://127.0.0.1:9088", "http://localhost:5173", "http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	repo := model.NewArticleRepository(db)
	articleH := handler.NewArticleHandler(repo)
	authH := handler.NewAuthHandler(repo)
	siteH := handler.NewSiteHandler(repo)
	toolsH := handler.NewToolsHandler(repo)
	friendLinkH := handler.NewFriendLinkHandler(repo)

	// 初始化邮件服务
	emailSvc := service.NewEmailService(&service.EmailConfig{})
	siteH.SetEmailService(emailSvc)
	articleH.SetEmailService(emailSvc)
	articleH.SetToolsHandler(toolsH)

	// 初始化 IP 地理位置查询服务
	wd, _ := os.Getwd()
	var projectRoot string
	if strings.HasSuffix(filepath.Base(wd), "server") {
		projectRoot = filepath.Dir(wd)
	} else {
		projectRoot = wd
	}
	xdbPath := filepath.Join(projectRoot, "server", "data", "ip2region.xdb")
	service.InitIPRegion(xdbPath)

	// 注册访问记录中间件（过滤 API 和静态资源请求）
	r.Use(func(c *gin.Context) {
		path := c.Request.URL.Path
		// 只记录页面访问，不记录 API 和静态资源
		if !strings.HasPrefix(path, "/api") &&
			!strings.HasPrefix(path, "/assets") &&
			!strings.HasPrefix(path, "/admin") &&
			path != "/favicon.ico" &&
			!strings.Contains(path, ".js") &&
			!strings.Contains(path, ".css") &&
			!strings.Contains(path, ".png") &&
			!strings.Contains(path, ".jpg") &&
			!strings.Contains(path, ".svg") &&
			!strings.Contains(path, ".ico") {
			toolsH.VisitLog(c)
			return
		}
		c.Next()
	})

	// API 路由
	api := r.Group("/api")
	{
		api.POST("/login", authH.Login)

		pub := api.Group("/public")
		{
			pub.GET("/articles", articleH.PublicList)
			pub.GET("/articles/:id", articleH.PublicGet)
			pub.POST("/articles/:id/view", articleH.View)
			pub.GET("/categories", articleH.GetCategories)
			pub.GET("/tags", articleH.GetTags)
			pub.GET("/configs", siteH.GetConfigs)
			pub.POST("/messages", siteH.CreateMessage)
			pub.GET("/friend-links", friendLinkH.PublicList)
		}

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/profile", authH.Profile)
			auth.PUT("/profile", authH.UpdateProfile)
			auth.POST("/change-password", authH.ChangePassword)
			auth.GET("/articles", articleH.List)
			auth.GET("/articles/trash", articleH.ListTrash)          // 回收站列表
			auth.GET("/articles/:id", articleH.GetByID)
			auth.POST("/articles", articleH.Create)
			auth.PUT("/articles/:id", articleH.Update)
			auth.DELETE("/articles/:id", articleH.Delete)
			auth.POST("/articles/:id/restore", articleH.Restore)     // 恢复文章
			auth.DELETE("/articles/:id/hard", articleH.HardDelete)   // 彻底删除
			auth.GET("/categories", articleH.GetCategories)
			auth.GET("/tags", articleH.GetTags)
			auth.GET("/configs", siteH.GetConfigs)
			auth.PUT("/configs", siteH.UpdateConfigs)
			auth.GET("/messages", siteH.ListMessages)
			auth.POST("/messages/reply", siteH.ReplyMessage)
			auth.DELETE("/messages/:id", siteH.DeleteMessage)
			auth.GET("/stats", siteH.Stats)

			// 访问记录
			auth.GET("/visits", toolsH.ListVisitLogs)
			auth.GET("/visits/stats", toolsH.GetVisitStats)

			// Sitemap 和 IndexNow
			auth.GET("/sitemap", toolsH.GenerateSitemap)
			auth.POST("/indexnow/push", toolsH.PushIndexNow)

			// 友情链接
			auth.GET("/friend-links", friendLinkH.List)
			auth.POST("/friend-links", friendLinkH.Create)
			auth.PUT("/friend-links/:id", friendLinkH.Update)
			auth.DELETE("/friend-links/:id", friendLinkH.Delete)
		}

		// 公开的 Sitemap
		api.GET("/sitemap.xml", toolsH.GenerateSitemap)

		// IndexNow Key 文件
		r.GET("/:key", toolsH.IndexNowKeyFile)
	}

	adminDist := filepath.Join(projectRoot, "admin", "dist")
	webDist := filepath.Join(projectRoot, "web", "dist")
	adminIndex := filepath.Join(adminDist, "index.html")
	webIndex := filepath.Join(webDist, "index.html")

	// 1) 先注册静态文件（优先级高于 NoRoute）
	if info, err := os.Stat(adminDist); err == nil && info.IsDir() {
		r.Static("/admin", adminDist)
	}
	if info, err := os.Stat(webDist); err == nil && info.IsDir() {
		// 用 http.FileSystem 正确设置 MIME type
		r.StaticFS("/assets", http.Dir(filepath.Join(webDist, "assets")))
		r.StaticFile("/index.html", webIndex)
		// 其他 web 根目录下的静态资源
		fs := http.Dir(webDist)
		r.GET("/favicon.ico", func(c *gin.Context) {
			c.FileFromFS("favicon.ico", fs)
		})
	}

	// 2) NoRoute 作为 SPA fallback
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 先尝试 web/dist 中是否存在对应文件（处理非 /assets 开头的静态资源）
		webFs := http.Dir(webDist)
		f, err := webFs.Open(strings.TrimPrefix(path, "/"))
		if err == nil {
			stat, _ := f.Stat()
			if stat != nil && !stat.IsDir() {
				f.Close()
				c.FileFromFS(path, http.Dir(webDist))
				return
			}
			f.Close()
		}

		// SPA fallback
		if strings.HasPrefix(path, "/admin") {
			if _, err := os.Stat(adminIndex); err == nil {
				c.File(adminIndex)
				return
			}
		}
		if _, err := os.Stat(webIndex); err == nil {
			c.File(webIndex)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "页面不存在"})
	})

	return r
}
