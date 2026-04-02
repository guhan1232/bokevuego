package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"bokeui/internal/model"
	"bokeui/internal/service"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	repo     *model.ArticleRepository
	emailSvc *service.EmailService
	toolsH   *ToolsHandler
}

func NewArticleHandler(repo *model.ArticleRepository) *ArticleHandler {
	return &ArticleHandler{repo: repo}
}

func (h *ArticleHandler) SetEmailService(svc *service.EmailService) {
	h.emailSvc = svc
}

func (h *ArticleHandler) SetToolsHandler(t *ToolsHandler) {
	h.toolsH = t
}

func (h *ArticleHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	q := model.ArticleQuery{Page: page, PageSize: pageSize}
	q.Keyword = c.Query("keyword")
	q.Category = c.Query("category")

	if s := c.Query("status"); s != "" {
		status, _ := strconv.Atoi(s)
		q.Status = &status
	}

	result, err := h.repo.List(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, a := range result.Items {
		a.Content = "" // 列表不返回内容
	}

	c.JSON(http.StatusOK, gin.H{
		"items": result.Items,
		"total": result.Total,
		"page":  page,
	})
}

func (h *ArticleHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	article, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	tags, _ := h.repo.GetTags(id)
	article.Content = article.Content // 详情页返回内容

	c.JSON(http.StatusOK, gin.H{
		"article": article,
		"tags":    tags,
	})
}

func (h *ArticleHandler) Create(c *gin.Context) {
	var a model.Article
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	a.Author = c.GetString("username")
	id, err := h.repo.Create(&a)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 处理标签
	var req struct {
		Tags []string `json:"tags"`
	}
	c.ShouldBindJSON(&req)
	if len(req.Tags) > 0 {
		h.repo.SetTags(id, req.Tags)
	}

	// 异步：发送邮件通知 + 自动推送 IndexNow
	go func() {
		if a.Status == 1 {
			h.notifyNewArticle(a, c.Request.Host)
			h.autoPushIndexNow(c.Request.Host)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *ArticleHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var a model.Article
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	a.ID = id
	if err := h.repo.Update(&a); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		Tags []string `json:"tags"`
	}
	c.ShouldBindJSON(&req)
	if len(req.Tags) > 0 {
		h.repo.SetTags(id, req.Tags)
	}

	// 异步：发布/更新时自动推送 IndexNow
	go func() {
		if a.Status == 1 {
			h.autoPushIndexNow(c.Request.Host)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// notifyNewArticle 发送新文章邮件通知
func (h *ArticleHandler) notifyNewArticle(a model.Article, host string) {
	if h.emailSvc == nil {
		return
	}
	smtpTo, _ := h.repo.GetConfig("smtp_to")
	if smtpTo == "" {
		return
	}
	siteTitle, _ := h.repo.GetConfig("site_title")
	if siteTitle == "" {
		siteTitle = "BokeUI 博客"
	}
	url := fmt.Sprintf("http://%s/article/%d", host, a.ID)
	err := h.emailSvc.SendWithTemplate(
		smtpTo,
		"["+siteTitle+"] 新文章发布 - "+a.Title,
		"article",
		map[string]interface{}{
			"Title":   a.Title,
			"Content": a.Summary,
			"URL":     url,
		},
	)
	if err != nil {
		log.Println("文章邮件通知发送失败:", err)
	}
}

// autoPushIndexNow 自动推送 IndexNow
func (h *ArticleHandler) autoPushIndexNow(host string) {
	enabled, _ := h.repo.GetConfig("indexnow_enabled")
	if enabled != "true" {
		return
	}
	if h.toolsH == nil {
		return
	}
	// 构造一个模拟的 gin.Context 不太方便，直接在这里实现推送逻辑
	key, _ := h.repo.GetConfig("indexnow_key")
	if key == "" {
		return
	}
	baseURL := "http://" + host

	articles, err := h.repo.GetPublishedArticles()
	if err != nil {
		log.Println("IndexNow 自动推送失败:", err)
		return
	}

	urls := []string{baseURL + "/", baseURL + "/articles", baseURL + "/contact"}
	for _, a := range articles {
		urls = append(urls, fmt.Sprintf("%s/article/%d", baseURL, a.ID))
	}

	// 构建简单 JSON payload
	payload := fmt.Sprintf(`{"host":"%s","key":"%s","keyLocation":"%s/%s.txt","urlList":[`, host, key, baseURL, key)
	for i, u := range urls {
		if i > 0 {
			payload += ","
		}
		payload += `"` + u + `"`
	}
	payload += "]}"

	resp, err := http.Post("https://api.indexnow.org/indexnow", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Println("IndexNow 自动推送失败:", err)
		return
	}
	defer resp.Body.Close()
	log.Println("IndexNow 自动推送完成，状态:", resp.StatusCode)
}

func (h *ArticleHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ListTrash 列出回收站文章
func (h *ArticleHandler) ListTrash(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	q := model.ArticleQuery{Page: page, PageSize: pageSize}
	q.Keyword = c.Query("keyword")

	result, err := h.repo.ListTrash(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, a := range result.Items {
		a.Content = "" // 列表不返回内容
	}

	c.JSON(http.StatusOK, gin.H{
		"items": result.Items,
		"total": result.Total,
		"page":  page,
	})
}

// Restore 恢复文章
func (h *ArticleHandler) Restore(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.repo.Restore(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "恢复成功"})
}

// HardDelete 彻底删除文章
func (h *ArticleHandler) HardDelete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.repo.HardDelete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "彻底删除成功"})
}

func (h *ArticleHandler) GetCategories(c *gin.Context) {
	cats, err := h.repo.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": cats})
}

func (h *ArticleHandler) GetTags(c *gin.Context) {
	tags, err := h.repo.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": tags})
}

func (h *ArticleHandler) View(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	h.repo.IncrementViews(id)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// ===== 前台公开接口 =====

func (h *ArticleHandler) PublicList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "9"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 9
	}

	q := model.ArticleQuery{Page: page, PageSize: pageSize, Status: intPtr(1)}
	q.Keyword = c.Query("keyword")
	q.Category = c.Query("category")

	result, err := h.repo.List(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, a := range result.Items {
		a.Content = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"items": result.Items,
		"total": result.Total,
		"page":  page,
	})
}

func (h *ArticleHandler) PublicGet(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	article, err := h.repo.GetByID(id)
	if err != nil || article.Status != 1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	tags, _ := h.repo.GetTags(id)
	c.JSON(http.StatusOK, gin.H{
		"article": article,
		"tags":    tags,
	})
}

func intPtr(v int) *int { return &v }
