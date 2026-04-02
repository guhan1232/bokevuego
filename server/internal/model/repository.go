package model

import "database/sql"

type ArticleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) Create(a *Article) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO articles (title, summary, content, cover, category, author, status) VALUES (?,?,?,?,?,?,?)",
		a.Title, a.Summary, a.Content, a.Cover, a.Category, a.Author, a.Status,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ArticleRepository) Update(a *Article) error {
	_, err := r.db.Exec(
		"UPDATE articles SET title=?, summary=?, content=?, cover=?, category=?, author=?, status=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		a.Title, a.Summary, a.Content, a.Cover, a.Category, a.Author, a.Status, a.ID,
	)
	return err
}

// Delete 软删除文章（移入回收站）
func (r *ArticleRepository) Delete(id int64) error {
	_, err := r.db.Exec("UPDATE articles SET deleted_at = CURRENT_TIMESTAMP WHERE id=? AND deleted_at IS NULL", id)
	return err
}

// SoftDelete 软删除文章（移入回收站）
func (r *ArticleRepository) SoftDelete(id int64) error {
	_, err := r.db.Exec("UPDATE articles SET deleted_at = CURRENT_TIMESTAMP WHERE id=? AND deleted_at IS NULL", id)
	return err
}

// Restore 恢复文章（从回收站恢复）
func (r *ArticleRepository) Restore(id int64) error {
	_, err := r.db.Exec("UPDATE articles SET deleted_at = NULL WHERE id=?", id)
	return err
}

// HardDelete 彻底删除文章（物理删除）
func (r *ArticleRepository) HardDelete(id int64) error {
	_, err := r.db.Exec("DELETE FROM articles WHERE id=?", id)
	return err
}

