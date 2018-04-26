CREATE TABLE IF NOT EXISTS `users` (
  `uuid` int(11) NOT NULL AUTO_INCREMENT,
  `displayName` varchar(256) COLLATE utf8_bin NOT NULL,
  `name` varchar(128) COLLATE utf8_bin NOT NULL,
  `password` char(128) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
COMMIT;