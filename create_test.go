package sorm

import (
	"math"
	"testing"
	"time"
)

type User struct {
	ID                int64      `json:"id" sorm:"column:id;primary_key:true"`
	Age               int64      `json:"age" sorm:"column:age"`
	Name              string     `json:"name" sorm:"column:name"`
	Email             string     `json:"email" sorm:"column:email"`
	Birthday          *time.Time `json:"birthday" sorm:"column:birthday"`
	ShippingAddressID int64      `json:"shipping_address_d" sorm:"column:shipping_address_d"`
	Latitude          float64    `json:"latitude" sorm:"column:latitude"`
	CompanyID         *int       `json:"company_id" sorm:"column:company_id"`
	Status            bool       `json:"status" sorm:"column:status"`
	Model
}

func TestCreate(t *testing.T) {
	var float = math.Pi
	var now = time.Now()
	var user = Make(&User{Name: "user_1", Age: 18, Birthday: &now, Latitude: float, Status: true}).(*User)
	var db *DB
	if d, ok := TestData.Load(TestConnectKey); ok {
		db = d.(*DB)
	}

	if !db.NewRecord(user) {
		t.Error("User should be new record before create")
	}

	if db.Save(&user).RowsAffected != 1 {
		t.Error("There should be one record be affected when create record")
	}

	if db.NewRecord(user) {
		t.Error("User should not new record after save")
	}

	var newUser User
	if err := db.First(&newUser, user.ID).Error; err != nil {
		t.Errorf("No error should happen, but got %v", err)
	} else if newUser.Age != user.Age {
		t.Errorf("User's Age should be saved (int)")
	} else if newUser.Latitude != float {
		t.Errorf("Float64 should not be changed after save")
	}

	db.Table(user).Update(user.Name, "user_2")
	db.First(&user, user.ID)
}
