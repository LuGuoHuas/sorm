# sorm
A Simple ORM

## example
```go
// connecting to database
var db = sorm.Open("postgres", "host=myhost port=myport user=gorm dbname=gorm password=mypassword")

// Declaring Models
type ObjectModel1 struct {
    Field1  sorm.Varchar    `json:"field_1"`
    Field2  sorm.Interger   `json:"field_2"`
}

type ObjectModel2 struct {
    Field1  sorm.Interger   `json:"field_1"`
    Field2  sorm.Varchar    `json:"field_2"`
}

var obj1 = ObjectModel1{
    Field1: "field_1",
    Field2: 2,
}

var obj2 = ObjectModel2{
    Field1: 1,
    Field2: "filed_2",
}


// Create Record
var result, err = db.obj1.Create()

// Query
var result, err = db.obj1.Find()
var result, err = db.obj1.Filter(sorm.Eq(obj1.Field1, "field")).Group(obj.Field1).Limit(10).Find()

// Update
var result, err = db.obj1.Filter(sorm.Lte(obj1.Field2, 3), sorm.Eq(obj1.Field1, "field")).Update()

// Join
var result, err = db.obj1.Join(obj2).On(obj1.Field1, obj2.Field2).On(obj1.Field2, obj2.Field1).Filter(sorm.Lte(obj1.Field2, 2)).Find()
```
