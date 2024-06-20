package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Services_20240620_094435 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Services_20240620_094435{}
	m.Created = "20240620_094435"

	migration.Register("Services_20240620_094435", m)
}

// Run the migrations
func (m *Services_20240620_094435) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE service(`service_id` int(11) NOT NULL AUTO_INCREMENT,`service_name` varchar(255) NOT NULL,`service_description` varchar(500) DEFAULT NULL,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT NULL,`modified_by` int(11) DEFAULT NULL,`active` int(11) DEFAULT 1,PRIMARY KEY (`service_id`))")
}

// Reverse the migrations
func (m *Services_20240620_094435) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `services`")
}
