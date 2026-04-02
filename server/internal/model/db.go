package model

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Article struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Summary   string     `json:"summary"`
	Content   string     `json:"content"`
	Cover     string     `json:"cover"`
	Category  string     `json:"category"`
	Author    string     `json:"author"`
	Status    int        `json:"status"` // 0 draft, 1 published
	Views     int64      `json:"views"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // 软删除时间
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"` // admin, editor
}

type SiteConfig struct {
	ID    int64  `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Message struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Content   string    `json:"content"`
	Reply     string    `json:"reply"`
	CreatedAt time.Time `json:"created_at"`
}

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ArticleTag struct {
	ArticleID int64 `json:"article_id"`
	TagID     int64 `json:"tag_id"`
}

// VisitLog 访问记录
type VisitLog struct {
	ID        int64     `json:"id"`
	Path      string    `json:"path"`
	IP        string    `json:"ip"`
	Region    string    `json:"region"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer"`
	CreatedAt time.Time `json:"created_at"`
}

// FriendLink 友情链接
type FriendLink struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	Logo    string `json:"logo"`
	Sort    int    `json:"sort"`
	Visible bool   `json:"visible"`
}

func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	return db, nil
}

func AutoMigrate(db *sql.DB) {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			nickname TEXT DEFAULT '',
			avatar TEXT DEFAULT '',
			role TEXT DEFAULT 'editor'
		)`,
		`CREATE TABLE IF NOT EXISTS articles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			summary TEXT DEFAULT '',
			content TEXT DEFAULT '',
			cover TEXT DEFAULT '',
			category TEXT DEFAULT '',
			author TEXT DEFAULT '',
			status INTEGER DEFAULT 0,
			views INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS article_tags (
			article_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			PRIMARY KEY (article_id, tag_id),
			FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			content TEXT NOT NULL,
			reply TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS site_configs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key TEXT UNIQUE NOT NULL,
			value TEXT DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS visit_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			path TEXT NOT NULL,
			ip TEXT DEFAULT '',
			user_agent TEXT DEFAULT '',
			referer TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_visit_logs_created_at ON visit_logs(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_visit_logs_path ON visit_logs(path)`,
		`CREATE TABLE IF NOT EXISTS friend_links (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			url TEXT NOT NULL,
			logo TEXT DEFAULT '',
			sort INTEGER DEFAULT 0,
			visible INTEGER DEFAULT 1
		)`,
	}

	for _, s := range statements {
		db.Exec(s)
	}

	// 增量迁移：为 visit_logs 添加 region 字段
	db.Exec("ALTER TABLE visit_logs ADD COLUMN region TEXT DEFAULT ''")

	// 增量迁移：为 articles 添加 deleted_at 字段（回收站功能）
	db.Exec("ALTER TABLE articles ADD COLUMN deleted_at DATETIME DEFAULT NULL")

	// 默认管理员
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 {
		pw, _ := HashPassword("admin123")
		db.Exec("INSERT INTO users (username, password, nickname, role) VALUES (?, ?, ?, ?)",
			"admin", pw, "管理员", "admin")
	}

	// 默认站点配置
	defaults := map[string]string{
		"site_title":       "BokeUI 博客",
		"site_subtitle":    "一个简洁优雅的个人博客",
		"site_footer":      "Powered by BokeUI",
		"site_icp":         "",
		"site_bg_image":    "",
		"site_logo":        "",
		"site_favicon":     "",
		"nav_bg":           "#1a1a2e",
		"hero_bg":          "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
		"smtp_host":        "",
		"smtp_port":        "587",
		"smtp_user":        "",
		"smtp_pass":        "",
		"smtp_to":          "",
		"indexnow_key":     "",
		"indexnow_enabled": "false",
	}
	for k, v := range defaults {
		db.Exec("INSERT OR IGNORE INTO site_configs (key, value) VALUES (?, ?)", k, v)
	}
}
