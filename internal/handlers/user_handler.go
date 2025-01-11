/**
 * UserHandler
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package handlers

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "strconv"

  "frog-go/internal/models"
  "frog-go/internal/services"
  "frog-go/pkg/code"
  "frog-go/pkg/response"

  "github.com/gorilla/mux"
)

type UserHandler struct {
  userService *services.UserService
}

func NewUserHandler(us *services.UserService) *UserHandler {
  return &UserHandler{
    userService: us,
  }
}

// CreateUser 处理 POST 请求，根据JSON创建用户
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
  var user models.User
  decoder := json.NewDecoder(r.Body)
  // 打印日志检查请求数据
  log.Println("UserHandler->CreateUser:user:", user)
  if err := decoder.Decode(&user); err != nil {
    response.Error(w, http.StatusBadRequest, "Invalid input: "+err.Error()) // 使用 Error 方法返回错误响应
    return
  }

  userId, err := h.userService.CreateUser(&user)
  if err != nil {
    // 如果用户已存在等错误，返回业务错误
    if err == code.ErrEmailExists {
      response.Error(w, http.StatusBadRequest, user.Name+err.Error())
      return
    }

    if err == code.ErrUserCreateFail {
      response.Error(w, http.StatusBadRequest, user.Name+err.Error())
      return
    }

    // 其他错误返回 500
    response.Error(w, http.StatusInternalServerError, "Error: "+err.Error())
    return
  }

  user.Id = int(userId)
  w.WriteHeader(http.StatusCreated)
  response.Success(w, user)
}

// UpdateUser 处理 PUT请求，更新用户信息
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
  var user models.User

  err := json.NewDecoder(r.Body).Decode(&user)
  log.Println("user:", user)
  if err != nil {
    response.Error(w, http.StatusBadRequest, "Invalid request payload") // 请求体解析错误
    return
  }

  // 调用服务层更新用户
  err = h.userService.UpdateUser(&user)
  if err != nil {
    switch err {
    case code.ErrEmailExists:
      // 邮箱已存在，返回 400 错误
      response.Error(w, http.StatusBadRequest, fmt.Sprintf("邮箱 '%s' 已存在，请换一个。", user.Email))
    case code.ErrUserUpdateFail:
      // 更新失败，返回 500 错误
      response.Error(w, http.StatusInternalServerError, fmt.Sprintf("更新用户失败: %v", err))
    default:
      // 其他未知错误
      response.Error(w, http.StatusInternalServerError, fmt.Sprintf("未知错误: %v", err))
    }
    return
  }

  // 更新成功
  response.Success(w, fmt.Sprintf("User id %d updated successfully.", user.Id))
}

// GetUserById 处理 GET 请求，根据ID返回用户
func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  idStr := vars["id"]
  id, err := strconv.Atoi(idStr)
  if err != nil || id <= 0 {
    response.Error(w, http.StatusBadRequest, "invalid user ID")
    return
  }

  user, err := h.userService.GetUserById(id)
  if err != nil {
    // defaultUser := &models.User{
    //   Id: id,
    // }
    // response.Success(w, defaultUser)
    response.Error(w, http.StatusBadRequest, "Not found user ID")
    return
  }

  response.Success(w, user)
}

// GetAllUsers 处理 GET 请求，返回所有用户
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
  users, err := h.userService.GetAllUsers()
  if err != nil {
    response.Error(w, http.StatusBadRequest, "Error fetching users")
    return
  }
  response.Success(w, users)
}

// FindPagedUsers 处理 GET 请求，分页查询用户
func (h *UserHandler) FindPagedUsers(w http.ResponseWriter, r *http.Request) {
  // 获取查询参数 "page" 和 "size"
  pageStr := r.URL.Query().Get("page")
  sizeStr := r.URL.Query().Get("size")
  conditionStr := r.URL.Query().Get("condition")

  // 默认分页参数
  if pageStr == "" {
    pageStr = "1"
  }
  if sizeStr == "" {
    sizeStr = "10"
  }

  page, err := strconv.Atoi(pageStr)
  if err != nil || page <= 0 {
    response.Error(w, http.StatusBadRequest, "页码无效")
    return
  }

  size, err := strconv.Atoi(sizeStr)
  if err != nil || size <= 0 {
    response.Error(w, http.StatusBadRequest, "每页数量无效")
    return
  }

  // 调用 service 层获取分页用户数据
  users, pagination, err := h.userService.FindPaged(page, size, conditionStr)
  if err != nil {
    response.Error(w, http.StatusInternalServerError, "Error fetching users: "+err.Error())
    return
  }

  result := response.PaginationWrapper("users", users, *pagination)
  response.Success(w, result)
}

// DeleteUser 处理 DELETE 请求，删除指定用户
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
  idStr := mux.Vars(r)["id"]
  id, err := strconv.Atoi(idStr)
  if err != nil {
    response.Error(w, http.StatusBadRequest, "Invalid user ID")
    return
  }

  result, err := h.userService.DeleteUser(id)
  if err != nil {
    response.Error(w, http.StatusInternalServerError, "Error deleting user")
    return
  }
  if result > 0 {
    response.Success(w, fmt.Sprintf("User id %d was deleted.", id))
  } else {
    response.Success(w, fmt.Sprintf("User id %d was not exist.", id))
  }
}

// GroupByHandler 处理按字段分组统计请求
func (h *UserHandler) GroupByHandler(w http.ResponseWriter, r *http.Request) {
  // 从请求中获取字段
  field := r.URL.Query().Get("field")
  if field == "" {
    response.Error(w, http.StatusBadRequest, "Field parameter is required")
    return
  }

  groupedData, err := h.userService.GroupBy(field)
  if err != nil {
    msg := fmt.Sprintf("Error grouping by field: %v", err)
    response.Error(w, http.StatusInternalServerError, msg)
    return
  }

  // 返回分组结果
  response.Success(w, groupedData)
}