func (r *ArticleRepository) GetByID(id int64) (*Article, error) {
	a := &Article{}
	err := r.db.QueryRow(
		"SELECT id, title, summary, content, cover, category, author, status, views, created_at, updated_at, deleted_at FROM articles WHERE id=? AND deleted_at IS NULL",
		id,
	).Scan(&a.ID, &a.Title, &a.Summary, &a.Content, &a.Cover, &a.Category, &a.Author, &a.Status, &a.Views, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// GetTrashByID 获取回收站中的文章
func (r *ArticleRepository) GetTrashByID(id int64) (*Article, error) {
	a := &Article{}
	err := r.db.QueryRow(
		"SELECT id, title, summary, content, cover, category, author, status, views, created_at, updated_at, deleted_at FROM articles WHERE id=? AND deleted_at IS NOT NULL",
		id,
	).Scan(&a.ID, &a.Title, &a.Summary, &a.Content, &a.Cover, &a.Category, &a.Author, &a.Status, &a.Views, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *ArticleRepository) IncrementViews(id int64) error {
	_, err := r.db.Exec("UPDATE articles SET views = views + 1 WHERE id=? AND deleted_at IS NULL", id)
	return err
}

type ArticleQuery struct {
	Page     int
	PageSize int
	Status   *int
	Category string
	Keyword  string
}

type ArticleListResult struct {
	Items []*Article
	Total int64
}

func (r *ArticleRepository) List(q ArticleQuery) (*ArticleListResult, error) {
	var conditions []string
	var args []interface{}

	// 默认过滤已删除的文章
	conditions = append(conditions, "deleted_at IS NULL")

	if q.Status != nil {
		conditions = append(conditions, "status=?")
		args = append(args, *q.Status)
	}
	if q.Category != "" {
		conditions = append(conditions, "category=?")
		args = append(args, q.Category)
	}
	if q.Keyword != "" {
		conditions = append(conditions, "(title LIKE ? OR summary LIKE ?)")
		args = append(args, "%"+q.Keyword+"%", "%"+q.Keyword+"%")
	}

	where := " WHERE " + conditions[0]
	for _, c := range conditions[1:] {
		where += " AND " + c
	}

	var total int64
	r.db.QueryRow("SELECT COUNT(*) FROM articles"+where, args...).Scan(&total)

	offset := (q.Page - 1) * q.PageSize
	args = append(args, offset, q.PageSize)
	rows, err := r.db.Query(
		"SELECT id, title, summary, content, cover, category, author, status, views, created_at, updated_at FROM articles"+where+" ORDER BY created_at DESC LIMIT ?,?",
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Article
	for rows.Next() {
		a := &Article{}
		rows.Scan(&a.ID, &a.Title, &a.Summary, &a.Content, &a.Cover, &a.Category, &a.Author, &a.Status, &a.Views, &a.CreatedAt, &a.UpdatedAt)
		items = append(items, a)
	}

	return &ArticleListResult{Items: items, Total: total}, nil
}

// ListTrash 列出回收站中的文章
func (r *ArticleRepository) ListTrash(q ArticleQuery) (*ArticleListResult, error) {
	var conditions []string
	var args []interface{}

	// 只查询已删除的文章
	conditions = append(conditions, "deleted_at IS NOT NULL")

	if q.Keyword != "" {
		conditions = append(conditions, "(title LIKE ? OR summary LIKE ?)")
		args = append(args, "%"+q.Keyword+"%", "%"+q.Keyword+"%")
	}

	where := " WHERE " + conditions[0]
	for _, c := range conditions[1:] {
		where += " AND " + c
	}

	var total int64
	r.db.QueryRow("SELECT COUNT(*) FROM articles"+where, args...).Scan(&total)

	offset := (q.Page - 1) * q.PageSize
	args = append(args, offset, q.PageSize)
	rows, err := r.db.Query(
		"SELECT id, title, summary, content, cover, category, author, status, views, created_at, updated_at, deleted_at FROM articles"+where+" ORDER BY deleted_at DESC LIMIT ?,?",
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Article
	for rows.Next() {
		a := &Article{}
		rows.Scan(&a.ID, &a.Title, &a.Summary, &a.Content, &a.Cover, &a.Category, &a.Author, &a.Status, &a.Views, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt)
		items = append(items, a)
	}

	return &ArticleListResult{Items: items, Total: total}, nil
}

func (r *ArticleRepository) GetCategories() ([]string, error) {
	rows, err := r.db.Query("SELECT DISTINCT category FROM articles WHERE category != '' AND deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cats []string
	for rows.Next() {
		var c string
		rows.Scan(&c)
		cats = append(cats, c)
	}
	return cats, nil
}

// Tag methods
func (r *ArticleRepository) SetTags(articleID int64, tagNames []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.Exec("DELETE FROM article_tags WHERE article_id=?", articleID)

	for _, name := range tagNames {
		var tagID int64
		err := tx.QueryRow("INSERT OR IGNORE INTO tags (name) VALUES (?) RETURNING id", name).Scan(&tagID)
		if err != nil {
			err = tx.QueryRow("SELECT id FROM tags WHERE name=?", name).Scan(&tagID)
			if err != nil {
				continue
			}
		}
		tx.Exec("INSERT OR IGNORE INTO article_tags (article_id, tag_id) VALUES (?, ?)", articleID, tagID)
	}

	return tx.Commit()
}

func (r *ArticleRepository) GetTags(articleID int64) ([]string, error) {
	rows, err := r.db.Query(
		"SELECT t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE at.article_id=?",
		articleID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []string
	for rows.Next() {
		var t string
		rows.Scan(&t)
		tags = append(tags, t)
	}
	return tags, nil
}

func (r *ArticleRepository) GetAllTags() ([]string, error) {
	rows, err := r.db.Query("SELECT name FROM tags ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []string
	for rows.Next() {
		var t string
		rows.Scan(&t)
		tags = append(tags, t)
	}
	return tags, nil
}

// Message methods
func (r *ArticleRepository) CreateMessage(m *Message) (int64, error) {
	result, err := r.db.Exec("INSERT INTO messages (name, email, content) VALUES (?,?,?)", m.Name, m.Email, m.Content)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ArticleRepository) ListMessages(page, pageSize int) ([]*Message, int64, error) {
	var total int64
	r.db.QueryRow("SELECT COUNT(*) FROM messages").Scan(&total)

	offset := (page - 1) * pageSize
	rows, err := r.db.Query("SELECT id, name, email, content, reply, created_at FROM messages ORDER BY created_at DESC LIMIT ?,?", offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []*Message
	for rows.Next() {
		m := &Message{}
		rows.Scan(&m.ID, &m.Name, &m.Email, &m.Content, &m.Reply, &m.CreatedAt)
		items = append(items, m)
	}
	return items, total, nil
}

func (r *ArticleRepository) ReplyMessage(id int64, reply string) error {
	_, err := r.db.Exec("UPDATE messages SET reply=? WHERE id=?", reply, id)
	return err
}

func (r *ArticleRepository) DeleteMessage(id int64) error {
	_, err := r.db.Exec("DELETE FROM messages WHERE id=?", id)
	return err
}

// Site config
func (r *ArticleRepository) GetConfig(key string) (string, error) {
	var v string
	err := r.db.QueryRow("SELECT value FROM site_configs WHERE key=?", key).Scan(&v)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return v, err
}

func (r *ArticleRepository) GetAllConfigs() (map[string]string, error) {
	rows, err := r.db.Query("SELECT key, value FROM site_configs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	configs := map[string]string{}
	for rows.Next() {
		var k, v string
		rows.Scan(&k, &v)
		configs[k] = v
	}
	return configs, nil
}

func (r *ArticleRepository) SetConfig(key, value string) error {
	_, err := r.db.Exec("INSERT OR REPLACE INTO site_configs (key, value) VALUES (?, ?)", key, value)
	return err
}

// User
func (r *ArticleRepository) GetUserByUsername(username string) (*User, error) {
	u := &User{}
	err := r.db.QueryRow("SELECT id, username, password, nickname, avatar, role FROM users WHERE username=?", username).
		Scan(&u.ID, &u.Username, &u.Password, &u.Nickname, &u.Avatar, &u.Role)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *ArticleRepository) GetUserByID(id int64) (*User, error) {
	u := &User{}
	err := r.db.QueryRow("SELECT id, username, nickname, avatar, role FROM users WHERE id=?", id).
		Scan(&u.ID, &u.Username, &u.Nickname, &u.Avatar, &u.Role)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *ArticleRepository) UpdatePassword(userID int64, hashedPassword string) error {
	_, err := r.db.Exec("UPDATE users SET password=? WHERE id=?", hashedPassword, userID)
	return err
}

func (r *ArticleRepository) UpdateUserProfile(userID int64, nickname, avatar string) error {
	_, err := r.db.Exec("UPDATE users SET nickname=?, avatar=? WHERE id=?", nickname, avatar, userID)
	return err
}

// Dashboard stats
func (r *ArticleRepository) GetStats() (map[string]interface{}, error) {
	var totalArticles, publishedArticles, totalMessages, draftArticles, trashArticles int
	var totalViews int

	r.db.QueryRow("SELECT COUNT(*) FROM articles WHERE deleted_at IS NULL").Scan(&totalArticles)
	r.db.QueryRow("SELECT COUNT(*) FROM articles WHERE status=1 AND deleted_at IS NULL").Scan(&publishedArticles)
	r.db.QueryRow("SELECT COUNT(*) FROM messages").Scan(&totalMessages)
	r.db.QueryRow("SELECT COUNT(*) FROM articles WHERE status=0 AND deleted_at IS NULL").Scan(&draftArticles)
	r.db.QueryRow("SELECT COUNT(*) FROM articles WHERE deleted_at IS NOT NULL").Scan(&trashArticles)
	r.db.QueryRow("SELECT COALESCE(SUM(views), 0) FROM articles WHERE deleted_at IS NULL").Scan(&totalViews)

	// 访问统计
	var todayVisits, totalVisits int
	r.db.QueryRow("SELECT COUNT(*) FROM visit_logs WHERE date(created_at) = date('now')").Scan(&todayVisits)
	r.db.QueryRow("SELECT COUNT(*) FROM visit_logs").Scan(&totalVisits)

	return map[string]interface{}{
		"total_articles":     totalArticles,
		"published_articles": publishedArticles,
		"total_messages":     totalMessages,
		"draft_articles":     draftArticles,
		"trash_articles":     trashArticles,
		"total_views":        totalViews,
		"today_visits":       todayVisits,
		"total_visits":       totalVisits,
	}, nil
}

// Visit log methods
func (r *ArticleRepository) CreateVisitLog(log *VisitLog) error {
	_, err := r.db.Exec(
		"INSERT INTO visit_logs (path, ip, region, user_agent, referer) VALUES (?, ?, ?, ?, ?)",
		log.Path, log.IP, log.Region, log.UserAgent, log.Referer,
	)
	return err
}

func (r *ArticleRepository) ListVisitLogs(page, pageSize int, path string) ([]*VisitLog, int64, error) {
	var conditions []string
	var args []interface{}

	if path != "" {
		conditions = append(conditions, "path LIKE ?")
		args = append(args, "%"+path+"%")
	}

	where := ""
	if len(conditions) > 0 {
		where = " WHERE " + conditions[0]
	}

	var total int64
	r.db.QueryRow("SELECT COUNT(*) FROM visit_logs"+where, args...).Scan(&total)

	offset := (page - 1) * pageSize
	args = append(args, offset, pageSize)

	query := "SELECT id, path, ip, region, user_agent, referer, created_at FROM visit_logs" + where + " ORDER BY created_at DESC LIMIT ?,?"
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []*VisitLog
	for rows.Next() {
		log := &VisitLog{}
		rows.Scan(&log.ID, &log.Path, &log.IP, &log.Region, &log.UserAgent, &log.Referer, &log.CreatedAt)
		items = append(items, log)
	}

	return items, total, nil
}

func (r *ArticleRepository) GetVisitStats() (map[string]interface{}, error) {
	stats := map[string]interface{}{}

	// 今日访问
	var todayVisits int
	r.db.QueryRow("SELECT COUNT(*) FROM visit_logs WHERE date(created_at) = date('now')").Scan(&todayVisits)
	stats["today"] = todayVisits

	// 昨日访问
	var yesterdayVisits int
	r.db.QueryRow("SELECT COUNT(*) FROM visit_logs WHERE date(created_at) = date('now', '-1 day')").Scan(&yesterdayVisits)
	stats["yesterday"] = yesterdayVisits

	// 本周访问
	var weekVisits int
	r.db.QueryRow("SELECT COUNT(*) FROM visit_logs WHERE date(created_at) >= date('now', '-7 days')").Scan(&weekVisits)
	stats["week"] = weekVisits

	// 总访问
	var totalVisits int
	r.db.QueryRow("SELECT COUNT(*) FROM visit_logs").Scan(&totalVisits)
	stats["total"] = totalVisits

	// 热门页面
	rows, err := r.db.Query(`
		SELECT path, COUNT(*) as cnt 
		FROM visit_logs 
		WHERE date(created_at) >= date('now', '-7 days')
		GROUP BY path 
		ORDER BY cnt DESC 
		LIMIT 10
	`)
	if err == nil {
		defer rows.Close()
		var topPages []map[string]interface{}
		for rows.Next() {
			var path string
			var cnt int
			rows.Scan(&path, &cnt)
			topPages = append(topPages, map[string]interface{}{"path": path, "count": cnt})
		}
		stats["top_pages"] = topPages
	}

	return stats, nil
}

// 获取所有已发布文章URL（用于Sitemap）
func (r *ArticleRepository) GetPublishedArticles() ([]*Article, error) {
	rows, err := r.db.Query(
		"SELECT id, title, created_at, updated_at FROM articles WHERE status=1 AND deleted_at IS NULL ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Article
	for rows.Next() {
		a := &Article{}
		rows.Scan(&a.ID, &a.Title, &a.CreatedAt, &a.UpdatedAt)
		items = append(items, a)
	}
	return items, nil
}

// FriendLink methods
func (r *ArticleRepository) CreateFriendLink(link *FriendLink) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO friend_links (name, url, logo, sort, visible) VALUES (?, ?, ?, ?, ?)",
		link.Name, link.URL, link.Logo, link.Sort, link.Visible,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ArticleRepository) UpdateFriendLink(link *FriendLink) error {
	_, err := r.db.Exec(
		"UPDATE friend_links SET name=?, url=?, logo=?, sort=?, visible=? WHERE id=?",
		link.Name, link.URL, link.Logo, link.Sort, link.Visible, link.ID,
	)
	return err
}

func (r *ArticleRepository) DeleteFriendLink(id int64) error {
	_, err := r.db.Exec("DELETE FROM friend_links WHERE id=?", id)
	return err
}

func (r *ArticleRepository) GetFriendLinkByID(id int64) (*FriendLink, error) {
	link := &FriendLink{}
	err := r.db.QueryRow("SELECT id, name, url, logo, sort, visible FROM friend_links WHERE id=?", id).
		Scan(&link.ID, &link.Name, &link.URL, &link.Logo, &link.Sort, &link.Visible)
	if err != nil {
		return nil, err
	}
	return link, nil
}

func (r *ArticleRepository) ListFriendLinks() ([]*FriendLink, error) {
	rows, err := r.db.Query("SELECT id, name, url, logo, sort, visible FROM friend_links ORDER BY sort ASC, id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*FriendLink
	for rows.Next() {
		link := &FriendLink{}
		rows.Scan(&link.ID, &link.Name, &link.URL, &link.Logo, &link.Sort, &link.Visible)
		items = append(items, link)
	}
	return items, nil
}

func (r *ArticleRepository) ListVisibleFriendLinks() ([]*FriendLink, error) {
	rows, err := r.db.Query("SELECT id, name, url, logo, sort, visible FROM friend_links WHERE visible=1 ORDER BY sort ASC, id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*FriendLink
	for rows.Next() {
		link := &FriendLink{}
		rows.Scan(&link.ID, &link.Name, &link.URL, &link.Logo, &link.Sort, &link.Visible)
		items = append(items, link)
	}
	return items, nil
}
