-- CREATE TABLE "room" ---------------------------
CREATE TABLE `room` ( 
	`id` BigInt( 20 ) UNSIGNED AUTO_INCREMENT NOT NULL,
	`name` VarChar( 50 ) NOT NULL,
	`host_user` VarChar( 50 ) NOT NULL,
	`participant` JSON NULL,
	`capacity` BigInt( 20 ) UNSIGNED NOT NULL DEFAULT 5,
	`created_at` DateTime NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` DateTime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_name` UNIQUE( `name` ))
ENGINE = InnoDB;
----------------------------------------------------------------