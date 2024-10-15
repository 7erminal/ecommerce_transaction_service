package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type AlterCreatedByColumn_20241015_065319 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AlterCreatedByColumn_20241015_065319{}
	m.Created = "20241015_065319"

	migration.Register("AlterCreatedByColumn_20241015_065319", m)
}

// Run the migrations
func (m *AlterCreatedByColumn_20241015_065319) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("ALTER TABLE orders ADD CONSTRAINT fk_userid FOREIGN KEY (created_by) REFERENCES users(user_id);")
}

// Reverse the migrations
func (m *AlterCreatedByColumn_20241015_065319) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
