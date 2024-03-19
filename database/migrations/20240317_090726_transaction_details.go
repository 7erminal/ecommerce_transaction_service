package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type TransactionDetails_20240317_090726 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &TransactionDetails_20240317_090726{}
	m.Created = "20240317_090726"

	migration.Register("TransactionDetails_20240317_090726", m)
}

// Run the migrations
func (m *TransactionDetails_20240317_090726) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE transaction_details(`transaction_detail_id` int(11) NOT NULL AUTO_INCREMENT,`transaction_id` int(11) NOT NULL,`amount` int(11) DEFAULT NULL,`comment` varchar(255) DEFAULT NULL,`sender_account_number` varchar(255) DEFAULT NULL,`recipient_account_number` varchar(255) DEFAULT NULL,`status_code` varchar(20) DEFAULT NULL,`status_message` varchar(255) DEFAULT NULL,`sender_id` int(11) DEFAULT NULL,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT 1,`modified_by` int(11) DEFAULT 1,`active` int(11) DEFAULT 1,PRIMARY KEY (`transaction_detail_id`), FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id) ON UPDATE CASCADE ON DELETE NO ACTION)")
}

// Reverse the migrations
func (m *TransactionDetails_20240317_090726) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `transaction_details`")
}
