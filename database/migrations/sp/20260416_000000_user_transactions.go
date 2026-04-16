package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type UserTransactions_20260416_000000 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserTransactions_20260416_000000{}
	m.Created = "20260416_000000"

	migration.Register("UserTransactions_20260416_000000", m)
}

// Run the migrations
func (m *UserTransactions_20260416_000000) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE user_transactions(
		` + "`transaction_id`" + ` VARCHAR(255) NOT NULL,
		` + "`service_id`" + ` INT(11) NOT NULL,
		` + "`request_id`" + ` INT(11) NOT NULL,
		` + "`transaction_customer_reference`" + ` INT(11) DEFAULT NULL,
		` + "`amount`" + ` DOUBLE DEFAULT NULL,
		` + "`transacting_currency`" + ` VARCHAR(255) DEFAULT NULL,
		` + "`source_channel`" + ` VARCHAR(255) DEFAULT NULL,
		` + "`source`" + ` VARCHAR(255) DEFAULT NULL,
		` + "`destination`" + ` VARCHAR(255) DEFAULT NULL,
		` + "`package`" + ` VARCHAR(255) DEFAULT NULL,
		` + "`charge`" + ` DOUBLE DEFAULT NULL,
		` + "`commission`" + ` DOUBLE DEFAULT NULL,
		` + "`external_reference_number`" + ` VARCHAR(255) DEFAULT NULL,
		` + "`is_async`" + ` TINYINT(1) DEFAULT NULL,
		` + "`corp_id`" + ` VARCHAR(255) DEFAULT NULL,
		` + "`status_id`" + ` INT(11) NOT NULL,
		` + "`extra_details_1`" + ` VARCHAR(255) DEFAULT NULL,
		` + "`extra_details_2`" + ` VARCHAR(255) DEFAULT NULL,
		` + "`extra_details_3`" + ` VARCHAR(255) DEFAULT NULL,
		` + "`client_response_code`" + ` VARCHAR(255) DEFAULT NULL,
		` + "`date_created`" + ` DATETIME DEFAULT CURRENT_TIMESTAMP,
		` + "`date_modified`" + ` DATETIME DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
		` + "`created_by`" + ` INT(11) DEFAULT NULL,
		` + "`modified_by`" + ` INT(11) DEFAULT NULL,
		` + "`active`" + ` INT(11) DEFAULT NULL,
		PRIMARY KEY (` + "`transaction_id`" + `),
		FOREIGN KEY (` + "`service_id`" + `) REFERENCES services(` + "`service_id`" + `) ON UPDATE CASCADE ON DELETE NO ACTION,
		FOREIGN KEY (` + "`request_id`" + `) REFERENCES user_requests(` + "`user_request_id`" + `) ON UPDATE CASCADE ON DELETE NO ACTION,
		FOREIGN KEY (` + "`transaction_customer_reference`" + `) REFERENCES customers(` + "`customer_id`" + `) ON UPDATE CASCADE ON DELETE SET NULL,
		FOREIGN KEY (` + "`status_id`" + `) REFERENCES status_codes(` + "`status_id`" + `) ON UPDATE CASCADE ON DELETE NO ACTION,
		FOREIGN KEY (` + "`created_by`" + `) REFERENCES auth_users(` + "`user_id`" + `) ON UPDATE CASCADE ON DELETE SET NULL,
		FOREIGN KEY (` + "`modified_by`" + `) REFERENCES auth_users(` + "`user_id`" + `) ON UPDATE CASCADE ON DELETE SET NULL
	)`)
}

// Reverse the migrations
func (m *UserTransactions_20260416_000000) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `user_transactions`")
}
