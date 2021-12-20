CREATE DATABASE IF NOT EXISTS `iam`;
USE `iam`;

DROP TABLE IF EXISTS `policy`;
CREATE TABLE `policy` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `instanceID` varchar(20) DEFAULT NULL,
    `name` varchar(45) NOT NULL,
    `username` varchar(255) NOT NULL,
    `policyShadow` longtext DEFAULT NULL,
    `extendShadow` longtext DEFAULT NULL,
    `createdAt` timestamp NOT NULL DEFAULT current_timestamp(),
    `updatedAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    PRIMARY KEY(`id`);
    UNIQUE KEY `uniq_name_username` (`name`, `username`),
    UNIQUE KEY `instanceID_UNIQUE` (`instanceID`),
    KEY `fk_policy_user_idx` (`username`)
);