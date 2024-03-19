package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type StatusTypes_20240317_090120 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &StatusTypes_20240317_090120{}
	m.Created = "20240317_090120"

	migration.Register("StatusTypes_20240317_090120", m)
}

// Run the migrations
func (m *StatusTypes_20240317_090120) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE status_types(`status_type_id` int(11) NOT NULL AUTO_INCREMENT,`status` varchar(50) NOT NULL,`active` int(11) DEFAULT 1,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT 1,`modified_by` int(11) DEFAULT 1,PRIMARY KEY (`status_type_id`))")
}

// Reverse the migrations
func (m *StatusTypes_20240317_090120) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `status_types`")
}
