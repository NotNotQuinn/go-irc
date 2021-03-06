CREATE TABLE IF NOT EXISTS `wb_core`.`user` (
	`ID` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	`Name` VARCHAR(50) NOT NULL COLLATE 'utf8_general_ci',
	`Twitch_ID` INT(16) UNSIGNED NOT NULL,
	`First_Seen` DATETIME NOT NULL DEFAULT current_timestamp() COMMENT 'When user was first added to the database',
	PRIMARY KEY (`ID`) USING BTREE,
	UNIQUE INDEX `Name` (`Name`) USING BTREE,
	UNIQUE INDEX `Twitch_ID` (`Twitch_ID`) USING BTREE
)
COMMENT='Users stored in the system'
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;
