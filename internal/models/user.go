/**
 * User Model Definition
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */
package models

import (
  "time"
)

type User struct {
  Id          int       `json:"id" db:"id"`
  Name        string    `json:"name" db:"name"`
  Email       string    `json:"email" db:"email"`
  Phone       string    `json:"phone" db:"phone"`
  Wechat      *string   `json:"wechat" db:"wechat"`
  Address     *string   `json:"address" db:"address"`
  CreatedTime time.Time `json:"created_time" db:"created_time"`
  UpdatedTime time.Time `json:"updated_time" db:"updated_time"`
}
