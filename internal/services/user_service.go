/**
 * User Service
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package services

import (
  "fmt"

  "github.com/frog-engine/frog-go/internal/models"
  "github.com/frog-engine/frog-go/internal/repositories"
  "github.com/frog-engine/frog-go/pkg/code"
  "github.com/frog-engine/frog-go/pkg/logger"
)

type UserService struct {
  userRepo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
  return &UserService{
    userRepo: repo,
  }
}

func (s *UserService) CreateUser(user *models.User) (int64, error) {
  // 检查邮箱名是否已存在
  conditions := map[string]interface{}{
    "email": user.Email,
  }

  exists, err := s.userRepo.ExistsByConditions(conditions, nil)
  if err != nil {
    logger.Println("Error checking if email exists:", err)
    return -1, err // 查询出错，返回错误
  }

  if exists {
    logger.Println("Email already exists:", user.Email)
    return -1, code.ErrEmailExists // 用户名已存在
  }

  // 创建用户
  userId, err := s.userRepo.Create(user)
  if err != nil {
    logger.Println("Error creating user:", err)
    return -1, code.ErrUserCreateFail // 创建失败
  }

  return userId, nil // 创建成功
}

func (s *UserService) UpdateUser(user *models.User) error {
  if user.Id <= 0 {
    logger.Println("Invalid user ID:", user.Id)
    return code.ErrUserInvalid
  }

  // 检查邮箱是否已存在（排除当前用户）
  conditions := map[string]interface{}{
    "email": user.Email,
  }
  excludeConditions := map[string]interface{}{
    "id": user.Id, // 排除当前用户
  }

  exists, err := s.userRepo.ExistsByConditions(conditions, excludeConditions)
  if err != nil {
    logger.Println("Error checking if email exists during update:", err)
    return code.ErrDatabase // 数据库查询错误
  }

  if exists {
    logger.Println("Email already exists during update:", user.Email)
    return code.ErrUserNotFound // 用户不存在，更新失败
  }

  // 执行更新操作
  err = s.userRepo.Update(user)
  if err != nil {
    logger.Println("Error updating user:", err)
    return code.ErrUserUpdateFail // 更新失败
  }

  return nil // 更新成功
}

func (s *UserService) GetUserById(id int) (*models.User, error) {
  user, err := s.userRepo.FindByID(id)

  if err != nil {
    logger.Println("Error in GetUserById:", err)
    return nil, err
  }

  logger.Println("Successfully retrieved user:", user)

  return user, nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
  users, err := s.userRepo.FindAll()
  if err != nil {
    logger.Printf("Error finding users: %v", err)
    return nil, err
  }
  return users, nil
}

// FindPaged 查询分页用户
func (s *UserService) FindPaged(page, size int, condition string) ([]models.User, *models.Pagination, error) {
  // 可指定查询的字段
  fields := []string{}
  users, pagination, err := s.userRepo.FindPaged(page, size, fields, condition)
  if err != nil {
    logger.Printf("Error fetching users: %v", err)
    return nil, pagination, fmt.Errorf("error fetching users: %v", err)
  }
  return users, pagination, nil
}

func (s *UserService) DeleteUser(id int) (int64, error) {
  result, err := s.userRepo.Delete(id)
  if err != nil {
    logger.Printf("Error deleting user with ID %d: %v", id, err)
  }
  return result, err
}

func (s *UserService) GroupBy(field string) (map[string]int, error) {
  groupedData, err := s.userRepo.GroupBy(field)
  if err != nil {
    logger.Println("Error during grouping:", err)
    return nil, err
  }
  return groupedData, nil
}
