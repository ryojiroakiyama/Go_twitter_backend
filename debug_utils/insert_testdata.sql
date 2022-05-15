INSERT INTO `account` (`username`, `password_hash`) VALUES
('benben', '$2a$10$T3C9WgYroD2SWAQegbB0qOzVC4XbqnWHHd9srL5DQ2ixbSj.Y4MDO');

SET @id1 = LAST_INSERT_ID();

INSERT INTO `status` (`account_id`, `content`) VALUES
(@id1, 'status of account1');

INSERT INTO `account` (`username`, `password_hash`) VALUES
('jonjon', '$2a$10$T3C9WgYroD2SWAQegbB0qOzVC4XbqnWHHd9srL5DQ2ixbSj.Y4MDO');

SET @id2 = LAST_INSERT_ID();

INSERT INTO `status` (`account_id`, `content`) VALUES
(@id2, 'status of account2');

INSERT INTO `account` (`username`, `password_hash`) VALUES
('sonson', '$2a$10$T3C9WgYroD2SWAQegbB0qOzVC4XbqnWHHd9srL5DQ2ixbSj.Y4MDO');

SET @id3 = LAST_INSERT_ID();

INSERT INTO `status` (`account_id`, `content`) VALUES
(@id3, 'status of account3');

INSERT INTO `relationship` (`user_id`, `follow_id`) VALUES
(@id1, @id2);

INSERT INTO `relationship` (`user_id`, `follow_id`) VALUES
(@id1, @id3);

INSERT INTO `relationship` (`user_id`, `follow_id`) VALUES
(@id2, @id3);
