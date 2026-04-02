package handler

import (
	"net/http"
	"strconv"

	"bokeui/internal/model"

	"github.com/gin-gonic/gin"
)

type FriendLinkHandler struct {
	repo *model.ArticleRepository
}

func NewFriendLinkHandler(repo *model.ArticleRepository) *FriendLinkHandler {
	return &FriendLinkHandler{repo: repo}
}

// List 获取所有友情链接（管理端）
func (h *FriendLinkHandler) List(c *gin.Context) {
	items, err := h.repo.ListFriendLinks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// PublicList 获取可见友情链接（公开接口）
func (h *FriendLinkHandler) PublicList(c *gin.Context) {
	items, err := h.repo.ListVisibleFriendLinks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// Create 创建友情链接
func (h *FriendLinkHandler) Create(c *gin.Context) {
	var link model.FriendLink
	if err := c.ShouldBindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	id, err := h.repo.CreateFriendLink(&link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// Update 更新友情链接
func (h *FriendLinkHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var link model.FriendLink
	if err := c.ShouldBindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	link.ID = id
	if err := h.repo.UpdateFriendLink(&link); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete 删除友情链接
func (h *FriendLinkHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.repo.DeleteFriendLink(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
