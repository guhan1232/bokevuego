package handler

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"bokeui/internal/model"
	"bokeui/internal/service"

	"github.com/gin-gonic/gin"
)

type ToolsHandler struct {
	repo *model.ArticleRepository
}

func NewToolsHandler(repo *model.ArticleRepository) *ToolsHandler {
	return &ToolsHandler{repo: repo}
}

// VisitLog 记录访问日志
func (h *ToolsHandler) VisitLog(c *gin.Context) {
	ip := c.ClientIP()
	log := &model.VisitLog{
		Path:      c.Request.URL.Path,
		IP:        ip,
		Region:    service.GetRegion(ip),
		UserAgent: c.Request.UserAgent(),
		Referer:   c.Request.Referer(),
	}
	h.repo.CreateVisitLog(log)
	c.Next()
}

// ListVisitLogs 获取访问记录列表
func (h *ToolsHandler) ListVisitLogs(c *gin.Context) {
	page := 1
	pageSize := 20
	path := c.Query("path")

	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := c.Query("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	items, total, err := h.repo.ListVisitLogs(page, pageSize, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"total": total,
		"page":  page,
	})
}

// GetVisitStats 获取访问统计
func (h *ToolsHandler) GetVisitStats(c *gin.Context) {
	stats, err := h.repo.GetVisitStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// Sitemap URL结构
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

type URL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod"`
	ChangeFreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

// GenerateSitemap 生成站点地图
func (h *ToolsHandler) GenerateSitemap(c *gin.Context) {
	baseURL := c.Query("base_url")
	if baseURL == "" {
		baseURL = "http://localhost:9088"
	}

	articles, err := h.repo.GetPublishedArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	urlset := URLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  []URL{},
	}

	// 首页
	urlset.URLs = append(urlset.URLs, URL{
		Loc:        baseURL + "/",
		LastMod:    time.Now().Format("2006-01-02"),
		ChangeFreq: "daily",
		Priority:   "1.0",
	})

	// 文章列表页
	urlset.URLs = append(urlset.URLs, URL{
		Loc:        baseURL + "/articles",
		LastMod:    time.Now().Format("2006-01-02"),
		ChangeFreq: "daily",
		Priority:   "0.8",
	})

	// 留言页
	urlset.URLs = append(urlset.URLs, URL{
		Loc:        baseURL + "/contact",
		LastMod:    time.Now().Format("2006-01-02"),
		ChangeFreq: "monthly",
		Priority:   "0.5",
	})

	// 文章详情页
	for _, a := range articles {
		lastMod := a.UpdatedAt.Format("2006-01-02")
		if lastMod == "" || strings.HasPrefix(lastMod, "0001") {
			lastMod = a.CreatedAt.Format("2006-01-02")
		}
		urlset.URLs = append(urlset.URLs, URL{
			Loc:        fmt.Sprintf("%s/article/%d", baseURL, a.ID),
			LastMod:    lastMod,
			ChangeFreq: "weekly",
			Priority:   "0.6",
		})
	}

	output, err := xml.MarshalIndent(urlset, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/xml", []byte(xml.Header+string(output)))
}

// IndexNow 推送
func (h *ToolsHandler) PushIndexNow(c *gin.Context) {
	var req struct {
		BaseURL string `json:"base_url"`
		Key     string `json:"key"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.Key == "" {
		// 从配置获取
		req.Key, _ = h.repo.GetConfig("indexnow_key")
	}

	if req.Key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未配置 IndexNow Key"})
		return
	}

	if req.BaseURL == "" {
		req.BaseURL = "http://localhost:9088"
	}

	// 获取所有文章
	articles, err := h.repo.GetPublishedArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建URL列表
	urls := []string{
		req.BaseURL + "/",
		req.BaseURL + "/articles",
		req.BaseURL + "/contact",
	}

	for _, a := range articles {
		urls = append(urls, fmt.Sprintf("%s/article/%d", req.BaseURL, a.ID))
	}

	// 推送到 IndexNow API
	payload := map[string]interface{}{
		"host":        strings.TrimPrefix(req.BaseURL, "http://"),
		"key":         req.Key,
		"keyLocation": fmt.Sprintf("%s/%s.txt", req.BaseURL, req.Key),
		"urlList":     urls,
	}

	jsonData, _ := jsonMarshal(payload)
	resp, err := http.Post("https://api.indexnow.org/indexnow", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "推送失败: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 200 || resp.StatusCode == 202 {
		c.JSON(http.StatusOK, gin.H{
			"message":  "推送成功",
			"count":    len(urls),
			"response": string(body),
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "推送失败",
			"status":   resp.StatusCode,
			"response": string(body),
		})
	}
}

// 生成 IndexNow Key 文件
func (h *ToolsHandler) IndexNowKeyFile(c *gin.Context) {
	key := c.Param("key")
	// 验证 key 是否匹配配置
	configKey, _ := h.repo.GetConfig("indexnow_key")
	if configKey == "" || key != configKey+".txt" {
		c.Status(http.StatusNotFound)
		return
	}
	c.String(http.StatusOK, configKey)
}

// 简单的 JSON marshal（避免导入 encoding/json）
func jsonMarshal(v interface{}) ([]byte, error) {
	switch val := v.(type) {
	case map[string]interface{}:
		var parts []string
		for k, v := range val {
			switch v := v.(type) {
			case string:
				parts = append(parts, fmt.Sprintf(`"%s":"%s"`, k, v))
			case int:
				parts = append(parts, fmt.Sprintf(`"%s":%d`, k, v))
			case []string:
				var arr []string
				for _, s := range v {
					arr = append(arr, fmt.Sprintf(`"%s"`, s))
				}
				parts = append(parts, fmt.Sprintf(`"%s":[%s]`, k, strings.Join(arr, ",")))
			}
		}
		return []byte("{" + strings.Join(parts, ",") + "}"), nil
	}
	return nil, fmt.Errorf("unsupported type")
}
