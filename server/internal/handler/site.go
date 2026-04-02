package handler

import (
	"log"
	"net/http"

	"bokeui/internal/model"
	"bokeui/internal/service"

	"github.com/gin-gonic/gin"
)

type SiteHandler struct {
	repo     *model.ArticleRepository
	emailSvc *service.EmailService
}

func NewSiteHandler(repo *model.ArticleRepository) *SiteHandler {
	return &SiteHandler{repo: repo}
}

func (h *SiteHandler) SetEmailService(svc *service.EmailService) {
	h.emailSvc = svc
}

func (h *SiteHandler) GetConfigs(c *gin.Context) {
	configs, err := h.repo.GetAllConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, configs)
}

func (h *SiteHandler) UpdateConfigs(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	for key, value := range req {
		h.repo.SetConfig(key, value)
	}

	c.JSON(http.StatusOK, gin.H{"message": "保存成功"})
}

// Dashboard
func (h *SiteHandler) Stats(c *gin.Context) {
	stats, err := h.repo.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// Messages
func (h *SiteHandler) CreateMessage(c *gin.Context) {
	var m model.Message
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	id, err := h.repo.CreateMessage(&m)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 异步发送邮件通知
	if h.emailSvc != nil {
		go func() {
			smtpTo, _ := h.repo.GetConfig("smtp_to")
			if smtpTo == "" {
				return
			}
			siteTitle, _ := h.repo.GetConfig("site_title")
			if siteTitle == "" {
				siteTitle = "BokeUI 博客"
			}
			adminURL := "http://" + c.Request.Host + "/admin"
			err := h.emailSvc.SendWithTemplate(
				smtpTo,
				"["+siteTitle+"] 收到新留言 - "+m.Name,
				"message",
				map[string]interface{}{
					"Name":      m.Name,
					"Email":     m.Email,
					"Content":   m.Content,
					"AdminURL":  adminURL,
				},
			)
			if err != nil {
				log.Println("邮件通知发送失败:", err)
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *SiteHandler) ListMessages(c *gin.Context) {
	page, _ := c.GetQuery("page")
	if page == "" {
		page = "1"
	}
	p := 1
	for _, ch := range page {
		if ch >= '0' && ch <= '9' {
			p = p*10 + int(ch-'0')
		}
	}

	items, total, err := h.repo.ListMessages(p, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total, "page": p})
}

func (h *SiteHandler) ReplyMessage(c *gin.Context) {
	var req struct {
		ID    int64  `json:"id"`
		Reply string `json:"reply"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.repo.ReplyMessage(req.ID, req.Reply); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "回复成功"})
}

func (h *SiteHandler) DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	var i int64
	for _, ch := range id {
		i = i*10 + int64(ch-'0')
	}
	if err := h.repo.DeleteMessage(i); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
