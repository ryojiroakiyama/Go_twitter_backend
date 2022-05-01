CREATE TABLE `account` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL UNIQUE,
  `password_hash` varchar(255) NOT NULL,
  `display_name` varchar(255),
  `avatar` text,
  `header` text,
  `note` text,
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);

CREATE TABLE `status` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `account_id` bigint(20) NOT NULL,
  `content` text NOT NULL,
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_account_id` (`account_id`),
  CONSTRAINT `fk_status_account_id` FOREIGN KEY (`account_id`) REFERENCES  `account` (`id`)
);

CREATE TABLE `relationship` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL,
  `follow_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE (`user_id`, `follow_id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_follow_id` (`follow_id`),
  CONSTRAINT `fk_relationship_user_id` FOREIGN KEY (`user_id`) REFERENCES `account` (`id`),
  CONSTRAINT `fk_relationship_follow_id` FOREIGN KEY (`follow_id`) REFERENCES `account` (`id`)
);