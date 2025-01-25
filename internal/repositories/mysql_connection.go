/**
 * ConnectDatabase Definition
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */
package repositories

import (
  "database/sql"
  "fmt"

  "github.com/frog-engine/frog-go/pkg/logger"
  _ "github.com/go-sql-driver/mysql"
)

type DatabaseConfig struct {
  Host     string
  Port     int
  User     string
  Password string
  DBName   string
}

func ConnectDatabase(dbConfig *DatabaseConfig) (*sql.DB, error) {
  if dbConfig == nil {
    return nil, fmt.Errorf("数据库配置不能为空")
  }

  // 构建数据库连接字符串
  dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)

  // 连接数据库
  db, err := sql.Open("mysql", dsn)
  if err != nil {
    logger.Printf("无法连接到数据库，错误信息：%v", err)
    return nil, fmt.Errorf("无法连接到数据库: %w", err)
  }

  // 确保数据库连接可用
  if err := db.Ping(); err != nil {
    logger.Printf("数据库不可用，错误信息：%v", err)
    db.Close() // 连接失败时关闭数据库连接
    return nil, fmt.Errorf("数据库不可用: %w", err)
  }

  logger.Println("成功连接到数据库")

  return db, nil
}
