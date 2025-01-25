/**
 * HTTP code Helper
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package code

import "errors"

const (
  Success     = 200 // 成功
  Error       = 400 // 通用错误
  ServerError = 500 // 服务错误
)

var (
  ErrEmailExists    = errors.New("邮箱账号已存在")
  ErrUserNotFound   = errors.New("用户不存在")
  ErrUserInvalid    = errors.New("用户信息无效")
  ErrUserUpdateFail = errors.New("用户更新失败")
  ErrUserCreateFail = errors.New("用户创建失败")
  ErrDatabase       = errors.New("数据库操作错误")
)
