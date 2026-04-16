package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type UserRequests_20260415_234607 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserRequests_20260415_234607{}
	m.Created = "20260415_234607"

	migration.Register("UserRequests_20260415_234607", m)
}

// Run the migrations
func (m *UserRequests_20260415_234607) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE user_requests(`user_request_id` int(11) NOT NULL AUTO_INCREMENT,`api_request_id` varchar(255) NOT NULL,`user_id` int(11) DEFAULT NULL,`request` longtext  DEFAULT NULL,`request_type` varchar(100) DEFAULT NULL,`request_status` varchar(100) DEFAULT NULL,`request_amount` float DEFAULT 0,`request_response` longtext  DEFAULT NULL,`callback_response` longtext  DEFAULT NULL,`request_date` datetime DEFAULT CURRENT_TIMESTAMP,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,PRIMARY KEY (`user_request_id`), FOREIGN KEY (user_id) REFERENCES auth_users(user_id) ON UPDATE CASCADE ON DELETE NO ACTION)")
}

// Reverse the migrations
func (m *UserRequests_20260415_234607) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `user_requests`")
}
