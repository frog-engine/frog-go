/**
 * Repository Utilities
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/frog-engine/frog-go/pkg/logger"
)

// GetStructFieldsByTag 根据给定的标签名获取结构体中的字段名称，并返回一个 map，
// 其中 db 标签值为 key，字段名为 value。
 func GetStructFieldsByTag(v interface{}, tagName string) (map[string]string, error) {
	 // 获取结构体的反射类型
	 val := reflect.ValueOf(v)
	 if val.Kind() != reflect.Struct {
		 return nil, fmt.Errorf("expected a struct type, got %s", val.Kind())
	 }
 
	 // 初始化结果的 map
	 fieldsMap := make(map[string]string)
	 
	 // 遍历结构体的所有字段
	 for i := 0; i < val.NumField(); i++ {
		 field := val.Type().Field(i)
		 // 获取指定标签的值
		 tagValue := field.Tag.Get(tagName)
		 if tagValue != "" {
			 // 将 db 标签值作为 key，字段名作为 value
			 fieldsMap[tagValue] = field.Name
		 }
	 }
 
	 return fieldsMap, nil
 }
 
 
 // MapScannerToEntity 将数据从 *sql.Row 或 *sql.Rows 映射到实体
 func MapScannerToEntity[T any](scanner interface {
	 Scan(dest ...interface{}) error
 }, entity *T) (*T, error) {
	 // 确保 entity 是指针且指向结构体
	 entityValue := reflect.ValueOf(entity)
	 if entityValue.Kind() != reflect.Ptr || entityValue.Elem().Kind() != reflect.Struct {
		 return nil, errors.New("entity 必须是指针且指向结构体")
	 }
 
	 // 获取结构体字段数量
	 numFields := entityValue.Elem().NumField()
	 scanArgs := make([]interface{}, numFields)
 
	 // 为每个字段准备扫描地址
	 for i := 0; i < numFields; i++ {
		 field := entityValue.Elem().Field(i)
		 if field.CanAddr() {
			 scanArgs[i] = field.Addr().Interface()
		 } else {
			 // 如果字段不可寻址，创建一个占位变量
			 var placeholder interface{}
			 scanArgs[i] = &placeholder
		 }
	 }
 
	 // 处理 scanner 类型
	 if err := scanner.Scan(scanArgs...); err != nil {
		 logger.Printf("扫描数据到结构体时发生错误: %v", err)
		 return nil, fmt.Errorf("无法扫描数据到实体: %w", err)
	 }
 
	 return entity, nil
 }
 
 
 // StructToSlice 提取结构体字段值并生成 SQL 的 SET 子句
 func StructToSlice(entity interface{}, tagName string, excludeFields ...string) (string, []interface{}, error) {
	 v := reflect.ValueOf(entity)
	 t := reflect.TypeOf(entity)
 
	 // 确保输入是结构体或指向结构体的指针
	 if v.Kind() == reflect.Ptr {
		 v = v.Elem()
		 t = t.Elem()
	 }
 
	 if v.Kind() != reflect.Struct {
		 return "", nil, errors.New("输入必须是结构体或指向结构体的指针")
	 }
 
	 var setClauses []string
	 var args []interface{}
 
	 // 将 excludeFields 转换为 map，方便快速查找
	 excludeMap := make(map[string]struct{}, len(excludeFields))
	 for _, field := range excludeFields {
		 excludeMap[field] = struct{}{}
	 }
 
	 // 遍历结构体的字段
	 for i := 0; i < t.NumField(); i++ {
		 field := t.Field(i)
		 dbTag := field.Tag.Get(tagName)
 
		 // 跳过没有 db 标签或需要排除的字段
		 if dbTag == "" {
			 continue
		 }
		 if _, excluded := excludeMap[dbTag]; excluded {
			 continue
		 }
 
		 // 添加到 SET 子句
		 setClauses = append(setClauses, dbTag+" = ?")
		 args = append(args, v.Field(i).Interface())
	 }
 
	 // 合并 SET 子句为字符串
	 return strings.Join(setClauses, ", "), args, nil
 }
 
 // StructToSQL 动态提取结构体字段和值，生成字段名和值
 func StructToSQL(entity interface{}, excludeFields ...string) (columns []string, placeholders []string, values []interface{}, err error) {
	 v := reflect.ValueOf(entity)
	 t := reflect.TypeOf(entity)
 
	 // 确保输入是结构体或指向结构体的指针
	 if v.Kind() == reflect.Ptr {
		 v = v.Elem()
		 t = t.Elem()
	 }
 
	 if v.Kind() != reflect.Struct {
		 return nil, nil, nil, fmt.Errorf("input must be a struct or a pointer to a struct")
	 }
 
	 // 将排除字段转换为 map 方便查找
	 excludeMap := make(map[string]struct{}, len(excludeFields))
	 for _, field := range excludeFields {
		 excludeMap[field] = struct{}{}
	 }
 
	 // 遍历结构体字段
	 for i := 0; i < t.NumField(); i++ {
		 field := t.Field(i)
		 dbTag := field.Tag.Get("db") // 假设 db 标签用于字段名映射
 
		 // 跳过没有 db 标签或需要排除的字段
		 if dbTag == "" {
			 continue
		 }
		 if _, excluded := excludeMap[dbTag]; excluded {
			 continue
		 }
 
		 // 获取字段值
		 fieldValue := v.Field(i)
 
		 // 添加列名和占位符
		 columns = append(columns, dbTag)
		 placeholders = append(placeholders, "?")
		 values = append(values, fieldValue.Interface())
	 }
 
	 return columns, placeholders, values, nil
 }
 
 // StructToFields 使用反射获取结构体的字段，返回字段名（db标签）以及字段值
 func StructToFields(entity interface{}, excludeFields ...string) (columns []string, values []interface{}, err error) {
	 v := reflect.ValueOf(entity)
	 t := reflect.TypeOf(entity)
 
	 // 确保输入是结构体或指向结构体的指针
	 if v.Kind() == reflect.Ptr {
		 v = v.Elem()
		 t = t.Elem()
	 }
 
	 if v.Kind() != reflect.Struct {
		 return nil, nil, fmt.Errorf("input must be a struct or a pointer to a struct")
	 }
 
	 // 将排除字段转换为 map 方便查找
	 excludeMap := make(map[string]struct{}, len(excludeFields))
	 for _, field := range excludeFields {
		 excludeMap[field] = struct{}{}
	 }
 
	 // 遍历结构体字段
	 for i := 0; i < t.NumField(); i++ {
		 field := t.Field(i)
		 dbTag := field.Tag.Get("db") // 假设 db 标签用于字段名映射
 
		 // 跳过没有 db 标签或需要排除的字段
		 if dbTag == "" {
			 continue
		 }
		 if _, excluded := excludeMap[dbTag]; excluded {
			 continue
		 }
 
		 // 获取字段值
		 fieldValue := v.Field(i)
 
		 // 添加列名和对应值
		 columns = append(columns, dbTag)
		 values = append(values, fieldValue.Interface())
	 }
 
	 return columns, values, nil
 }
 
// MapRowToStruct 根据字段名映射 row 中的值到结构体字段
func MapRowToStruct(row map[string]interface{}, entity interface{}) error {
	// 确保 entity 是一个指针
	entityVal := reflect.ValueOf(entity)
	if entityVal.Kind() != reflect.Ptr || entityVal.IsNil() {
			return fmt.Errorf("entity must be a non-nil pointer to a struct")
	}

	entityVal = entityVal.Elem()
	entityType := entityVal.Type()

	// 遍历 struct 字段，直接根据 db 标签将字段名映射到字段索引
	for i := 0; i < entityType.NumField(); i++ {
			field := entityType.Field(i)
			dbTag := field.Tag.Get("db")
			if dbTag == "" {
					continue
			}

			if value, exists := row[dbTag]; exists {
					fieldVal := entityVal.Field(i)

					// 跳过不可设置的字段
					if !fieldVal.CanSet() || value == nil {
							continue
					}

					val := reflect.ValueOf(value)

					// 如果是指针类型字段，确保初始化指针
					if fieldVal.Kind() == reflect.Ptr {
							if fieldVal.IsNil() {
									fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
							}
							fieldVal = fieldVal.Elem() // 获取指针指向的值
					}

					// 只处理基础类型字段
					if val.Type().AssignableTo(fieldVal.Type()) {
							fieldVal.Set(val)
					} else if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
					}
			}
	}

	return nil
}

 
 // MapRowToStructReflect 使用反射根据字段动态填充结构体
 func MapRowToStructReflect(row map[string]interface{}, entity interface{}, fields []string) error {
	 // 确保 entity 是一个非空指针
	 entityValue := reflect.ValueOf(entity)
	 if entityValue.Kind() != reflect.Ptr || entityValue.IsNil() {
		 return fmt.Errorf("entity must be a non-nil pointer to a struct")
	 }
 
	 entityValue = entityValue.Elem()
	 entityType := entityValue.Type()
 
	 // 如果 fields 为空，则获取所有结构体字段名称
	 if len(fields) == 0 {
		 for i := 0; i < entityType.NumField(); i++ {
			 fields = append(fields, entityType.Field(i).Name)
		 }
	 }
 
	 // 遍历字段名称，按名称从 row 中获取值
	 for _, fieldName := range fields {
		 // 获取字段
		 structField := entityValue.FieldByName(fieldName)
		 if !structField.IsValid() || !structField.CanSet() {
			 logger.Printf("Field %s is invalid or cannot be set", fieldName)
			 continue
		 }
 
		 // 获取字段的 db 标签（如果有）
		 fieldInfo, _ := entityType.FieldByName(fieldName)
		 dbTag := fieldInfo.Tag.Get("db")
		 key := dbTag
		 if key == "" {
			 // 如果没有 db 标签，默认用字段名
			 key = fieldName
		 }
 
		 // 从 row 获取值
		 value, exists := row[key]
		 if !exists {
			 logger.Printf("Key %s not found in row", key)
			 continue
		 }
 
		 val := reflect.ValueOf(value)
		 if !val.IsValid() {
			 logger.Printf("Value for key %s is invalid", key)
			 continue
		 }
 
		 // 处理 nil 值（指针、接口、切片、map、chan、func）
		 switch val.Kind() {
		 case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
			 if val.IsNil() {
				 logger.Printf("Value for key %s is nil", key)
				 continue
			 }
		 }
 
		 // 类型匹配或转换赋值
		 if val.Type().AssignableTo(structField.Type()) {
			 structField.Set(val)
		 } else if val.Type().ConvertibleTo(structField.Type()) {
			 structField.Set(val.Convert(structField.Type()))
		 } else if structField.Kind() == reflect.Ptr && val.Type().AssignableTo(structField.Type().Elem()) {
			 // 基本类型转指针类型
			 ptrValue := reflect.New(structField.Type().Elem())
			 ptrValue.Elem().Set(val)
			 structField.Set(ptrValue)
		 } else {
			 logger.Printf("type mismatch for field %s: expected %s, got %s", fieldName, structField.Type(), val.Type())
			 return nil
		 }
	 }
	 return nil
 }