/*
 Navicat PostgreSQL Data Transfer

 Source Server         : localhost
 Source Server Type    : PostgreSQL
 Source Server Version : 90618
 Source Host           : localhost:5432
 Source Catalog        : digger
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 90618
 File Encoding         : 65001

 Date: 07/09/2020 21:21:18
*/


-- ----------------------------
-- Sequence structure for seq_config_snapshot
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."seq_config_snapshot";
CREATE SEQUENCE "public"."seq_config_snapshot" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for seq_field
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."seq_field";
CREATE SEQUENCE "public"."seq_field" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for seq_plugin
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."seq_plugin";
CREATE SEQUENCE "public"."seq_plugin" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for seq_project
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."seq_project";
CREATE SEQUENCE "public"."seq_project" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for seq_result
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."seq_result";
CREATE SEQUENCE "public"."seq_result" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for seq_schedule_queue
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."seq_schedule_queue";
CREATE SEQUENCE "public"."seq_schedule_queue" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for seq_stage
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."seq_stage";
CREATE SEQUENCE "public"."seq_stage" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for seq_statistic
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."seq_statistic";
CREATE SEQUENCE "public"."seq_statistic" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for seq_task
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."seq_task";
CREATE SEQUENCE "public"."seq_task" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for t_config
-- ----------------------------
DROP TABLE IF EXISTS "public"."t_config";
CREATE TABLE "public"."t_config" (
  "key" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "value" varchar(500) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."t_config"."key" IS 'config_key';
COMMENT ON TABLE "public"."t_config" IS '配置表';

-- ----------------------------
-- Table structure for t_config_snapshot
-- ----------------------------
DROP TABLE IF EXISTS "public"."t_config_snapshot";
CREATE TABLE "public"."t_config_snapshot" (
  "id" int4 NOT NULL DEFAULT nextval('seq_config_snapshot'::regclass),
  "project_id" int4 NOT NULL,
  "config" json NOT NULL,
  "create_time" timestamp(6) NOT NULL DEFAULT now()
)
;
COMMENT ON COLUMN "public"."t_config_snapshot"."project_id" IS '所属项目';
COMMENT ON COLUMN "public"."t_config_snapshot"."config" IS '任务启动时刻的配置';
COMMENT ON COLUMN "public"."t_config_snapshot"."create_time" IS '创建时间';
COMMENT ON TABLE "public"."t_config_snapshot" IS '配置快照';

-- ----------------------------
-- Table structure for t_field
-- ----------------------------
DROP TABLE IF EXISTS "public"."t_field";
CREATE TABLE "public"."t_field" (
  "id" int4 NOT NULL DEFAULT nextval('seq_field'::regclass),
  "project_id" int4 NOT NULL,
  "stage_id" int4 NOT NULL,
  "is_array" bool NOT NULL,
  "is_html" bool NOT NULL,
  "css" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "attr" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "remark" varchar(255) COLLATE "pg_catalog"."default",
  "next_stage" varchar(255) COLLATE "pg_catalog"."default",
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "plugin" varchar(50) COLLATE "pg_catalog"."default",
  "xpath" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."t_field"."id" IS 'ID';
COMMENT ON COLUMN "public"."t_field"."project_id" IS '所属项目';
COMMENT ON COLUMN "public"."t_field"."stage_id" IS '所属阶段';
COMMENT ON COLUMN "public"."t_field"."is_array" IS '是否为数组';
COMMENT ON COLUMN "public"."t_field"."is_html" IS '是否为html';
COMMENT ON COLUMN "public"."t_field"."css" IS 'css选择器';
COMMENT ON COLUMN "public"."t_field"."attr" IS '字段属性';
COMMENT ON COLUMN "public"."t_field"."remark" IS '字段备注';
COMMENT ON COLUMN "public"."t_field"."next_stage" IS '下一阶段';
COMMENT ON COLUMN "public"."t_field"."name" IS '字段名';
COMMENT ON COLUMN "public"."t_field"."plugin" IS '插件';
COMMENT ON COLUMN "public"."t_field"."xpath" IS 'xpath选择器';
COMMENT ON TABLE "public"."t_field" IS '字段表';

-- ----------------------------
-- Table structure for t_plugin
-- ----------------------------
DROP TABLE IF EXISTS "public"."t_plugin";
CREATE TABLE "public"."t_plugin" (
  "id" int4 NOT NULL DEFAULT nextval('seq_plugin'::regclass),
  "name" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "script" text COLLATE "pg_catalog"."default" NOT NULL,
  "project_id" int4 NOT NULL
)
;
COMMENT ON COLUMN "public"."t_plugin"."name" IS '插件名称';
COMMENT ON COLUMN "public"."t_plugin"."script" IS '插件脚本内容';
COMMENT ON COLUMN "public"."t_plugin"."project_id" IS '所属项目';
COMMENT ON TABLE "public"."t_plugin" IS '插件表';

-- ----------------------------
-- Table structure for t_project
-- ----------------------------
DROP TABLE IF EXISTS "public"."t_project";
CREATE TABLE "public"."t_project" (
  "id" int4 NOT NULL DEFAULT nextval('seq_project'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "display_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "remark" varchar(255) COLLATE "pg_catalog"."default",
  "settings" json,
  "create_time" timestamp(6) NOT NULL DEFAULT now(),
  "update_time" timestamp(6) NOT NULL DEFAULT now(),
  "tags" varchar(255) COLLATE "pg_catalog"."default",
  "start_urls" json NOT NULL,
  "start_stage" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "headers" json,
  "node_affinity" json,
  "cron" varchar(255) COLLATE "pg_catalog"."default",
  "enable_cron" bool NOT NULL DEFAULT false
)
;
COMMENT ON COLUMN "public"."t_project"."name" IS '项目名称（唯一）';
COMMENT ON COLUMN "public"."t_project"."display_name" IS '项目显示名称';
COMMENT ON COLUMN "public"."t_project"."remark" IS '备注';
COMMENT ON COLUMN "public"."t_project"."settings" IS '配置';
COMMENT ON COLUMN "public"."t_project"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."t_project"."update_time" IS '更新时间';
COMMENT ON COLUMN "public"."t_project"."tags" IS '标签';
COMMENT ON COLUMN "public"."t_project"."start_urls" IS '开始地址';
COMMENT ON COLUMN "public"."t_project"."start_stage" IS '开始阶段';
COMMENT ON COLUMN "public"."t_project"."headers" IS '自定义headers';
COMMENT ON COLUMN "public"."t_project"."node_affinity" IS '节点亲和标签配置';
COMMENT ON COLUMN "public"."t_project"."cron" IS '定时任务配置，cron表达式';
COMMENT ON COLUMN "public"."t_project"."enable_cron" IS '是否启用定时任务';
COMMENT ON TABLE "public"."t_project" IS '项目表';

-- ----------------------------
-- Table structure for t_queue
-- ----------------------------
DROP TABLE IF EXISTS "public"."t_queue";
CREATE TABLE "public"."t_queue" (
  "id" int4 NOT NULL DEFAULT nextval('seq_schedule_queue'::regclass),
  "task_id" int4 NOT NULL,
  "stage_name" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "url" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "status" int2 NOT NULL DEFAULT 0,
  "middle_data" json,
  "expire" int8 NOT NULL DEFAULT 0
)
;
COMMENT ON COLUMN "public"."t_queue"."task_id" IS '所属任务id';
COMMENT ON COLUMN "public"."t_queue"."stage_name" IS '阶段stage名称';
COMMENT ON COLUMN "public"."t_queue"."url" IS '爬取链接地址';
COMMENT ON COLUMN "public"."t_queue"."status" IS '状态（0：未处理，1：已处理，2：异常）';
COMMENT ON COLUMN "public"."t_queue"."middle_data" IS '中间数据';
COMMENT ON COLUMN "public"."t_queue"."expire" IS '过期时间';
COMMENT ON TABLE "public"."t_queue" IS '调度任务表';

-- ----------------------------
-- Table structure for t_result
-- ----------------------------
DROP TABLE IF EXISTS "public"."t_result";
CREATE TABLE "public"."t_result" (
  "id" int4 NOT NULL DEFAULT nextval('seq_result'::regclass),
  "task_id" int4 NOT NULL,
  "result" json
)
;
COMMENT ON COLUMN "public"."t_result"."task_id" IS '所属任务id';
COMMENT ON COLUMN "public"."t_result"."result" IS '结果json文件';
COMMENT ON TABLE "public"."t_result" IS '爬虫结果表';

-- ----------------------------
-- Table structure for t_stage
-- ----------------------------
DROP TABLE IF EXISTS "public"."t_stage";
CREATE TABLE "public"."t_stage" (
  "id" int4 NOT NULL DEFAULT nextval('seq_stage'::regclass),
  "project_id" int4 NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "is_list" bool NOT NULL,
  "list_css" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "page_css" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "page_attr" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "plugins" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "is_unique" bool NOT NULL DEFAULT false,
  "page_xpath" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "list_xpath" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."t_stage"."project_id" IS '所属项目id';
COMMENT ON COLUMN "public"."t_stage"."name" IS '阶段名称';
COMMENT ON COLUMN "public"."t_stage"."is_list" IS '是否是list类型';
COMMENT ON COLUMN "public"."t_stage"."list_css" IS 'list类型的css选择器';
COMMENT ON COLUMN "public"."t_stage"."page_css" IS '分页css选择器';
COMMENT ON COLUMN "public"."t_stage"."page_attr" IS '分页css选择器属性';
COMMENT ON COLUMN "public"."t_stage"."plugins" IS '插件';
COMMENT ON COLUMN "public"."t_stage"."is_unique" IS '该阶段的链接是否全局唯一';
COMMENT ON COLUMN "public"."t_stage"."page_xpath" IS '分页xpath选择器';
COMMENT ON COLUMN "public"."t_stage"."list_xpath" IS '列表的xpath选择器';
COMMENT ON TABLE "public"."t_stage" IS '阶段表';

-- ----------------------------
-- Table structure for t_statistic
-- ----------------------------
DROP TABLE IF EXISTS "public"."t_statistic";
CREATE TABLE "public"."t_statistic" (
  "id" int4 NOT NULL DEFAULT nextval('seq_statistic'::regclass),
  "data" json,
  "create_time" timestamp(6) NOT NULL DEFAULT now()
)
;
COMMENT ON COLUMN "public"."t_statistic"."data" IS '数据';
COMMENT ON COLUMN "public"."t_statistic"."create_time" IS '创建时间';
COMMENT ON TABLE "public"."t_statistic" IS '统计表';

-- ----------------------------
-- Table structure for t_task
-- ----------------------------
DROP TABLE IF EXISTS "public"."t_task";
CREATE TABLE "public"."t_task" (
  "id" int4 NOT NULL DEFAULT nextval('seq_task'::regclass),
  "project_id" int4 NOT NULL,
  "config_snapshot_id" int4 NOT NULL,
  "status" int4 NOT NULL DEFAULT 0,
  "result_count" int4 NOT NULL DEFAULT 0,
  "io_in" int4 NOT NULL DEFAULT 0,
  "io_out" int4 NOT NULL DEFAULT 0,
  "success_request" int4 NOT NULL DEFAULT 0,
  "error_request" int4 NOT NULL DEFAULT 0,
  "bind_node_mode" int4 NOT NULL DEFAULT 0,
  "create_time" timestamp(6) NOT NULL DEFAULT now()
)
;
COMMENT ON COLUMN "public"."t_task"."project_id" IS '所属项目';
COMMENT ON COLUMN "public"."t_task"."config_snapshot_id" IS '绑定的配置快照id';
COMMENT ON COLUMN "public"."t_task"."status" IS '任务状态（0：暂停，1：进行中，2：停止，3：已完成）';
COMMENT ON COLUMN "public"."t_task"."result_count" IS '结果数';
COMMENT ON COLUMN "public"."t_task"."io_in" IS '流入字节';
COMMENT ON COLUMN "public"."t_task"."io_out" IS '流出字节';
COMMENT ON COLUMN "public"."t_task"."success_request" IS '成功请求数';
COMMENT ON COLUMN "public"."t_task"."error_request" IS '错误请求数';
COMMENT ON COLUMN "public"."t_task"."bind_node_mode" IS '任务绑定模式（0：随机，1：特定）';
COMMENT ON COLUMN "public"."t_task"."create_time" IS '创建时间';
COMMENT ON TABLE "public"."t_task" IS '任务表';

-- ----------------------------
-- Table structure for t_user
-- ----------------------------
DROP TABLE IF EXISTS "public"."t_user";
CREATE TABLE "public"."t_user" (
  "id" int4 NOT NULL,
  "username" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "password" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "avatar" varchar(255) COLLATE "pg_catalog"."default",
  "create_time" timestamp(6) NOT NULL DEFAULT now(),
  "update_time" timestamp(6) NOT NULL DEFAULT now()
)
;
COMMENT ON COLUMN "public"."t_user"."username" IS '用户名或邮箱';
COMMENT ON COLUMN "public"."t_user"."password" IS '登录密码，使用md5加密';
COMMENT ON COLUMN "public"."t_user"."avatar" IS '头像';
COMMENT ON COLUMN "public"."t_user"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."t_user"."update_time" IS '更新时间';
COMMENT ON TABLE "public"."t_user" IS '用户表';

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."seq_config_snapshot"', 1, true);
SELECT setval('"public"."seq_field"', 1, true);
SELECT setval('"public"."seq_plugin"', 1, true);
SELECT setval('"public"."seq_project"', 1, true);
SELECT setval('"public"."seq_result"', 1, true);
SELECT setval('"public"."seq_schedule_queue"', 1, true);
SELECT setval('"public"."seq_stage"', 1, true);
SELECT setval('"public"."seq_task"', 1, true);
SELECT setval('"public"."seq_statistic"', 1, true);

-- ----------------------------
-- Primary Key structure for table t_config
-- ----------------------------
ALTER TABLE "public"."t_config" ADD CONSTRAINT "t_config_pkey" PRIMARY KEY ("key");

-- ----------------------------
-- Indexes structure for table t_config_snapshot
-- ----------------------------
CREATE INDEX "Index_project1" ON "public"."t_config_snapshot" USING btree (
  "project_id" "pg_catalog"."int4_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table t_config_snapshot
-- ----------------------------
ALTER TABLE "public"."t_config_snapshot" ADD CONSTRAINT "t_config_snapshot_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table t_field
-- ----------------------------
ALTER TABLE "public"."t_field" ADD CONSTRAINT "t_field_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table t_plugin
-- ----------------------------
CREATE UNIQUE INDEX "Index_plugin_name" ON "public"."t_plugin" USING btree (
  "project_id" "pg_catalog"."int4_ops" ASC NULLS LAST,
  "name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Indexes structure for table t_project
-- ----------------------------
CREATE UNIQUE INDEX "Index_p_name" ON "public"."t_project" USING btree (
  "name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table t_project
-- ----------------------------
ALTER TABLE "public"."t_project" ADD CONSTRAINT "t_project_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table t_queue
-- ----------------------------
CREATE INDEX "Index_taskid_status" ON "public"."t_queue" USING btree (
  "task_id" "pg_catalog"."int4_ops" ASC NULLS LAST,
  "status" "pg_catalog"."int2_ops" ASC NULLS LAST,
  "expire" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table t_queue
-- ----------------------------
ALTER TABLE "public"."t_queue" ADD CONSTRAINT "t_schedule_queue_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table t_stage
-- ----------------------------
CREATE UNIQUE INDEX "Index_name" ON "public"."t_stage" USING btree (
  "project_id" "pg_catalog"."int4_ops" ASC NULLS LAST,
  "name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table t_stage
-- ----------------------------
ALTER TABLE "public"."t_stage" ADD CONSTRAINT "t_stage_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table t_statistic
-- ----------------------------
CREATE INDEX "Index_sta_time" ON "public"."t_statistic" USING btree (
  "create_time" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);

-- ----------------------------
-- Indexes structure for table t_task
-- ----------------------------
CREATE INDEX "Index_pid_1" ON "public"."t_task" USING btree (
  "project_id" "pg_catalog"."int4_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table t_task
-- ----------------------------
ALTER TABLE "public"."t_task" ADD CONSTRAINT "t_task_pkey" PRIMARY KEY ("id");
