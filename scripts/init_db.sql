-- MySQL dump 10.13  Distrib 8.0.25, for Linux (x86_64)
--
-- Host: pancast.cgizr834pu3g.us-east-2.rds.amazonaws.com    Database: pancast
-- ------------------------------------------------------
-- Server version	8.0.20

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
SET @MYSQLDUMP_TEMP_LOG_BIN = @@SESSION.SQL_LOG_BIN;
SET @@SESSION.SQL_LOG_BIN= 0;

--
-- GTID state at the beginning of the backup 
--

SET @@GLOBAL.GTID_PURGED=/*!80000 '+'*/ '';

CREATE DATABASE pancast;
USE pancast;

--
-- Table structure for table `device`
--

DROP TABLE IF EXISTS `device`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `device` (
  `device_id` int NOT NULL,
  `secret_key` char(32) DEFAULT NULL,
  `clock_init` int DEFAULT '0',
  `clock_offset` int DEFAULT '0',
  PRIMARY KEY (`device_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `beacon`
--

DROP TABLE IF EXISTS `beacon`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `beacon` (
  `device_id` int NOT NULL,
  `location_id` char(8) DEFAULT 'DEADBEEF',
  PRIMARY KEY (`device_id`),
  KEY `idx_beacon_location_id` (`location_id`),
  CONSTRAINT `beacon_device_id` FOREIGN KEY (`device_id`) REFERENCES `device` (`device_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `dongle`
--

DROP TABLE IF EXISTS `dongle`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dongle` (
  `device_id` int NOT NULL,
  PRIMARY KEY (`device_id`),
  CONSTRAINT `dongle_device_id` FOREIGN KEY (`device_id`) REFERENCES `device` (`device_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `epi_entries`
--

DROP TABLE IF EXISTS `epi_entries`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `epi_entries` (
  `eph_id` binary(15) NOT NULL,
  `location_id` char(8) NOT NULL,
  `time_dongle` int DEFAULT '0',
  `time_beacon` int DEFAULT '0',
  `beacon_id` int NOT NULL,
  PRIMARY KEY (`eph_id`,`location_id`,`beacon_id`),
  KEY `epi_beacon_id_idx` (`location_id`),
  KEY `epi_beacon_id_idx1` (`beacon_id`),
  CONSTRAINT `epi_beacon_id` FOREIGN KEY (`beacon_id`) REFERENCES `beacon` (`device_id`),
  CONSTRAINT `epi_location_id` FOREIGN KEY (`location_id`) REFERENCES `beacon` (`location_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `risk_entries`
--

DROP TABLE IF EXISTS `risk_entries`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `risk_entries` (
  `eph_id` binary(15) NOT NULL,
  `location_id` char(8) NOT NULL,
  `time_dongle` int DEFAULT '0',
  `time_beacon` int DEFAULT '0',
  `beacon_id` int NOT NULL,
  PRIMARY KEY (`eph_id`,`location_id`,`beacon_id`),
  KEY `risk_beacon_id_idx` (`location_id`),
  KEY `risk_beacon_id_idx1` (`beacon_id`),
  CONSTRAINT `risk_beacon_id` FOREIGN KEY (`beacon_id`) REFERENCES `beacon` (`device_id`),
  CONSTRAINT `risk_location_id` FOREIGN KEY (`location_id`) REFERENCES `beacon` (`location_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

SET @@SESSION.SQL_LOG_BIN = @MYSQLDUMP_TEMP_LOG_BIN;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

