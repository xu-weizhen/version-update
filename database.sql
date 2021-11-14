create database update_version character set utf8mb4;

create user 'uv_query_usr'@'%' identified by '123456';
grant select on update_version.* to 'uv_query_usr'@'%';

create user 'uv_admin_usr'@'%' identified by '123456';
grant all on update_version.* to 'uv_admin_usr'@'%';

use update_version;

-- 创建安卓的版本规则包
CREATE TABLE rulesForAndroid(
	`id` INT primary key auto_increment, 
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
	`download_count` BIGINT not null default 0,
	`status` int not null default 0,
	-- PRIMARY KEY(`aid`,`platform`,`update_version_code`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
create index rulesForAndroid_index on rulesForAndroid(`aid`,`platform`,`update_version_code`)

-- 创建iOS的版本规则包
CREATE TABLE rulesForiOS(
	`id` INT primary key auto_increment, 
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
	`download_count` BIGINT not null default 0,
	`status` int not null default 0,
	-- PRIMARY KEY(`aid`,`platform`,`update_version_code` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
create index rulesForiOS_index on rulesForiOS(`aid`,`platform`,`update_version_code`)

-- 创建设备白名单列表
CREATE TABLE device_id(
	`aid` INT,
	`platform` VARCHAR(20),
	`update_version_code` varchar(40),
	`device_id_list` MEDIUMTEXT NOT NULL,
	PRIMARY KEY(`aid`,`platform`,`update_version_code` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8;