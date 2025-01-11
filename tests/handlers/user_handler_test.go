/**
 * Image Handler Tests
 *
 * @author jarryli@gmail.com
 * @date 2024-03-20
 */

package handlers_test

import (
  "bytes"
  "encoding/json"
  "fmt"
  "frog-go/internal/models"
  "io"
  "net/http"
  "strconv"
  "testing"
  "time"
)

// TestUser 测试用例
// $ go test -run TestUser
// Base URL
const baseURL = "http://127.0.0.1:8080/api/user"

// TestUserCRUD 测试用户的增删改查功能
func TestUser(t *testing.T) {
  var createdUserID int

  // ---------------------- 创建用户 (Create) ----------------------
  t.Run("CreateUser", func(t *testing.T) {
    wechat := "test_wechat"
    address := "Test Address"

    user := models.User{
      Name:    "测试用户",
      Email:   "testuser" + time.DateTime + "@example.com",
      Phone:   "1234567890",
      Wechat:  &wechat,
      Address: &address,
    }

    data, err := json.Marshal(user)
    if err != nil {
      t.Fatalf("Error marshalling data: %v", err)
    }

    req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(data))
    if err != nil {
      t.Fatalf("Error creating request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
      t.Fatalf("Error sending request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
      bodyBytes, err := io.ReadAll(resp.Body)
      if err != nil {
        t.Fatalf("Error reading response body: %v", err)
      }
      t.Fatalf("Expected status Created, got %v, body: %s", resp.StatusCode, bodyBytes)
    }

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      t.Fatalf("Error decoding response: %v", err)
    }

    createdUserID = int(result["data"].(map[string]interface{})["id"].(float64))
    t.Logf("User created successfully with ID: %d", createdUserID)
  })

  // ---------------------- 查询用户 (Read) ----------------------
  t.Run("GetUser", func(t *testing.T) {
    req, err := http.NewRequest("GET", baseURL+"/"+strconv.Itoa(createdUserID), nil)
    if err != nil {
      t.Fatalf("Error creating GET request: %v", err)
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
      t.Fatalf("Error sending GET request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
      t.Fatalf("Expected status OK, got %v", resp.StatusCode)
    }

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      t.Fatalf("Error decoding GET response: %v", err)
    }

    t.Logf("Fetched user: %+v", result)
  })

  // ---------------------- 更新用户 (Update) ----------------------
  t.Run("UpdateUser", func(t *testing.T) {
    updatedWechat := "updated_wechat"
    updatedAddress := "Updated Address"

    updatedUser := models.User{
      Id:      createdUserID,
      Name:    "更新用户",
      Email:   "updated@example.com",
      Phone:   "9876543210",
      Wechat:  &updatedWechat,
      Address: &updatedAddress,
    }

    data, err := json.Marshal(updatedUser)
    if err != nil {
      t.Fatalf("Error marshalling updated data: %v", err)
    }

    req, err := http.NewRequest("PUT", baseURL, bytes.NewBuffer(data))
    if err != nil {
      t.Fatalf("Error creating PUT request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
      t.Fatalf("Error sending PUT request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
      bodyBytes, err := io.ReadAll(resp.Body)
      if err != nil {
        t.Fatalf("Error reading response body: %v", err)
      }
      t.Fatalf("Expected status OK, got %v, body: %s", resp.StatusCode, bodyBytes)
    }

    t.Logf("User updated successfully: ID %d", createdUserID)
  })

  // ---------------------- 删除用户 (Delete) ----------------------
  t.Run("DeleteUser", func(t *testing.T) {
    req, err := http.NewRequest("DELETE", baseURL+"/"+strconv.Itoa(createdUserID), nil)
    if err != nil {
      t.Fatalf("Error creating DELETE request: %v", err)
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
      t.Fatalf("Error sending DELETE request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
      bodyBytes, err := io.ReadAll(resp.Body)
      if err != nil {
        t.Fatalf("Error reading response body: %v", err)
      }
      t.Fatalf("Expected status OK, got %v, body: %s", resp.StatusCode, bodyBytes)
    }

    t.Logf("User deleted successfully: ID %d", createdUserID)
  })

  // ---------------------- 查询用户列表 (List) ----------------------
  t.Run("ListUsers", func(t *testing.T) {
    // 1. 构建请求 URL
    url := baseURL + "/list?page=1&size=5&fields=name,email&condition=email%20like%20%27a%25%27"
    // fmt.Println("Request URL:", url)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
      t.Fatalf("Error creating GET request: %v", err)
    }

    // 2. 创建 HTTP 客户端并发送请求
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
      t.Fatalf("Error sending GET request: %v", err)
    }
    defer resp.Body.Close()

    // 读取响应内容
    // body, err := io.ReadAll(resp.Body)
    // fmt.Println(url, "Response Body:", string(body), resp.StatusCode)

    if err != nil {
      fmt.Println("Error reading response body:", err)
      return
    }
    // 打印响应的原始 JSON 内容

    // 3. 校验状态码
    if resp.StatusCode != http.StatusOK {
      t.Fatalf("Expected status OK, got %v", resp.StatusCode)
    }

    // 4. 解析返回数据
    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      t.Fatalf("Error decoding GET response: %v", err)
    }

    // 5. 校验返回数据格式
    data, ok := result["data"].(map[string]interface{})["users"].([]interface{})
    if !ok {
      t.Fatalf("Expected 'data' to be a list, got %T", result["data"])
    }

    total, ok := result["data"].(map[string]interface{})["total"].(float64) // JSON 数字默认解析为 float64
    if !ok {
      t.Fatalf("Expected 'total' to be a number, got %T", result["total"])
    }

    // 6. 打印返回结果
    t.Logf("Fetched user list: %+v", data)
    t.Logf("Total users: %v", total)

    // 7. 校验数据不为空（根据实际业务需求调整）
    if len(data) == 0 {
      t.Errorf("Expected at least one user in the list, but got zero")
    }
  })

  // ---------------------- 分组查询测试 (List) ----------------------
  t.Run("GroupUsers", func(t *testing.T) {
    // 1. 构建请求 URL，包含分组字段
    url := baseURL + "/group?field=name,email" // 传入要分组的字段
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
      t.Fatalf("Error creating GET request: %v", err)
    }

    // 2. 创建 HTTP 客户端并发送请求
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
      t.Fatalf("Error sending GET request: %v", err)
    }
    defer resp.Body.Close()

    // 3. 校验状态码
    if resp.StatusCode != http.StatusOK {
      t.Fatalf("Expected status OK, got %v", resp.StatusCode)
    }

    // 4. 解析返回数据
    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      t.Fatalf("Error decoding GET response: %v", err)
    }

    // 5. 校验返回数据格式
    groupedData, ok := result["data"].(map[string]interface{})
    if !ok {
      t.Fatalf("Expected 'users' to be a map, got %T", result["data"])
    }

    // 6. 打印返回结果
    t.Logf("Grouped data: %+v", groupedData)

    // 7. 校验返回数据格式
    if len(groupedData) == 0 {
      t.Errorf("Expected non-empty grouped data, but got zero")
    }
  })

}
