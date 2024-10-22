package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type AddColumnToOrdersTable_20241022_083054 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AddColumnToOrdersTable_20241022_083054{}
	m.Created = "20241022_083054"

	migration.Register("AddColumnToOrdersTable_20241022_083054", m)
}

// Run the migrations
func (m *AddColumnToOrdersTable_20241022_083054) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("ALTER TABLE orders ADD COLUMN order_number int DEFAULT NULL AFTER order_id;")
}

// Reverse the migrations
func (m *AddColumnToOrdersTable_20241022_083054) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
