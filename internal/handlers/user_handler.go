/**
 * UserHandler
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package handlers

import (
  "fmt"

  "strconv"

  "github.com/frog-engine/frog-go/internal/models"
  "github.com/frog-engine/frog-go/internal/services"
  "github.com/frog-engine/frog-go/pkg/code"
  "github.com/frog-engine/frog-go/pkg/logger"
  "github.com/frog-engine/frog-go/pkg/response"
  "github.com/gofiber/fiber/v3"
)

type UserHandler struct {
  userService *services.UserService
}

func NewUserHandler(us *services.UserService) *UserHandler {
  return &UserHandler{
    userService: us,
  }
}

// CreateUser 处理 POST 请求，创建新用户
func (h *UserHandler) CreateUser(c fiber.Ctx) error {
  var user models.User
  if err := c.Bind().Body(&user); err != nil {
    return response.Error(c, fiber.StatusBadRequest, "输入无效: "+err.Error())
  }

  userId, err := h.userService.CreateUser(&user)
  if err != nil {
    if err == code.ErrEmailExists {
      logger.Printf("创建用户 %s 失败:%v", user.Email, err)
      return response.Error(c, fiber.StatusBadRequest, user.Name+"邮箱已存在 "+err.Error())
    }
    return response.Error(c, fiber.StatusInternalServerError, err.Error())
  }

  user.Id = int(userId)
  return response.Success(c, user)
}

// UpdateUser 处理 PUT 请求，更新用户信息
func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
  var user models.User
  if err := c.Bind().Body(&user); err != nil {
    return response.Error(c, fiber.StatusBadRequest, "请求负载无效")
  }

  err := h.userService.UpdateUser(&user)
  if err != nil {
    switch err {
    case code.ErrEmailExists:
      return response.Error(c, fiber.StatusBadRequest, "邮箱已存在，请换一个。")
    case code.ErrUserUpdateFail:
      return response.Error(c, fiber.StatusInternalServerError, "更新用户失败")
    default:
      return response.Error(c, fiber.StatusInternalServerError, "未知错误")
    }
  }

  return response.Success(c, "用户更新成功。")
}

// GetUserById 处理 GET 请求，通过 ID 获取用户
func (h *UserHandler) GetUserById(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil || id <= 0 {
    return response.Error(c, fiber.StatusBadRequest, "无效的用户 ID")
  }

  user, err := h.userService.GetUserById(id)
  if err != nil {
    return response.Error(c, fiber.StatusNotFound, "未找到用户 ID")
  }

  return response.Success(c, user)
}

// GetAllUsers 处理 GET 请求，获取所有用户
func (h *UserHandler) GetAllUsers(c fiber.Ctx) error {
  users, err := h.userService.GetAllUsers()
  if err != nil {
    return response.Error(c, fiber.StatusInternalServerError, "获取用户失败")
  }
  return response.Success(c, users)
}

// FindPagedUsers 处理 GET 请求，分页查询用户
func (h *UserHandler) FindPagedUsers(c fiber.Ctx) error {
  page, err := strconv.Atoi(c.Query("page"))
  if err != nil || page <= 0 {
    page = models.DefaultPage
    return response.Error(c, fiber.StatusBadRequest, "页码无效")
  }

  size, err := strconv.Atoi(c.Query("size"))
  if err != nil || size <= 0 {
    size = models.DefaultSize
    return response.Error(c, fiber.StatusBadRequest, "每页数量无效")
  }

  condition := c.Query("condition", "")

  users, pagination, err := h.userService.FindPaged(page, size, condition)
  if err != nil {
    return response.Error(c, fiber.StatusInternalServerError, "获取用户失败: "+err.Error())
  }

  result := response.PaginationWrapper("users", users, *pagination)
  return response.Success(c, result)
}

// DeleteUser 处理 DELETE 请求，删除指定用户
func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
  idStr := c.Params("id")
  id, err := strconv.Atoi(idStr)
  if err != nil {
    return response.Error(c, fiber.StatusBadRequest, "无效的用户 ID")
  }

  result, err := h.userService.DeleteUser(id)
  if err != nil {
    return response.Error(c, fiber.StatusInternalServerError, "删除用户失败")
  }

  if result > 0 {
    return response.Success(c, fmt.Sprintf("用户 ID %d 已删除。", id))
  }
  return response.Error(c, fiber.StatusOK, fmt.Sprintf("用户 ID %d 不存在。", id))
}

// GroupByHandler 处理按字段分组统计请求
func (h *UserHandler) GroupByHandler(c fiber.Ctx) error {
  field := c.Query("field")
  if field == "" {
    return response.Error(c, fiber.StatusBadRequest, "字段参数是必需的")
  }

  groupedData, err := h.userService.GroupBy(field)
  if err != nil {
    msg := fmt.Sprintf("按字段分组时出错: %v", err)
    return response.Error(c, fiber.StatusInternalServerError, msg)
  }

  return response.Success(c, groupedData)
}
