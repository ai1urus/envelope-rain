use envelope;

truncate users;

DROP PROCEDURE IF EXISTS insert_user;

DELIMITER $$
CREATE PROCEDURE insert_user(IN start_num INT(10), IN max_num INT(10))
BEGIN
  DECLARE i INT DEFAULT 0;
  SET autocommit=0;
  REPEAT
  SET i=i+1;
  INSERT INTO users (uid, amount, cur_count) VALUES (start_num+i, 0, 0);
  UNTIL i=max_num
  END REPEAT;
  COMMIT;
END $$
DELIMITER ;

CALL insert_user(0, 100000);