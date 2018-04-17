CREATE TABLE IF NOT EXISTS `summary` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `snapshotdate` date NOT NULL,
  `cash` decimal(10,2),
  `loads` decimal(10,2),
  `smartmoney` decimal(10,2),
  `codes` decimal(10,2),
  `total` decimal(10,2),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci AUTO_INCREMENT=1 ;