package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type UserInsTransactions_20260416_000001 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserInsTransactions_20260416_000001{}
	m.Created = "20260416_000001"

	migration.Register("UserInsTransactions_20260416_000001", m)
}

// Run the migrations
func (m *UserInsTransactions_20260416_000001) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE user_ins_transactions(
		` + "`ins_transaction_id`" + ` VARCHAR(255) NOT NULL,
		` + "`user_transaction_id`" + ` VARCHAR(255) NOT NULL,
		` + "`amount`" + ` DOUBLE DEFAULT NULL,
		` + "`data`" + ` VARCHAR(255) DEFAULT NULL,
		` + "`sender_account_number`" + ` VARCHAR(255) DEFAULT NULL,
		` + "`recipient_account_number`" + ` VARCHAR(255) DEFAULT NULL,
		` + "`service_id`" + ` INT(11) DEFAULT NULL,
		` + "`status`" + ` INT(11) DEFAULT NULL,
		` + "`request`" + ` VARCHAR(255) DEFAULT NULL,
		` + "`response`" + ` VARCHAR(255) DEFAULT NULL,
		` + "`date_created`" + ` DATETIME DEFAULT CURRENT_TIMESTAMP,
		` + "`date_modified`" + ` DATETIME DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
		` + "`created_by`" + ` INT(11) DEFAULT NULL,
		` + "`modified_by`" + ` INT(11) DEFAULT NULL,
		` + "`active`" + ` INT(11) DEFAULT NULL,
		PRIMARY KEY (` + "`ins_transaction_id`" + `),
		FOREIGN KEY (` + "`user_transaction_id`" + `) REFERENCES user_transactions(` + "`transaction_id`" + `) ON UPDATE CASCADE ON DELETE NO ACTION,
		FOREIGN KEY (` + "`service_id`" + `) REFERENCES services(` + "`service_id`" + `) ON UPDATE CASCADE ON DELETE SET NULL,
		FOREIGN KEY (` + "`status`" + `) REFERENCES status_codes(` + "`status_id`" + `) ON UPDATE CASCADE ON DELETE SET NULL
	)`)
}

// Reverse the migrations
func (m *UserInsTransactions_20260416_000001) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `user_ins_transactions`")
}
