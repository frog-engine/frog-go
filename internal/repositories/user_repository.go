/**
 * UserRepository Definition
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

// UserRepository 定义用户数据操作接口
 type UserRepository interface {
	 FindByID(id int) (*models.User, error)
	 FindAll() ([]models.User, error)
	 ExistsByConditions(conditions map[string]interface{}, excludeConditions map[string]interface{}) (bool, error)
	 Create(user *models.User) (int64, error)
	 Update(user *models.User) error
	 Delete(id int) (int64, error)
	 FindPaged(page int, size int, fields []string, condition string) ([]models.User, *models.Pagination, error)
	 GroupBy(field string) (map[string]int, error)
 }
 
 // SQLUserRepository 提供基于 SQL 的 UserRepository 实现
 type SQLUserRepository struct {
	 baseRepo *BaseRepositoryImpl[models.User]
 }
 
 // NewSQLUserRepository 创建 SQLUserRepository 实例
 func NewSQLUserRepository(db *sql.DB) *SQLUserRepository {
	 return &SQLUserRepository{
		 baseRepo: NewBaseRepositoryImpl[models.User](db),
	 }
 }
 
 func (repo *SQLUserRepository) FindByID(id int) (*models.User, error) {
	 query := `SELECT * FROM users WHERE id = ?`
	 user, err := repo.baseRepo.FindByID(id, query)
	 if err != nil {
		 logger.Printf("FindByID error: %v\n", err)
		 return nil, err
	 }
	 logger.Printf("FindByID result: %+v\n", user)
	 return user, nil
 }
 
 // FindAll 查找所有用户
 func (repo *SQLUserRepository) FindAll() ([]models.User, error) {
	 query := `SELECT * FROM users`
	 users, err := repo.baseRepo.FindAll(query)
	 if err != nil {
		 logger.Printf("FindAll error: %v\n", err)
		 return nil, err
	 }
	 logger.Printf("FindAll result: %+v\n", users)
	 return users, nil
 }
 
 // ExistsByConditions 检查满足指定条件且排除某些条件的用户记录是否存在
 func (repo *SQLUserRepository) ExistsByConditions(
	 conditions map[string]interface{},
	 excludeConditions map[string]interface{},
 ) (bool, error) {
	 return repo.baseRepo.IsRecordExists("users", conditions, excludeConditions)
 }
 
 // Create 创建新用户
 func (repo *SQLUserRepository) Create(user *models.User) (int64, error) {
	 columns, placeholders, values, err := utils.StructToSQL(user, "id", "created_time", "updated_time")
	 if err != nil {
		 return -1, fmt.Errorf("failed to construct SQL: %w", err)
	 }
 
	 query := fmt.Sprintf("INSERT INTO users (%s) VALUES (%s)", strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	 lastInsertID, err := repo.baseRepo.Create(query, values...)
	 if err != nil {
		 return -1, fmt.Errorf("failed to create user: %w", err)
	 }
 
	 logger.Printf("Create result: Last Insert ID: %d\n", lastInsertID)
	 return lastInsertID, nil
 }
 
 // Update 更新用户
 func (repo *SQLUserRepository) Update(user *models.User) error {
	 setClause, args, err := utils.StructToSlice(user, "db", "created_time", "updated_time")
	 if err != nil {
		 return fmt.Errorf("failed to build SQL set clause: %v", err)
	 }
 
	 query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", setClause)
	 args = append(args, user.Id)
 
	 err = repo.baseRepo.Update(query, args)
	 if err != nil {
		 logger.Printf("Update error: %v\n", err)
	 }
	 return err
 }
 
 // Delete 删除用户
 func (repo *SQLUserRepository) Delete(id int) (int64, error) {
	 query := `DELETE FROM users WHERE id = ?`
	 result, err := repo.baseRepo.Delete(query, id)
	 if err != nil {
		logger.Printf("Delete error: %v, User ID: %d\n", err, id)
	 }
	 return result, err
 }
 
// FindPaged 分页查询用户
func (repo *SQLUserRepository) FindPaged(page int, size int, fields []string, condition string) ([]models.User, *models.Pagination, error) {
	pagination := models.NewPagination(page, size, 0)
	tableName := "users"

	results, total, err := repo.baseRepo.FindWithPagination(tableName, fields, condition, pagination)
	if err != nil {
		logger.Printf("FindPaged error: %v, Page: %d, Size: %d, Condition: %s\n", err, page, size, condition)
		return nil, pagination, err
	}
	pagination.Total = total

	var users []models.User
	for _, row := range results {
		var user models.User
		if err := utils.MapRowToStruct(row, &user); err != nil {
			return nil, pagination, err
		}
		users = append(users, user)
	}

	logger.Printf("FindPaged results: Users: %+v, Pagination: %+v\n", users, pagination)
	return users, pagination, nil
}
 
// GroupBy 根据字段进行分组统计
func (repo *SQLUserRepository) GroupBy(field string) (map[string]int, error) {
	groupedData, err := repo.baseRepo.GroupBy(models.User{}, "users", field)
	if err != nil {
		logger.Printf("GroupBy error for field %s: %v\n", field, err)
		return nil, err
	}
	logger.Printf("GroupBy result for field %s: %+v\n", field, groupedData)
	return groupedData, nil
}
 