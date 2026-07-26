-- Initial SQL migration for transaction_service models
-- Generated from: models/*.go
-- Dialect: MySQL 8+

SET NAMES utf8mb4;
SET time_zone = '+00:00';

CREATE TABLE IF NOT EXISTS status (
    status_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    status VARCHAR(128) NOT NULL,
    status_code VARCHAR(128) NOT NULL,
    date_created DATETIME NOT NULL,
    date_modified DATETIME NOT NULL,
    created_by INT NOT NULL,
    modified_by INT NOT NULL,
    active INT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS status_codes (
    status_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    status_code VARCHAR(50) NOT NULL,
    status_description VARCHAR(255) NOT NULL,
    active INT NOT NULL,
    date_created DATETIME NOT NULL,
    date_modified DATETIME NOT NULL,
    created_by INT NOT NULL,
    modified_by INT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS requests (
    request_id VARCHAR(36) PRIMARY KEY,
    api_request_id VARCHAR(255) NOT NULL,
    cust_id VARCHAR(255) NOT NULL,
    request VARCHAR(1000) NOT NULL,
    request_type VARCHAR(100) NOT NULL,
    request_status VARCHAR(255) NOT NULL,
    request_amount DOUBLE NOT NULL,
    request_response VARCHAR(1000) NOT NULL,
    callback_response VARCHAR(1000) NOT NULL,
    request_date DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_requests (
    user_request_id VARCHAR(36) PRIMARY KEY,
    api_request_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    request VARCHAR(1000) NOT NULL,
    request_type VARCHAR(100) NOT NULL,
    request_status VARCHAR(255) NOT NULL,
    request_amount DOUBLE NOT NULL,
    request_response VARCHAR(1000) NOT NULL,
    callback_response VARCHAR(1000) NOT NULL,
    request_date DATETIME NOT NULL,
    date_created DATETIME NOT NULL,
    date_modified DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS orders (
    order_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    customer_id VARCHAR(200) NULL,
    customer_name VARCHAR(255) NULL,
    customer_email VARCHAR(255) NULL,
    customer_phone VARCHAR(255) NULL,
    order_number VARCHAR(100) NOT NULL,
    quantity INT NOT NULL,
    cost FLOAT NOT NULL,
    order_desc VARCHAR(500) NOT NULL,
    order_location VARCHAR(255) NOT NULL,
    currency VARCHAR(100) NOT NULL,
    order_date DATETIME NOT NULL,
    order_end_date DATETIME NOT NULL,
    returned_date DATETIME NOT NULL,
    date_created DATETIME NOT NULL,
    date_modified DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    modified_by BIGINT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS transactions (
    request_id VARCHAR(36) PRIMARY KEY,
    order_id BIGINT NOT NULL,
    branch_id VARCHAR(255) NOT NULL,
    branch_name VARCHAR(255) NOT NULL,
    amount FLOAT NOT NULL,
    currency_id VARCHAR(255) NOT NULL,
    currency_symbol VARCHAR(255) NOT NULL,
    status_id BIGINT NOT NULL,
    date_created DATETIME NOT NULL,
    date_modified DATETIME NOT NULL,
    created_by INT NOT NULL,
    modified_by INT NOT NULL,
    active INT NOT NULL,
    service_id VARCHAR(255) NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    CONSTRAINT fk_transactions_order
        FOREIGN KEY (order_id)
        REFERENCES orders(order_id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_transactions_status
        FOREIGN KEY (status_id)
        REFERENCES status(status_id)
        ON DELETE RESTRICT,
    INDEX idx_transactions_order_id (order_id),
    INDEX idx_transactions_status_id (status_id),
    INDEX idx_transactions_service_id (service_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS bil_transactions (
    transaction_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    transaction_ref_number VARCHAR(255) NOT NULL UNIQUE,
    service_id VARCHAR(255) NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    biller_code VARCHAR(255) NOT NULL,
    request_id VARCHAR(36) NOT NULL,
    transaction_by VARCHAR(255) NOT NULL,
    amount DOUBLE NOT NULL,
    transacting_currency VARCHAR(255) NOT NULL,
    source_channel VARCHAR(255) NOT NULL,
    source VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    package VARCHAR(255) NOT NULL,
    charge DOUBLE NOT NULL,
    commission DOUBLE NOT NULL,
    external_reference_number VARCHAR(255) NOT NULL,
    is_async TINYINT(1) NOT NULL,
    corp_id VARCHAR(255) NOT NULL,
    status_id BIGINT NOT NULL,
    extra_details_1 VARCHAR(255) NOT NULL,
    extra_details_2 VARCHAR(255) NOT NULL,
    extra_details_3 VARCHAR(255) NOT NULL,
    client_response_code VARCHAR(255) NOT NULL,
    date_created DATETIME NOT NULL,
    date_modified DATETIME NOT NULL,
    created_by VARCHAR(255) NOT NULL,
    modified_by VARCHAR(255) NOT NULL,
    active INT NOT NULL,
    CONSTRAINT fk_bil_transactions_request
        FOREIGN KEY (request_id)
        REFERENCES requests(request_id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_bil_transactions_status
        FOREIGN KEY (status_id)
        REFERENCES status_codes(status_id)
        ON DELETE RESTRICT,
    INDEX idx_bil_transactions_service_id (service_id),
    INDEX idx_bil_transactions_request_id (request_id),
    INDEX idx_bil_transactions_status_id (status_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS bil_ins_transactions (
    bil_ins_transaction_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    bil_transaction_id BIGINT NOT NULL,
    amount DOUBLE NOT NULL,
    biller_id VARCHAR(255) NOT NULL,
    biller_code VARCHAR(255) NOT NULL,
    biller_name VARCHAR(255) NOT NULL,
    sender_account_number VARCHAR(255) NOT NULL,
    recipient_account_number VARCHAR(255) NOT NULL,
    network VARCHAR(150) NOT NULL,
    status BIGINT NOT NULL,
    request VARCHAR(255) NOT NULL,
    response VARCHAR(255) NOT NULL,
    date_created DATETIME NOT NULL,
    date_modified DATETIME NOT NULL,
    created_by INT NOT NULL,
    modified_by INT NOT NULL,
    active INT NOT NULL,
    CONSTRAINT fk_bil_ins_transactions_bil_transaction
        FOREIGN KEY (bil_transaction_id)
        REFERENCES bil_transactions(transaction_id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_bil_ins_transactions_status
        FOREIGN KEY (status)
        REFERENCES status_codes(status_id)
        ON DELETE RESTRICT,
    INDEX idx_bil_ins_transactions_bil_transaction_id (bil_transaction_id),
    INDEX idx_bil_ins_transactions_biller_id (biller_id),
    INDEX idx_bil_ins_transactions_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS order_items (
    order_item_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT NOT NULL,
    item_id VARCHAR(255) NOT NULL,
    item_name VARCHAR(255) NOT NULL,
    unit_price FLOAT NOT NULL,
    quantity INT NOT NULL,
    item_status BIGINT NOT NULL,
    total_price FLOAT NOT NULL,
    order_date DATETIME NOT NULL,
    date_created DATETIME NOT NULL,
    date_modified DATETIME NOT NULL,
    created_by BIGINT NOT NULL,
    modified_by BIGINT NOT NULL,
    comment LONGTEXT NOT NULL,
    CONSTRAINT fk_order_items_order
        FOREIGN KEY (order_id)
        REFERENCES orders(order_id)
        ON DELETE CASCADE,
    CONSTRAINT fk_order_items_status
        FOREIGN KEY (item_status)
        REFERENCES status(status_id)
        ON DELETE RESTRICT,
    INDEX idx_order_items_order_id (order_id),
    INDEX idx_order_items_item_status (item_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS transaction_details (
    transaction_detail_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    transaction_id VARCHAR(36) NOT NULL,
    amount FLOAT NOT NULL,
    comment VARCHAR(255) NULL,
    sender_account_number VARCHAR(255) NOT NULL,
    recipient_account_number VARCHAR(255) NOT NULL,
    status_code VARCHAR(20) NOT NULL,
    status_message VARCHAR(255) NOT NULL,
    sender_id BIGINT NOT NULL,
    transaction_type VARCHAR(255) NOT NULL,
    date_created DATETIME NOT NULL,
    date_modified DATETIME NOT NULL,
    created_by INT NOT NULL,
    modified_by INT NOT NULL,
    active INT NOT NULL,
    CONSTRAINT fk_transaction_details_transaction
        FOREIGN KEY (transaction_id)
        REFERENCES transactions(request_id)
        ON DELETE CASCADE,
    INDEX idx_transaction_details_transaction_id (transaction_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_transactions (
    transaction_id VARCHAR(255) NOT NULL,
    service_id VARCHAR(255) NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    request_id VARCHAR(36) NOT NULL,
    transaction_customer_reference VARCHAR(255) NULL,
    amount DOUBLE NOT NULL,
    transacting_currency VARCHAR(255) NOT NULL,
    source_channel VARCHAR(255) NOT NULL,
    source VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    package VARCHAR(255) NOT NULL,
    charge DOUBLE NOT NULL,
    commission DOUBLE NOT NULL,
    reference VARCHAR(255) NULL,
    external_reference_number VARCHAR(255) NOT NULL,
    is_async TINYINT(1) NOT NULL,
    status_id BIGINT NOT NULL,
    corp_id VARCHAR(255) NOT NULL,
    extra_details_1 VARCHAR(255) NOT NULL,
    extra_details_2 VARCHAR(255) NOT NULL,
    extra_details_3 VARCHAR(255) NOT NULL,
    client_response_code VARCHAR(255) NOT NULL,
    date_created DATETIME NOT NULL,
    date_modified DATETIME NOT NULL,
    created_by VARCHAR(255) NOT NULL,
    modified_by VARCHAR(255) NOT NULL,
    active INT NOT NULL,
    CONSTRAINT fk_user_transactions_request
        FOREIGN KEY (request_id)
        REFERENCES user_requests(user_request_id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_user_transactions_status
        FOREIGN KEY (status_id)
        REFERENCES status_codes(status_id)
        ON DELETE RESTRICT,
    UNIQUE KEY uq_user_transactions_transaction_id (transaction_id),
    INDEX idx_user_transactions_service_id (service_id),
    INDEX idx_user_transactions_request_id (request_id),
    INDEX idx_user_transactions_status_id (status_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_ins_transactions (
    ins_transaction_id VARCHAR(255) NOT NULL,
    user_transaction_id VARCHAR(255) NOT NULL,
    amount DOUBLE NOT NULL,
    data VARCHAR(255) NULL,
    sender_account_number VARCHAR(255) NOT NULL,
    recipient_account_number VARCHAR(255) NOT NULL,
    service_id VARCHAR(255) NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    status BIGINT NOT NULL,
    request VARCHAR(255) NULL,
    response VARCHAR(255) NULL,
    date_created DATETIME NOT NULL,
    date_modified DATETIME NOT NULL,
    created_by INT NOT NULL,
    modified_by INT NOT NULL,
    active INT NOT NULL,
    CONSTRAINT fk_user_ins_transactions_status
        FOREIGN KEY (status)
        REFERENCES status_codes(status_id)
        ON DELETE RESTRICT,
    UNIQUE KEY uq_user_ins_transactions_ins_transaction_id (ins_transaction_id),
    INDEX idx_user_ins_transactions_service_id (service_id),
    INDEX idx_user_ins_transactions_status (status),
    INDEX idx_user_ins_transactions_user_transaction_id (user_transaction_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Incremental schema updates for model changes
-- Convert created_by/modified_by from numeric types to string types.
ALTER TABLE orders
    MODIFY COLUMN created_by VARCHAR(255) NOT NULL,
    MODIFY COLUMN modified_by VARCHAR(255) NOT NULL;

ALTER TABLE transactions
    MODIFY COLUMN created_by VARCHAR(255) NOT NULL,
    MODIFY COLUMN modified_by VARCHAR(255) NOT NULL;

ALTER TABLE order_items
    MODIFY COLUMN created_by VARCHAR(255) NOT NULL,
    MODIFY COLUMN modified_by VARCHAR(255) NOT NULL;

ALTER TABLE transaction_details
    MODIFY COLUMN created_by VARCHAR(255) NOT NULL,
    MODIFY COLUMN modified_by VARCHAR(255) NOT NULL;
