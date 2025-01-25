/**
 * BaseRepository Definition
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */
package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/frog-engine/frog-go/internal/models"
	"github.com/frog-engine/frog-go/internal/utils"
	"github.com/frog-engine/frog-go/pkg/logger"
)

// BaseRepository 提供泛型接口
type BaseRepository[T any] interface {
    FindByID(id int, query string) (*T, error)
    FindAll(query string) ([]T, error)
		FindWithPagination(
			tableName string, fields []string, condition string,
			pagination *models.Pagination) ([]map[string]interface{}, int, error)
		IsRecordExists(tableName string, conditions map[string]interface{}) (bool, error)
    Create(query string, args ...interface{}) (int64, error)
    Update(query string, entity *T) error
    Delete(query string, id int) error
		GroupBy(entity interface{}, tableName, field string) (map[string]int, error)
}

// BaseRepositoryImpl 提供基本的数据库操作实现
type BaseRepositoryImpl[T any] struct {
    db *sql.DB // 数据库连接实例
}

// NewBaseRepositoryImpl 创建 BaseRepositoryImpl 实例
func NewBaseRepositoryImpl[T any](db *sql.DB) *BaseRepositoryImpl[T] {
    return &BaseRepositoryImpl[T]{db: db}
}

// mapRowToEntity 映射单行数据到实体
func (repo *BaseRepositoryImpl[T]) mapRowToEntity(row *sql.Row, entity *T) error {
	_, err := utils.MapScannerToEntity(row, entity)
	return err
}

