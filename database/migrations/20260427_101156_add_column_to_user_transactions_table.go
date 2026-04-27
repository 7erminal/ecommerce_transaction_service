package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type AddColumnToUserTransactionsTable_20260427_101156 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AddColumnToUserTransactionsTable_20260427_101156{}
	m.Created = "20260427_101156"

	migration.Register("AddColumnToUserTransactionsTable_20260427_101156", m)
}

// Run the migrations
func (m *AddColumnToUserTransactionsTable_20260427_101156) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("ALTER TABLE user_transactions ADD COLUMN reference varChar(200) default null after commission")
}

// Reverse the migrations
func (m *AddColumnToUserTransactionsTable_20260427_101156) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
