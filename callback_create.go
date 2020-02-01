package sorm

import (
	"fmt"
	"strings"
)

// Define callbacks for creating
func init() {
	DefaultCallback.Create().Register(BeginTransaction, beginTransactionCallback)
	DefaultCallback.Create().Register(BeforeCreate, beforeCreateCallback)
	DefaultCallback.Create().Register(SaveBeforeAssociations, saveBeforeAssociationsCallback)
	DefaultCallback.Create().Register(UpdateTimeStamp, updateTimeStampForCreateCallback)
	DefaultCallback.Create().Register(create, createCallback)
	DefaultCallback.Create().Register(ForceReloadAfterCreate, forceReloadAfterCreateCallback)
	DefaultCallback.Create().Register(SaveAfterAssociations, saveAfterAssociationsCallback)
	DefaultCallback.Create().Register(AfterCreate, afterCreateCallback)
	DefaultCallback.Create().Register(CommitOrRollbackTransaction, commitOrRollbackTransactionCallback)
}

// beforeCreateCallback will invoke `BeforeSave`, `BeforeCreate` method before creating
func beforeCreateCallback(scope *Scope) {
	if !scope.HasError() {
		scope.CallMethod("BeforeSave")
	}
	if !scope.HasError() {
		scope.CallMethod("BeforeCreate")
	}
}

// updateTimeStampForCreateCallback will set `CreatedAt`, `UpdatedAt` when creating
func updateTimeStampForCreateCallback(scope *Scope) {
	if !scope.HasError() {
		now := scope.db.nowFunc()

		if createdAtField, ok := scope.FieldByName("CreatedAt"); ok {
			if createdAtField.IsBlank {
				createdAtField.Set(now)
			}
		}

		if updatedAtField, ok := scope.FieldByName("UpdatedAt"); ok {
			if updatedAtField.IsBlank {
				updatedAtField.Set(now)
			}
		}
	}
}

