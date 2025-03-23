package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type AddColumnToOrderItemsTable_20250322_102731 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AddColumnToOrderItemsTable_20250322_102731{}
	m.Created = "20250322_102731"

	migration.Register("AddColumnToOrderItemsTable_20250322_102731", m)
}

// Run the migrations
func (m *AddColumnToOrderItemsTable_20250322_102731) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("ALTER TABLE order_items ADD COLUMN item_status int DEFAULT NULL AFTER quantity, ADD FOREIGN KEY (item_status) REFERENCES status(status_id);")
}

// Reverse the migrations
func (m *AddColumnToOrderItemsTable_20250322_102731) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
