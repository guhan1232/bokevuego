package service

import (
	"log"
	"strings"
	"sync"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var (
	regionOnce sync.Once
	regionInst *IPRegion
)

type IPRegion struct {
	searcher *xdb.Searcher
}

// InitIPRegion 初始化 IP 地理位置查询（全局单例，带 VectorIndex 缓存）
func InitIPRegion(dbPath string) *IPRegion {
	regionOnce.Do(func() {
		var err error

		// 加载 VectorIndex 缓存加速查询
		vIndex, err := xdb.LoadVectorIndexFromFile(dbPath)
		if err != nil {
			log.Printf("[ip2region] 加载 VectorIndex 失败，回退到文件查询: %v", err)
			regionInst = &IPRegion{}
			regionInst.searcher, err = xdb.NewWithFileOnly(xdb.IPv4, dbPath)
			if err != nil {
				log.Printf("[ip2region] 初始化失败: %v", err)
				return
			}
			return
		}

		regionInst = &IPRegion{}
		regionInst.searcher, err = xdb.NewWithVectorIndex(xdb.IPv4, dbPath, vIndex)
		if err != nil {
			log.Printf("[ip2region] 初始化 VectorIndex searcher 失败: %v", err)
			return
		}

		log.Println("[ip2region] IP 地理位置查询服务初始化完成")
	})
	return regionInst
}

// GetRegion 查询 IP 的地理位置，返回格式: "国家|区域|省份|城市|ISP"
func (r *IPRegion) GetRegion(ip string) string {
	if r == nil || r.searcher == nil {
		return ""
	}
	if ip == "" || ip == "127.0.0.1" || ip == "::1" || strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "172.") {
		return "内网IP"
	}

	region, err := r.searcher.SearchByStr(ip)
	if err != nil {
		return ""
	}

	// ip2region 返回格式: "国家|区域|省份|城市|ISP"
	// 清理 "0" 和空值，简化显示
	parts := strings.Split(region, "|")
	cleaned := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" && p != "0" {
			cleaned = append(cleaned, p)
		}
	}
	if len(cleaned) == 0 {
		return ""
	}
	return strings.Join(cleaned, " ")
}

// GetRegion 返回格式化地理位置: "省份 城市 ISP"
func GetRegion(ip string) string {
	if regionInst == nil {
		return ""
	}
	return regionInst.GetRegion(ip)
}
