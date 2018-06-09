CREATE TABLE IF NOT EXISTS `messages` (
  `id` int(10) unsigned NOT NULL,
  `body` varchar(200) DEFAULT NULL,
  `msg_box` varchar(20) NOT NULL,
  `address` varchar(20) NOT NULL,
  `synced` timestamp DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
