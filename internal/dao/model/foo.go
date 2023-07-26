package model

type Foo struct {
	ID   int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name string `json:"name" gorm:"type:varchar(255);not null"`
	Age  int64  `json:"age" gorm:"type:int"`
}
