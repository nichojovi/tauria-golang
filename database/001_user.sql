-- CREATE TABLE "user" ---------------------------
CREATE TABLE `user` ( 
	`id` BigInt( 20 ) UNSIGNED AUTO_INCREMENT NOT NULL,
	`user_name` VarChar( 50 ) NOT NULL,
	`password` VarChar( 50 ) NOT NULL,
	`full_name` VarChar( 50 ) NOT NULL,
	`email` VarChar( 50 ) NOT NULL,
	`phone` VarChar( 50 ) NOT NULL,
	`created_at` DateTime NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` DateTime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_user_name` UNIQUE( `user_name` ))
ENGINE = InnoDB;
----------------------------------------------------------------