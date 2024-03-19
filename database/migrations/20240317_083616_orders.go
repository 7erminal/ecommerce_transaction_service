package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Orders_20240317_083616 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Orders_20240317_083616{}
	m.Created = "20240317_083616"

	migration.Register("Orders_20240317_083616", m)
}

// Run the migrations
func (m *Orders_20240317_083616) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE orders(`order_id` int(11) NOT NULL AUTO_INCREMENT,`quantity` int(11) NOT NULL,`cost` int(11) NOT NULL,`currency` int(11) NOT NULL,`order_date` datetime NOT NULL,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT 1,`modified_by` int(11) DEFAULT 1,PRIMARY KEY (`order_id`))")
}

// Reverse the migrations
func (m *Orders_20240317_083616) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `orders`")
}