// mapRowsToEntity 映射多行数据到实体切片
func (repo *BaseRepositoryImpl[T]) mapRowsToEntity(rows *sql.Rows) ([]T, error) {
	var entities []T

	for rows.Next() {
		var entity T
		if _, err := utils.MapScannerToEntity(rows, &entity); err != nil {
			logger.Printf("Error mapping row to entity: %v", err)
			return nil, err
		}
		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}

// FindByID 查找指定ID的实体
func (repo *BaseRepositoryImpl[T]) FindByID(id int, query string) (*T, error) {
	logger.Printf("BaseRepositoryImpl->FindByID: ID=%d, Query=%s\n", id, query)
	row := repo.db.QueryRow(query, id)
	var entity T
	if err := repo.mapRowToEntity(row, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindAll 查找所有实体，不推荐直接使用，此处仅为示例
func (repo *BaseRepositoryImpl[T]) FindAll(query string) ([]T, error) {
	logger.Printf("BaseRepositoryImpl->FindAll: Query=%s\n", query)
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities, err := repo.mapRowsToEntity(rows)
	if err != nil {
		return nil, err
	}

	// 如果没有数据，返回空切片，而不是 nil
	if len(entities) == 0 {
		return []T{}, nil
	}

	return entities, nil
}

// IsRecordExists 根据指定条件检查数据库记录是否存在，支持排除条件。
// 
// 参数：
//   - tableName: 数据表名称（如 "users"）。
//   - conditions: 查询条件，键为字段名，值为字段值（如 {"email": "test@example.com"}）。
//   - excludeConditions: 排除条件，键为字段名，值为字段值（如 {"id": 1}，表示排除 id=1 的记录）。
//
// 返回值：
//   - bool: 如果满足条件的记录存在，返回 true；否则返回 false。
//   - error: 查询过程中出现的错误信息。
func (repo *BaseRepositoryImpl[T]) IsRecordExists(
  tableName string,
  conditions map[string]interface{},
  excludeConditions map[string]interface{},
) (bool, error) {
  var whereClauses []string
  var args []interface{}

  // 添加查询条件
  for field, value := range conditions {
    whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", field))
    args = append(args, value)
  }

  // 添加排除条件 (如 id != ?)
  for field, value := range excludeConditions {
    whereClauses = append(whereClauses, fmt.Sprintf("%s != ?", field))
    args = append(args, value)
  }

  query := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE %s", tableName, strings.Join(whereClauses, " AND "))
  var count int
  err := repo.db.QueryRow(query, args...).Scan(&count)
  if err != nil {
    return false, err
  }

  return count > 0, nil
}

// Create 创建实体
func (repo *BaseRepositoryImpl[T]) Create(query string, args ...interface{}) (int64, error) {
	logger.Printf("BaseRepositoryImpl->Create: Query=%s, Args=%v\n", query, args)
	// 执行插入 SQL 查询
	result, err := repo.db.Exec(query, args...)
	if err != nil {
			logger.Printf("Error executing insert query: %v\n", err)
			return 0, fmt.Errorf("failed to execute insert query: %w", err)
	}

	// 获取插入的 ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
			logger.Printf("Error getting last insert id: %v\n", err)
			return 0, fmt.Errorf("failed to retrieve last insert id: %w", err)
	}

	return lastInsertID, nil
}

// Update 更新实体
func (repo *BaseRepositoryImpl[T]) Update(query string, args []interface{}) error {
	logger.Printf("BaseRepositoryImpl>Update: Query=%s, Args=%v\n", query, args)
  _, err := repo.db.Exec(query, args...)
  return err
}

// Delete 删除实体
func (repo *BaseRepositoryImpl[T]) Delete(query string, id int) (int64, error)  {
	logger.Printf("BaseRepositoryImpl>Delete: Query=%s, ID=%d\n", query, id)
		result, err := repo.db.Exec(query, id)
		if err != nil {
			return 0, err
		}
		affectedRows, _ := result.RowsAffected()
		return affectedRows, nil
}

func (repo *BaseRepositoryImpl[T]) FindWithPagination(
	tableName string, fields []string, condition string,
	pagination *models.Pagination) ([]map[string]interface{}, int, error) {

	// 计算分页的 offset 和 limit
	offset := (pagination.Page - 1) * pagination.Size
	limit := pagination.Size

	// 设置查询字段
	selectFields := "*"
	if len(fields) > 0 {
		selectFields = strings.Join(fields, ", ")
	}

	// 动态构建 SQL 查询语句
	whereClause := ""
	if condition != "" {
		whereClause = "WHERE " + condition
	}

	// 1. 查询总记录数（无事务）
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", tableName, whereClause)
	err := repo.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		logger.Printf("统计总数出错: %v\n", err)
		total = 0 // 出错时总数设为0
	}

	// 2. 查询分页数据（无事务）
	query := fmt.Sprintf(`
		SELECT %s 
		FROM %s %s
		LIMIT ? OFFSET ?`, selectFields, tableName, whereClause)

	logger.Printf("BaseRepositoryImpl->FindWithPagination: Query=%s, Limit=%d, Offset=%d\n", query, limit, offset)

	rows, err := repo.db.Query(query, limit, offset)
	if err != nil {
		logger.Printf("分页查询出错: %v\n", err)
		return []map[string]interface{}{}, total, nil // 返回空数据但不中断程序
	}
	defer rows.Close()

	// 3. 处理查询结果
	var results []map[string]interface{}
	columns, err := rows.Columns()
	if err != nil {
		logger.Printf("获取列名出错: %v\n", err)
		return []map[string]interface{}{}, total, nil
	}

	for rows.Next() {
		row := make(map[string]interface{})
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			logger.Printf("数据扫描出错: %v\n", err)
			continue // 跳过出错的行
		}

		for i, col := range columns {
			if b, ok := values[i].([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = values[i]
			}
		}
		results = append(results, row)
	}

	if err = rows.Err(); err != nil {
		logger.Printf("遍历数据行出错: %v\n", err)
	}

	// 返回结果和总记录数
	return results, total, nil
}

// GroupBy 根据字段进行分组统计
func (repo *BaseRepositoryImpl[T]) GroupBy(entity interface{}, tableName, field string) (map[string]int, error) {
	// 获取有效字段
	fieldsMap, err := utils.GetStructFieldsByTag(entity, "db")
	if err != nil {
		logger.Printf("Error getting valid fields: %v\n", err)
		return nil, err
	}
	logger.Printf("BaseRepository->GroupBy:fieldsMap: %v\n", fieldsMap)

	// 校验字段是否合法，判断map中的key是否存在
	validFields := strings.Split(field, ",") // 支持多个字段，逗号分隔
	for _, validField := range validFields {
		validField = strings.TrimSpace(validField)
		if _, exists := fieldsMap[validField]; !exists {
			return nil, fmt.Errorf("invalid field for grouping: %s", validField)
		}
	}

	// 构建 SQL 查询
	selectFields := strings.Join(validFields, ", ")
	groupByClause := strings.Join(validFields, ", ")
	query := fmt.Sprintf("SELECT %s, COUNT(*) FROM %s GROUP BY %s", selectFields, tableName, groupByClause)
	logger.Printf("BaseRepository->GroupBy:query: %s\n", query)

	// 执行 SQL 查询
	rows, err := repo.db.Query(query)
	if err != nil {
		logger.Printf("Error executing query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	// 解析结果
	groupedData := make(map[string]int)
	for rows.Next() {
		// 创建与查询字段数目一致的 slice 用于扫描结果
		groupValues := make([]interface{}, len(validFields))
		for i := range groupValues {
			groupValues[i] = new(string) // 假设字段值为 string 类型
		}

		// 扫描查询结果
		var count int
		if err := rows.Scan(append(groupValues, &count)...); err != nil {
			logger.Printf("Error scanning row: %v\n", err)
			return nil, err
		}

		// 生成分组键
		var keyParts []string
		for _, value := range groupValues {
			keyParts = append(keyParts, *(value.(*string))) // 假设字段值为 string
		}
		key := strings.Join(keyParts, ",")

		// 存储结果
		groupedData[key] = count
	}

	// 检查查询是否出错
	if err := rows.Err(); err != nil {
		logger.Printf("Error in rows: %v\n", err)
		return nil, err
	}

	return groupedData, nil
}