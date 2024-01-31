-- DDL (Data Definition Language)
DROP DATABASE IF EXISTS `product_test_db`;

CREATE DATABASE `product_test_db`;

USE `product_test_db`;

CREATE TABLE `products` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) DEFAULT NULL,
  `quantity` int DEFAULT NULL,
  `code_value` varchar(50) DEFAULT NULL,
  `is_published` varchar(50) DEFAULT NULL,
  `expiration` date DEFAULT NULL,
  `price` decimal(5,2) DEFAULT NULL,
  PRIMARY KEY (`id`),
)