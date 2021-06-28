
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

--
-- Initialize database if it does not already exist
--

CREATE DATABASE IF NOT EXISTS pancast;
USE pancast;

--
-- Table structure for table `beacon`
--

DROP TABLE IF EXISTS `beacon`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `beacon` (
  `device_id` int NOT NULL,
  `location_id` bigint DEFAULT NULL,
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
INSERT INTO `beacon` VALUES (1,1),(2,2),(3,3),(4,4),(5,5),(6,6),(7,7),(8,8),(9,9),(10,10),(11,11),(12,12),(13,13),(14,14),(15,15),(16,16),(17,17),(18,18),(19,19),(20,20);
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
  `secret_key` binary(32) DEFAULT NULL,
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
INSERT INTO `device` VALUES (0,_binary '�/��\�Fֱ�h\�l|*\�{	�Dc)\�/��M',27070951,0),(1,_binary '462D4A614E645267556B587032733576',0,0),(2,_binary '24432646294A404E635266556A586E32',0,0),(3,_binary '763979244226452948404D635166546A',0,0),(4,_binary '3373367638792F423F4528482B4D6251',0,0),(5,_binary '6B58703273357638782F413F4428472B',0,0),(6,_binary '5266556A586E32723575377821412544',0,0),(7,_binary '404D635166546A576E5A723475377721',0,0),(8,_binary '4528482B4D6251655468576D5A713474',0,0),(9,_binary '2F413F4428472B4B6250655368566D59',0,0),(10,_binary '753778214125442A472D4B6150645367',0,0),(11,_binary '5A7234743777217A25432A462D4A614E',0,0),(12,_binary '68576D5A7134743677397A2443264629',0,0),(13,_binary '50655368566D59713374367639792442',0,0),(14,_binary '2D4B6150645367566B59703373367638',0,0),(15,_binary '432A462D4A614E645267556B58703273',0,0),(16,_binary '397A24432646294A404E635266556A58',0,0),(17,_binary '7336763979244226452948404D635166',0,0),(18,_binary '59703373357638792F423F4528482B4D',0,0),(19,_binary '67556B58703272357538782F413F4428',0,0),(20,_binary '4E635266556A586E3272347537782141',0,0),(21,_binary 'e�\n�ä\�K���\�>�i�6�\\\�-.\�j\�\�ʔ����',27070974,0),(22,_binary 'WI�\�N��x&87;\�?Zd\\BMR\�B\�Z�݅�a`\�\r',27070974,0),(23,_binary '\�z\�)qhe\�W��\��\�$\�^(_\��bn��\�5�\�=',27070973,0),(24,_binary '��/G?3��}uN��*\�o\"�\�\�#z\�\�1',27070974,0),(25,_binary 'I�O�\"J��@��2|0\���I\�&�`@>nRY�N',27079424,0),(26,_binary 'M2��켧\�\�\�\�Ԫ�a�\�jY�2�+\��',27080575,0);
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
INSERT INTO `dongle` VALUES (0),(21),(22),(23),(24),(25),(26);
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
  `location_id` bigint DEFAULT NULL,
  `time_dongle` int DEFAULT '0',
  `time_beacon` int DEFAULT '0',
  `beacon_id` int NOT NULL,
  PRIMARY KEY (`eph_id`,`beacon_id`),
  KEY `epi_beacon_id_idx1` (`beacon_id`),
  KEY `epi_location_id_idx` (`location_id`),
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
  `location_id` bigint DEFAULT '0',
  `time_beacon` int DEFAULT '0',
  `time_dongle` int DEFAULT '0',
  `beacon_id` int NOT NULL,
  PRIMARY KEY (`eph_id`,`beacon_id`),
  KEY `risk_beacon_id_idx1` (`beacon_id`),
  KEY `risk_location_id_idx` (`location_id`),
  CONSTRAINT `risk_beacon_id` FOREIGN KEY (`beacon_id`) REFERENCES `beacon` (`device_id`),
  CONSTRAINT `risk_location_id` FOREIGN KEY (`location_id`) REFERENCES `beacon` (`location_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `risk_entries`
--

LOCK TABLES `risk_entries` WRITE;
/*!40000 ALTER TABLE `risk_entries` DISABLE KEYS */;
INSERT INTO `risk_entries` VALUES (_binary '00d17c102c46b0d',3,50000000,0,3),(_binary '04ad5f4414ab19f',8,50000000,0,8),(_binary '0c788e58b0b9b94',16,50000000,0,16),(_binary '174063dc93a8354',12,50000000,0,12),(_binary '19c4945cac1f5d6',9,50000000,0,9),(_binary '1dad8ea6111f7db',18,50000000,0,18),(_binary '4a10c16f3637b5f',7,50000000,0,7),(_binary '5f8465db3b720c8',11,50000000,0,11),(_binary '6a7f72f9571481b',6,50000000,0,6),(_binary '6b3e950b5210610',5,50000000,0,5),(_binary '6c8993217953c24',1,50000000,0,1),(_binary '735b5d91ec28891',10,50000000,0,10),(_binary '75ad9d470973990',20,50000000,0,20),(_binary 'a03774dd6e90df6',15,50000000,0,15),(_binary 'b36adc822470024',2,50000000,0,2),(_binary 'bfe27b6fb0cace0',19,50000000,0,19),(_binary 'deadbeefdeadbee',1,50000000,0,1),(_binary 'ed11689a21d86b5',14,50000000,0,14),(_binary 'f319bafc4ec5f31',4,50000000,0,4),(_binary 'fb3d8063b0e6587',13,50000000,0,13),(_binary 'fdaa3482f063e17',17,50000000,0,17);
/*!40000 ALTER TABLE `risk_entries` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
