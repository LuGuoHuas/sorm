[![GitHub license](https://img.shields.io/github/license/LuGHuaaa/sorm)](https://github.com/LuGHuaaa/sorm/blob/master/LICENSE)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/LuGHuaaa/sorm)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/LuGHuaaa/sorm/Go%20-%20Test)



# sorm
A Simple ORM

## example
```go
// connecting to database
var db = sorm.Open("postgres", "host=myhost port=myport user=gorm dbname=gorm password=mypassword")
var err error

// Declaring Models
type ObjectModel1 struct {
    sorm.Model
    Field1  sorm.Varchar    `json:"field_1"`
    Field2  sorm.Interger   `json:"field_2"`
}

type ObjectModel2 struct {
    sorm.Model
    Field1  sorm.Interger   `json:"field_1"`
    Field2  sorm.Varchar    `json:"field_2"`
}

var obj1 = sorm.Make(&ObjectModel1{
    Field1: "field_1",
    Field2: 2,
}).(*ObjectModel1)

var obj2 = sorm.Make(&ObjectModel2{
    Field1: 1,
    Field2: "filed_2",
}).(*ObjectModel2)


// Create Record
err = db.Create(obj1).Error

// Query
err = db.Table(obj1).Find().Error
err = db.Table(obj1).Filter(sorm.Eq(obj1.Field1, "field")).Group(obj.Field1).Limit(10).Find(&result).Error

// Update
err = db.Table(obj1).Filter(sorm.Lte(obj1.Field2, 3), sorm.Eq(obj1.Field1, "field")).Update(obj1.Field2, "field").Error

// Delete
err = db.Table(obj1).Delete()

// Join
err = db.Table(obj1).Join(obj2).On(obj1.Field1, obj2.Field2).On(obj1.Field2, obj2.Field1).Filter(sorm.Lte(obj1.Field2, 2)).Find().Error
```
