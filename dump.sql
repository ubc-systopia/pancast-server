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

--
-- Table structure for table `beacon`
--

CREATE DATABASE pancast;
USE pancast;

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
-- Dumping data for table `beacon`
--

LOCK TABLES `beacon` WRITE;
/*!40000 ALTER TABLE `beacon` DISABLE KEYS */;
INSERT INTO `beacon` VALUES (1,'LOC00001'),(2,'LOC00002'),(3,'LOC00003'),(4,'LOC00004'),(5,'LOC00005'),(6,'LOC00006'),(7,'LOC00007'),(8,'LOC00008'),(9,'LOC00009'),(10,'LOC00010'),(11,'LOC00011'),(12,'LOC00012'),(13,'LOC00013'),(14,'LOC00014'),(15,'LOC00015'),(16,'LOC00016'),(17,'LOC00017'),(18,'LOC00018'),(19,'LOC00019'),(20,'LOC00020');
/*!40000 ALTER TABLE `beacon` ENABLE KEYS */;
UNLOCK TABLES;

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
-- Dumping data for table `device`
--

LOCK TABLES `device` WRITE;
/*!40000 ALTER TABLE `device` DISABLE KEYS */;
INSERT INTO `device` VALUES (1,'462D4A614E645267556B587032733576',0,0),(2,'24432646294A404E635266556A586E32',0,0),(3,'763979244226452948404D635166546A',0,0),(4,'3373367638792F423F4528482B4D6251',0,0),(5,'6B58703273357638782F413F4428472B',0,0),(6,'5266556A586E32723575377821412544',0,0),(7,'404D635166546A576E5A723475377721',0,0),(8,'4528482B4D6251655468576D5A713474',0,0),(9,'2F413F4428472B4B6250655368566D59',0,0),(10,'753778214125442A472D4B6150645367',0,0),(11,'5A7234743777217A25432A462D4A614E',0,0),(12,'68576D5A7134743677397A2443264629',0,0),(13,'50655368566D59713374367639792442',0,0),(14,'2D4B6150645367566B59703373367638',0,0),(15,'432A462D4A614E645267556B58703273',0,0),(16,'397A24432646294A404E635266556A58',0,0),(17,'7336763979244226452948404D635166',0,0),(18,'59703373357638792F423F4528482B4D',0,0),(19,'67556B58703272357538782F413F4428',0,0),(20,'4E635266556A586E3272347537782141',0,0);
/*!40000 ALTER TABLE `device` ENABLE KEYS */;
UNLOCK TABLES;

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
-- Dumping data for table `dongle`
--

LOCK TABLES `dongle` WRITE;
/*!40000 ALTER TABLE `dongle` DISABLE KEYS */;
/*!40000 ALTER TABLE `dongle` ENABLE KEYS */;
UNLOCK TABLES;

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
-- Dumping data for table `epi_entries`
--

LOCK TABLES `epi_entries` WRITE;
/*!40000 ALTER TABLE `epi_entries` DISABLE KEYS */;
/*!40000 ALTER TABLE `epi_entries` ENABLE KEYS */;
UNLOCK TABLES;

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

--
-- Dumping data for table `risk_entries`
--

LOCK TABLES `risk_entries` WRITE;
/*!40000 ALTER TABLE `risk_entries` DISABLE KEYS */;
INSERT INTO `risk_entries` VALUES (_binary '00d17c102c46b0d','LOC00003',0,0,3),(_binary '04ad5f4414ab19f','LOC00008',0,0,8),(_binary '0c788e58b0b9b94','LOC00016',0,0,16),(_binary '174063dc93a8354','LOC00012',0,0,12),(_binary '19c4945cac1f5d6','LOC00009',0,0,9),(_binary '1dad8ea6111f7db','LOC00018',0,0,18),(_binary '4a10c16f3637b5f','LOC00007',0,0,7),(_binary '5f8465db3b720c8','LOC00011',0,0,11),(_binary '6a7f72f9571481b','LOC00006',0,0,6),(_binary '6b3e950b5210610','LOC00005',0,0,5),(_binary '6c8993217953c24','LOC00001',0,0,1),(_binary '735b5d91ec28891','LOC00010',0,0,10),(_binary '75ad9d470973990','LOC00020',0,0,20),(_binary 'a03774dd6e90df6','LOC00015',0,0,15),(_binary 'b36adc822470024','LOC00002',0,0,2),(_binary 'bfe27b6fb0cace0','LOC00019',0,0,19),(_binary 'ed11689a21d86b5','LOC00014',0,0,14),(_binary 'f319bafc4ec5f31','LOC00004',0,0,4),(_binary 'fb3d8063b0e6587','LOC00013',0,0,13),(_binary 'fdaa3482f063e17','LOC00017',0,0,17);
/*!40000 ALTER TABLE `risk_entries` ENABLE KEYS */;
UNLOCK TABLES;
SET @@SESSION.SQL_LOG_BIN = @MYSQLDUMP_TEMP_LOG_BIN;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-06-02 10:00:17

