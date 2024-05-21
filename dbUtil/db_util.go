package dbUtil

import "gorm.io/gorm"

// Pagination 分页
func Pagination(db *gorm.DB, page, size int) *gorm.DB {
	return db.Limit(size).Offset((page - 1) * size)
}