// createCallback the callback used to insert data into database
func createCallback(scope *Scope) {
	if !scope.HasError() {
		defer scope.trace(NowFunc())

		var (
			columns, placeholders        []string
			blankColumnsWithDefaultValue []string
		)

		for _, field := range scope.Fields() {
			if scope.changeableField(field) {
				if field.IsNormal && !field.IsIgnored {
					if field.IsBlank && field.HasDefaultValue {
						blankColumnsWithDefaultValue = append(blankColumnsWithDefaultValue, scope.Quote(field.DBName))
						scope.InstanceSet("sorm:blank_columns_with_default_value", blankColumnsWithDefaultValue)
					} else if !field.IsPrimaryKey || !field.IsBlank {
						columns = append(columns, scope.Quote(field.DBName))
						placeholders = append(placeholders, scope.AddToVars(field.Field.Interface()))
					}
				} else if field.Relationship != nil && field.Relationship.Kind == "belongs_to" {
					for _, foreignKey := range field.Relationship.ForeignDBNames {
						if foreignField, ok := scope.FieldByName(foreignKey); ok && !scope.changeableField(foreignField) {
							columns = append(columns, scope.Quote(foreignField.DBName))
							placeholders = append(placeholders, scope.AddToVars(foreignField.Field.Interface()))
						}
					}
				}
			}
		}

		var (
			returningColumn = "*"
			quotedTableName = scope.QuotedTableName()
			primaryField    = scope.PrimaryField()
			extraOption     string
			insertModifier  string
		)

		if str, ok := scope.Get("sorm:insert_option"); ok {
			extraOption = fmt.Sprint(str)
		}
		if str, ok := scope.Get("sorm:insert_modifier"); ok {
			insertModifier = strings.ToUpper(fmt.Sprint(str))
			if insertModifier == "INTO" {
				insertModifier = ""
			}
		}

		if primaryField != nil {
			returningColumn = scope.Quote(primaryField.DBName)
		}

		lastInsertIDOutputInterstitial := scope.Dialect().LastInsertIDOutputInterstitial(quotedTableName, returningColumn, columns)
		var lastInsertIDReturningSuffix string
		if lastInsertIDOutputInterstitial == "" {
			lastInsertIDReturningSuffix = scope.Dialect().LastInsertIDReturningSuffix(quotedTableName, returningColumn)
		}

		if len(columns) == 0 {
			scope.Raw(fmt.Sprintf(
				"INSERT%v INTO %v %v%v%v",
				addExtraSpaceIfExist(insertModifier),
				quotedTableName,
				scope.Dialect().DefaultValueStr(),
				addExtraSpaceIfExist(extraOption),
				addExtraSpaceIfExist(lastInsertIDReturningSuffix),
			))
		} else {
			scope.Raw(fmt.Sprintf(
				"INSERT%v INTO %v (%v)%v VALUES (%v)%v%v",
				addExtraSpaceIfExist(insertModifier),
				scope.QuotedTableName(),
				strings.Join(columns, ","),
				addExtraSpaceIfExist(lastInsertIDOutputInterstitial),
				strings.Join(placeholders, ","),
				addExtraSpaceIfExist(extraOption),
				addExtraSpaceIfExist(lastInsertIDReturningSuffix),
			))
		}

		// execute create sql: no primaryField
		if primaryField == nil {
			if result, err := scope.SQLDB().Exec(scope.SQL, scope.SQLVars...); scope.Err(err) == nil {
				// set rows affected count
				scope.db.RowsAffected, _ = result.RowsAffected()

				// set primary value to primary field
				if primaryField != nil && primaryField.IsBlank {
					if primaryValue, err := result.LastInsertId(); scope.Err(err) == nil {
						scope.Err(primaryField.Set(primaryValue))
					}
				}
			}
			return
		}

		// execute create sql: lastInsertID implemention for majority of dialects
		if lastInsertIDReturningSuffix == "" && lastInsertIDOutputInterstitial == "" {
			if result, err := scope.SQLDB().Exec(scope.SQL, scope.SQLVars...); scope.Err(err) == nil {
				// set rows affected count
				scope.db.RowsAffected, _ = result.RowsAffected()

				// set primary value to primary field
				if primaryField != nil && primaryField.IsBlank {
					if primaryValue, err := result.LastInsertId(); scope.Err(err) == nil {
						scope.Err(primaryField.Set(primaryValue))
					}
				}
			}
			return
		}

		// execute create sql: dialects with additional lastInsertID requirements (currently postgres & mssql)
		if primaryField.Field.CanAddr() {
			if err := scope.SQLDB().QueryRow(scope.SQL, scope.SQLVars...).Scan(primaryField.Field.Addr().Interface()); scope.Err(err) == nil {
				primaryField.IsBlank = false
				scope.db.RowsAffected = 1
			}
		} else {
			scope.Err(ErrUnaddressable)
		}
		return
	}
}

// forceReloadAfterCreateCallback will reload columns that having default value, and set it back to current object
func forceReloadAfterCreateCallback(scope *Scope) {
	if blankColumnsWithDefaultValue, ok := scope.InstanceGet("sorm:blank_columns_with_default_value"); ok {
		db := scope.DB().New().Table(scope.TableName()).Select(blankColumnsWithDefaultValue.([]string))
		for _, field := range scope.Fields() {
			if field.IsPrimaryKey && !field.IsBlank {
				db = db.Where(fmt.Sprintf("%v = ?", field.DBName), field.Field.Interface())
			}
		}
		db.Scan(scope.Value)
	}
}

// afterCreateCallback will invoke `AfterCreate`, `AfterSave` method after creating
func afterCreateCallback(scope *Scope) {
	if !scope.HasError() {
		scope.CallMethod("AfterCreate")
	}
	if !scope.HasError() {
		scope.CallMethod("AfterSave")
	}
}
