CREATE TABLE IF NOT EXISTS `loads` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `trans_datetime` datetime NOT NULL,
  `trans_code` varchar(30) NOT NULL,
  `mobile_number` varchar(12) DEFAULT NULL,
  `amount` decimal(10,2) NOT NULL,
  `fee` boolean DEFAULT 0,
  `details` varchar(200) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci AUTO_INCREMENT=1 ;
