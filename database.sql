-- 创建安卓的版本规则包
CREATE TABLE rulesForAndroid(
	`aid` INT,
	`platform` VARCHAR(20),
	`update_version_code` varchar(40),
	`max_update_version_code` varchar(40) not null,
	`min_update_version_code` varchar(40) not null,
	`max_os_api` int not null,
	`min_os_api` int not null,
	`cpu_arch` varchar(20) not null,
	`channel` varchar(40) not null,
	`download_url` varchar(255) not null,
	`md5` varchar(255) not null,
	`title` varchar(127) not null,
	`update_tips` varchar(255) not null,
	PRIMARY KEY(`aid`,`platform`,`update_version_code` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 创建iOS的版本规则包
CREATE TABLE rulesForiOS(
	`aid` INT,
	`platform` VARCHAR(20),
	`update_version_code` varchar(40),
	`max_update_version_code` varchar(40) not null,
	`min_update_version_code` varchar(40) not null,
	`cpu_arch` varchar(20) not null,
	`channel` varchar(40) not null,
	`download_url` varchar(255) not null,
	`md5` varchar(255) not null,
	`title` varchar(127) not null,
	`update_tips` varchar(255) not null,
	PRIMARY KEY(`aid`,`platform`,`update_version_code` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8;






