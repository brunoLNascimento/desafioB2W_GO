CREATE DATABASE desafioGo;

CREATE TABLE `planets` (
	`ID` INT(11) NOT NULL AUTO_INCREMENT,
	`PLANET_NAME` VARCHAR(80) NOT NULL,
	`PLANET_TERRAIN` VARCHAR(80) NOT NULL,
	`PLANET_FILMS` INT(11) NOT NULL,
	PRIMARY KEY (`ID`)
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
AUTO_INCREMENT=1
;
