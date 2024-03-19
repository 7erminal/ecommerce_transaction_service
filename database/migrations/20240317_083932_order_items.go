package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type OrderItems_20240317_083932 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &OrderItems_20240317_083932{}
	m.Created = "20240317_083932"

	migration.Register("OrderItems_20240317_083932", m)
}

// Run the migrations
func (m *OrderItems_20240317_083932) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE order_items(`order_item_id` int(11) NOT NULL AUTO_INCREMENT,`order_id` int(11) NOT NULL,`item_id` int(11) NOT NULL,`quantity` int(11) NOT NULL,`order_date` datetime NOT NULL,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT 1,`modified_by` int(11) DEFAULT 1,PRIMARY KEY (`order_item_id`), FOREIGN KEY (order_id) REFERENCES orders(order_id) ON UPDATE CASCADE ON DELETE NO ACTION, FOREIGN KEY (item_id) REFERENCES items(item_id) ON UPDATE CASCADE ON DELETE NO ACTION)")
}

// Reverse the migrations
func (m *OrderItems_20240317_083932) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `order_items`")
}
