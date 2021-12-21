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
    PRIMARY KEY(`id`),
    UNIQUE KEY `uniq_name_username` (`name`, `username`),
    UNIQUE KEY `instanceID_UNIQUE` (`instanceID`),
    KEY `fk_policy_user_idx` (`username`),
    CONSTRAINT `fk_policy_user` FOREIGN KEY (`username`) REFERENCES `user` (`name`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=47 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `policy_audit`;
CREATE TABLE `policy_audit` (
    `id` bigint(20) unsigned NOT NULL,
    `instanceID` varchar(20) DEFAULT NULL,
    `name` varchar(45) NOT NULL,
    `username` varchar(255) NOT NULL,
    `policyShadow` longtext DEFAULT NULL,
    `extendShadow` longtext DEFAULT NULL,
    `createdAt` timestamp NOT NULL DEFAULT current_timestamp(),
    `updatedAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    `deletedAt` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
    PRIMARY KEY(`id`),
    KEY `fk_policy_user_idx` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `secret`;
CREATE TABLE `secret` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `instanceID` varchar(20) DEFAULT NULL,
    `name` varchar(45) NOT NULL,
    `username` varchar(255) NOT NULL,
    `secretID` varchar(36) NOT NULL,
    `secretKey` varchar(255) NOT NULL,
    `expires` int(64) unsigned NOT NULL DEFAULT 1534308590,
    `description` varchar(255) NOT NULL,
    `extendShadow` longtext DEFAULT NULL,
    `createdAt` timestamp NOT NULL DEFAULT current_timestamp(),
    `updatedAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_name_username` (`name`, `username`),
    UNIQUE KEY `instanceID_UNIQUE` (`instanceID`),
    KEY `fk_secret_user_idx` (`username`),
    CONSTRAINT `fk_secret_user` FOREIGN KEY (`username`) REFERENCES `user` (`name`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `instanceID` varchar(20) DEFAULT NULL,
  `name` varchar(45) NOT NULL,
  `nickname` varchar(30) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(256) NOT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `isAdmin` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT '1: administrator\\\\n0: non-administrator',
  `extendShadow` longtext DEFAULT NULL,
  `createdAt` timestamp NOT NULL DEFAULT current_timestamp(),
  `updatedAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`),
  UNIQUE KEY `instanceID_UNIQUE` (`instanceID`)
) ENGINE=InnoDB AUTO_INCREMENT=38 DEFAULT CHARSET=utf8;