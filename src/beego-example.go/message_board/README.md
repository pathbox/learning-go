CREATE TABLE `leave_message` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL DEFAULT '0' COMMENT '来自用户表user的id',
  `content` text NOT NULL COMMENT '留言内容',
  `status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '是否展示 0 否 1是',
  `create_at` datetime NOT NULL COMMENT '创建时间',
  `update_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`) USING BTREE
) ENGINE=MyISAM AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='留言表';