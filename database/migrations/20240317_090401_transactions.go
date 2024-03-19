package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Transactions_20240317_084831 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Transactions_20240317_084831{}
	m.Created = "20240317_084831"

	migration.Register("Transactions_20240317_084831", m)
}

// Run the migrations
func (m *Transactions_20240317_084831) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE transactions(`transaction_id` int(11) NOT NULL AUTO_INCREMENT,`order_id` int(11) NOT NULL,`amount` int(11) DEFAULT NULL,`transacting_currency` int(11) NOT NULL,`status_id` int(11) NOT NULL,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT 1,`modified_by` int(11) DEFAULT 1,`active` int(11) DEFAULT 1,PRIMARY KEY (`transaction_id`), FOREIGN KEY (order_id) REFERENCES orders(order_id) ON UPDATE CASCADE ON DELETE NO ACTION, FOREIGN KEY (status_id) REFERENCES status_types(status_type_id) ON UPDATE CASCADE ON DELETE NO ACTION)")
}

// Reverse the migrations
func (m *Transactions_20240317_084831) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `transactions`")
}
