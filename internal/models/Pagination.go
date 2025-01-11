/**
 * Pagination Model Definition
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package models

// Pagination 用于分页查询的结构体
type Pagination struct {
  Page  int `json:"page"`
  Size  int `json:"size"`
  Total int `json:"total"`
}

// 默认分页参数
const DefaultPage = 1
const DefaultSize = 10

// NewPagination 创建一个新的分页对象
func NewPagination(page, size, total int) *Pagination {
  if page <= 0 {
    page = DefaultPage
  }
  if size <= 0 {
    size = DefaultSize
  }
  if total < 0 {
    total = 0
  }
  return &Pagination{Page: page, Size: size, Total: total}
}

// Offset 计算偏移量
func (p *Pagination) Offset() int {
  if p.Page <= 0 {
    p.Page = DefaultPage
  }
  if p.Size <= 0 {
    p.Size = DefaultSize
  }
  // 计算偏移量，假设页码从 1 开始，偏移量 = (页码 - 1) * 每页大小
  return (p.Page - 1) * p.Size
}

// TotalPages 计算总页数
func (p *Pagination) TotalPages() int {
  if p.Total == 0 {
    return 0
  }
  return (p.Total + p.Size - 1) / p.Size
}

// HasNextPage 判断是否有下一页
func (p *Pagination) HasNextPage() bool {
  return p.Page < p.TotalPages()
}

// HasPreviousPage 判断是否有上一页
func (p *Pagination) HasPreviousPage() bool {
  return p.Page > 1
}

// FirstPage 返回第一页的分页对象
func (p *Pagination) FirstPage() *Pagination {
  return &Pagination{Page: 1, Size: p.Size, Total: p.Total}
}

// LastPage 返回最后一页的分页对象
func (p *Pagination) LastPage() *Pagination {
  return &Pagination{Page: p.TotalPages(), Size: p.Size, Total: p.Total}
}

// NextPage 返回下一页的分页对象
func (p *Pagination) NextPage() *Pagination {
  if p.HasNextPage() {
    return &Pagination{Page: p.Page + 1, Size: p.Size, Total: p.Total}
  }
  return p
}

// PreviousPage 返回上一页的分页对象
func (p *Pagination) PreviousPage() *Pagination {
  if p.HasPreviousPage() {
    return &Pagination{Page: p.Page - 1, Size: p.Size, Total: p.Total}
  }
  return p
}
