CREATE TABLE `t_flounder_callcenter_call_record` (
  `record_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `callid` char(32) NOT NULL DEFAULT '' COMMENT 'CallID',
  `remote_url` varchar(512) NOT NULL DEFAULT '' COMMENT 'RecordFile',
  `call_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `call_date` char(10) NOT NULL DEFAULT '' COMMENT '日期',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `call_sheet_id` char(40) NOT NULL DEFAULT '' COMMENT '通话记录编号',
  PRIMARY KEY (`record_id`),
  KEY `callid` (`callid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='采集数据';
