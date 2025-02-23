// Copyright 2022 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sysview

import (
	"context"
	"fmt"
	"time"

	"github.com/matrixorigin/matrixone/pkg/logutil"
	ie "github.com/matrixorigin/matrixone/pkg/util/internalExecutor"
	"github.com/matrixorigin/matrixone/pkg/util/trace"
)

const (
	MysqlDBConst       = "mysql"
	InformationDBConst = "information_schema"
	sqlCreateDBConst   = "create database if not exists "
	sqlUseDbConst      = "use "
)

var (
	InitMysqlSysTables = []string{
		`CREATE TABLE IF NOT EXISTS user (
			Host char(255)  NOT NULL DEFAULT '',
			User char(32)  NOT NULL DEFAULT '',
			Select_priv varchar(10) NOT NULL DEFAULT 'N',
			Insert_priv varchar(10) NOT NULL DEFAULT 'N',
			Update_priv varchar(10) NOT NULL DEFAULT 'N',
			Delete_priv varchar(10) NOT NULL DEFAULT 'N',
			Create_priv varchar(10) NOT NULL DEFAULT 'N',
			Drop_priv varchar(10)  NOT NULL DEFAULT 'N',
			Reload_priv varchar(10)  NOT NULL DEFAULT 'N',
			Shutdown_priv varchar(10)  NOT NULL DEFAULT 'N',
			Process_priv varchar(10)  NOT NULL DEFAULT 'N',
			File_priv varchar(10)  NOT NULL DEFAULT 'N',
			Grant_priv varchar(10)  NOT NULL DEFAULT 'N',
			References_priv varchar(10)  NOT NULL DEFAULT 'N',
			Index_priv varchar(10)  NOT NULL DEFAULT 'N',
			Alter_priv varchar(10)  NOT NULL DEFAULT 'N',
			Show_db_priv varchar(10)  NOT NULL DEFAULT 'N',
			Super_priv varchar(10)  NOT NULL DEFAULT 'N',
			Create_tmp_table_priv varchar(10)  NOT NULL DEFAULT 'N',
			Lock_tables_priv varchar(10)  NOT NULL DEFAULT 'N',
			Execute_priv varchar(10)  NOT NULL DEFAULT 'N',
			Repl_slave_priv varchar(10)  NOT NULL DEFAULT 'N',
			Repl_client_priv varchar(10)  NOT NULL DEFAULT 'N',
			Create_view_priv varchar(10)  NOT NULL DEFAULT 'N',
			Show_view_priv varchar(10)  NOT NULL DEFAULT 'N',
			Create_routine_priv varchar(10)  NOT NULL DEFAULT 'N',
			Alter_routine_priv varchar(10)  NOT NULL DEFAULT 'N',
			Create_user_priv varchar(10)  NOT NULL DEFAULT 'N',
			Event_priv varchar(10)  NOT NULL DEFAULT 'N',
			Trigger_priv varchar(10)  NOT NULL DEFAULT 'N',
			Create_tablespace_priv varchar(10)  NOT NULL DEFAULT 'N',
			ssl_type varchar(10)  NOT NULL DEFAULT '',
			ssl_cipher blob NOT NULL,
			x509_issuer blob NOT NULL,
			x509_subject blob NOT NULL,
			max_questions int unsigned NOT NULL DEFAULT '0',
			max_updates int unsigned NOT NULL DEFAULT '0',
			max_connections int unsigned NOT NULL DEFAULT '0',
			max_user_connections int unsigned NOT NULL DEFAULT '0',
			plugin char(64)  NOT NULL DEFAULT 'caching_sha2_password',
			authentication_string text ,
			password_expired varchar(10)  NOT NULL DEFAULT 'N',
			password_last_changed timestamp NULL DEFAULT NULL,
			password_lifetime smallint unsigned DEFAULT NULL,
			account_locked varchar(10)  NOT NULL DEFAULT 'N',
			Create_role_priv varchar(10)  NOT NULL DEFAULT 'N',
			Drop_role_priv varchar(10)  NOT NULL DEFAULT 'N',
			Password_reuse_history smallint unsigned DEFAULT NULL,
			Password_reuse_time smallint unsigned DEFAULT NULL,
			Password_require_current varchar(10)  DEFAULT NULL,
			User_attributes json DEFAULT NULL,
			PRIMARY KEY (Host,User)
		  );`,
		`CREATE TABLE IF NOT EXISTS db (
			Host char(255) NOT NULL DEFAULT '',
			Db char(64)  NOT NULL DEFAULT '',
			User char(32)  NOT NULL DEFAULT '',
			Select_priv varchar(10)  NOT NULL DEFAULT 'N',
			Insert_priv varchar(10)  NOT NULL DEFAULT 'N',
			Update_priv varchar(10)  NOT NULL DEFAULT 'N',
			Delete_priv varchar(10)  NOT NULL DEFAULT 'N',
			Create_priv varchar(10)  NOT NULL DEFAULT 'N',
			Drop_priv varchar(10)  NOT NULL DEFAULT 'N',
			Grant_priv varchar(10)  NOT NULL DEFAULT 'N',
			References_priv varchar(10)  NOT NULL DEFAULT 'N',
			Index_priv varchar(10)  NOT NULL DEFAULT 'N',
			Alter_priv varchar(10)  NOT NULL DEFAULT 'N',
			Create_tmp_table_priv varchar(10)  NOT NULL DEFAULT 'N',
			Lock_tables_priv varchar(10)  NOT NULL DEFAULT 'N',
			Create_view_priv varchar(10)  NOT NULL DEFAULT 'N',
			Show_view_priv varchar(10)  NOT NULL DEFAULT 'N',
			Create_routine_priv varchar(10)  NOT NULL DEFAULT 'N',
			Alter_routine_priv varchar(10)  NOT NULL DEFAULT 'N',
			Execute_priv varchar(10)  NOT NULL DEFAULT 'N',
			Event_priv varchar(10)  NOT NULL DEFAULT 'N',
			Trigger_priv varchar(10)  NOT NULL DEFAULT 'N',
			PRIMARY KEY (Host,Db,User),
			KEY User (User)
		  );`,
		`CREATE TABLE IF NOT EXISTS procs_priv (
			Host char(255)  NOT NULL DEFAULT '',
			Db char(64)  NOT NULL DEFAULT '',
			User char(32)  NOT NULL DEFAULT '',
			Routine_name char(64)  NOT NULL DEFAULT '',
			Routine_type varchar(10)  NOT NULL,
			Grantor varchar(288)  NOT NULL DEFAULT '',
			Proc_priv varchar(10)  NOT NULL DEFAULT '',
			Timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (Host,Db,User,Routine_name,Routine_type),
			KEY Grantor (Grantor)
		  );`,
		`CREATE TABLE IF NOT EXISTS columns_priv (
			Host char(255)  NOT NULL DEFAULT '',
			Db char(64)  NOT NULL DEFAULT '',
			User char(32)  NOT NULL DEFAULT '',
			Table_name char(64)  NOT NULL DEFAULT '',
			Column_name char(64)  NOT NULL DEFAULT '',
			Timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			Column_priv varchar(10) NOT NULL DEFAULT '',
			PRIMARY KEY (Host,Db,User,Table_name,Column_name)
		  );`,
		`CREATE TABLE IF NOT EXISTS tables_priv (
			Host char(255)  NOT NULL DEFAULT '',
			Db char(64)  NOT NULL DEFAULT '',
			User char(32)  NOT NULL DEFAULT '',
			Table_name char(64)  NOT NULL DEFAULT '',
			Grantor varchar(288)  NOT NULL DEFAULT '',
			Timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			Table_priv varchar(10) NOT NULL DEFAULT '',
			Column_priv varchar(10) NOT NULL DEFAULT '',
			PRIMARY KEY (Host,Db,User,Table_name),
			KEY Grantor (Grantor)
		  );`,
	}
	InitInformationSchemaSysTables = []string{
		"CREATE TABLE IF NOT EXISTS KEY_COLUMN_USAGE(" +
			"CONSTRAINT_CATALOG varchar(64)," +
			"CONSTRAINT_SCHEMA varchar(64)," +
			"CONSTRAINT_NAME varchar(64)," +
			"TABLE_CATALOG varchar(64)," +
			"TABLE_SCHEMA varchar(64)," +
			"TABLE_NAME varchar(64)," +
			"COLUMN_NAME varchar(64)," +
			"ORDINAL_POSITION int unsigned," +
			"POSITION_IN_UNIQUE_CONSTRAINT int unsigned," +
			"REFERENCED_TABLE_SCHEMA varchar(64)," +
			"REFERENCED_TABLE_NAME varchar(64)," +
			"REFERENCED_COLUMN_NAME varchar(64)" +
			");",
		"CREATE VIEW COLUMNS AS select " +
			"'def' as TABLE_CATALOG," +
			"att_database as TABLE_SCHEMA," +
			"att_relname AS TABLE_NAME," +
			"attname AS COLUMN_NAME," +
			"attnum AS ORDINAL_POSITION," +
			"mo_show_visible_bin(att_default,1) as COLUMN_DEFAULT," +
			"(case when attnotnull != 0 then 'YES' else 'NO' end) as IS_NULLABLE," +
			"mo_show_visible_bin(atttyp,2) as DATA_TYPE," +
			"att_length as CHARACTER_MAXIMUM_LENGTH," +
			"att_length as CHARACTER_OCTET_LENGTH," +
			"0 as NUMERIC_PRECISION," +
			"0 as NUMERIC_SCALE," +
			"0 as DATETIME_PRECISION," +
			"'utf8' as CHARACTER_SET_NAME," +
			"'utf8_bin' as COLLATION_NAME," +
			"mo_show_visible_bin(atttyp,3) as COLUMN_TYPE," +
			"case when att_constraint_type = 'p' then 'PRI' else '' end as COLUMN_KEY," +
			"case when att_is_auto_increment = 1 then 'auto_increment' else '' end as EXTRA," +
			"'select,insert,update,references' as `PRIVILEGES`," +
			"att_comment as COLUMN_COMMENT," +
			"'' as GENERATION_EXPRESSION," +
			"0 as SRS_ID " +
			"from mo_catalog.mo_columns where att_relname!='%!%mo_increment_columns' and att_relname not like '__mo_cpkey%' and attname != '__mo_rowid'",
		//"CREATE TABLE IF NOT EXISTS COLUMNS(" +
		//	"TABLE_CATALOG varchar(64)," +
		//	"TABLE_SCHEMA varchar(64)," +
		//	"TABLE_NAME varchar(64)," +
		//	"COLUMN_NAME varchar(64)," +
		//	"ORDINAL_POSITION int unsigned," +
		//	"COLUMN_DEFAULT text," +
		//	"IS_NULLABLE varchar(3)," +
		//	"DATA_TYPE longtext," +
		//	"CHARACTER_MAXIMUM_LENGTH bigint," +
		//	"CHARACTER_OCTET_LENGTH bigint," +
		//	"NUMERIC_PRECISION bigint unsigned," +
		//	"NUMERIC_SCALE bigint unsigned," +
		//	"DATETIME_PRECISION int unsigned," +
		//	"CHARACTER_SET_NAME varchar(64)," +
		//	"COLLATION_NAME varchar(64)," +
		//	"COLUMN_TYPE mediumtext," +
		//	"COLUMN_KEY varchar(10)," +
		//	"EXTRA varchar(256)," +
		//	"`PRIVILEGES` varchar(154)," +
		//	"COLUMN_COMMENT text," +
		//	"GENERATION_EXPRESSION longtext," +
		//	"SRS_ID int unsigned" +
		//	");",

		"CREATE TABLE IF NOT EXISTS PROFILING (" +
			"QUERY_ID int NOT NULL DEFAULT '0'," +
			"SEQ int NOT NULL DEFAULT '0'," +
			"STATE varchar(30) NOT NULL DEFAULT ''," +
			"DURATION decimal(9,6) NOT NULL DEFAULT '0.000000'," +
			"CPU_USER decimal(9,6) DEFAULT NULL," +
			"CPU_SYSTEM decimal(9,6) DEFAULT NULL," +
			"CONTEXT_VOLUNTARY int DEFAULT NULL," +
			"CONTEXT_INVOLUNTARY int DEFAULT NULL," +
			"BLOCK_OPS_IN int DEFAULT NULL," +
			"BLOCK_OPS_OUT int DEFAULT NULL," +
			"MESSAGES_SENT int DEFAULT NULL," +
			"MESSAGES_RECEIVED int DEFAULT NULL," +
			"PAGE_FAULTS_MAJOR int DEFAULT NULL," +
			"PAGE_FAULTS_MINOR int DEFAULT NULL," +
			"SWAPS int DEFAULT NULL," +
			"SOURCE_FUNCTION varchar(30) DEFAULT NULL," +
			"SOURCE_FILE varchar(20) DEFAULT NULL," +
			"SOURCE_LINE int DEFAULT NULL" +
			");",
		"CREATE TABLE IF NOT EXISTS `PROCESSLIST` (" +
			"ID bigint unsigned NOT NULL DEFAULT '0'," +
			"USER varchar(32) NOT NULL DEFAULT ''," +
			"HOST varchar(261) NOT NULL DEFAULT ''," +
			"DB varchar(64) DEFAULT NULL," +
			"COMMAND varchar(16) NOT NULL DEFAULT ''," +
			"TIME int NOT NULL DEFAULT '0'," +
			"STATE varchar(64) DEFAULT NULL," +
			"INFO longtext" +
			");",
		"CREATE TABLE IF NOT EXISTS USER_PRIVILEGES (" +
			"GRANTEE varchar(292) NOT NULL DEFAULT ''," +
			"TABLE_CATALOG varchar(512) NOT NULL DEFAULT ''," +
			"PRIVILEGE_TYPE varchar(64) NOT NULL DEFAULT ''," +
			"IS_GRANTABLE varchar(3) NOT NULL DEFAULT ''" +
			");",
		//"CREATE TABLE IF NOT EXISTS SCHEMATA (" +
		//	"CATALOG_NAME varchar(64)," +
		//	"SCHEMA_NAME varchar(64)," +
		//	"DEFAULT_CHARACTER_SET_NAME varchar(64)," +
		//	"DEFAULT_COLLATION_NAME varchar(64)," +
		//	"SQL_PATH binary(0)," +
		//	"DEFAULT_ENCRYPTION varchar(10)" +
		//	");",
		"CREATE VIEW SCHEMATA AS SELECT " +
			"dat_catalog_name AS CATALOG_NAME," +
			"datname AS SCHEMA_NAME," +
			"'utf8mb4' AS DEFAULT_CHARACTER_SET_NAME," +
			"'utf8mb4_0900_ai_ci' AS DEFAULT_COLLATION_NAME," +
			"'' AS SQL_PATH," +
			"'NO' AS DEFAULT_ENCRYPTION " +
			"FROM mo_catalog.mo_database;",
		"CREATE TABLE IF NOT EXISTS CHARACTER_SETS(" +
			"CHARACTER_SET_NAME varchar(64)," +
			"DEFAULT_COLLATE_NAME varchar(64)," +
			"DESCRIPTION varchar(2048)," +
			"MAXLEN int unsigned" +
			");",
		"CREATE TABLE IF NOT EXISTS TRIGGERS(" +
			"TRIGGER_CATALOG varchar(64)," +
			"TRIGGER_SCHEMA varchar(64)," +
			"TRIGGER_NAME varchar(64)," +
			"EVENT_MANIPULATION varchar(10)," +
			"EVENT_OBJECT_CATALOG varchar(64)," +
			"EVENT_OBJECT_SCHEMA varchar(64)," +
			"EVENT_OBJECT_TABLE varchar(64)," +
			"ACTION_ORDER int unsigned," +
			"ACTION_CONDITION binary(0)," +
			"ACTION_STATEMENT longtext," +
			"ACTION_ORIENTATION varchar(3)," +
			"ACTION_TIMING varchar(10)," +
			"ACTION_REFERENCE_OLD_TABLE binary(0)," +
			"ACTION_REFERENCE_NEW_TABLE binary(0)," +
			"ACTION_REFERENCE_OLD_ROW varchar(3)," +
			"ACTION_REFERENCE_NEW_ROW varchar(3)," +
			"CREATED timestamp(2)," +
			"SQL_MODE varchar(10)," +
			"DEFINER varchar(288)," +
			"CHARACTER_SET_CLIENT varchar(64)," +
			"COLLATION_CONNECTION varchar(64)," +
			"DATABASE_COLLATION varchar(64)" +
			");",

		"CREATE VIEW IF NOT EXISTS TABLES AS " +
			"SELECT 'def' AS TABLE_CATALOG," +
			"reldatabase AS TABLE_SCHEMA," +
			"relname AS TABLE_NAME," +
			"'' AS TABLE_TYPE," +
			"'' AS ENGINE," +
			"0 AS VERSION," +
			"'' AS ROW_FORMAT," +
			"0 AS TABLE_ROWS," +
			"0 AS AVG_ROW_LENGTH," +
			"0 AS DATA_LENGTH," +
			"0 AS MAX_DATA_LENGTH," +
			"0 AS INDEX_LENGTH," +
			"0 AS DATA_FREE," +
			"0 AS `AUTO_INCREMENT`," +
			"created_time AS CREATE_TIME," +
			"created_time AS UPDATE_TIME," +
			"created_time AS CHECK_TIME," +
			"'' AS TABLE_COLLATION," +
			"0 AS CHECKSUM," +
			"'' AS CREATE_OPTIONS," +
			"rel_comment AS TABLE_COMMENT " +
			"FROM mo_catalog.mo_tables;",

		//"CREATE TABLE TABLES(" +
		//	"TABLE_CATALOG varchar(64)," +
		//	"TABLE_SCHEMA varchar(64)," +
		//	"TABLE_NAME varchar(64)," +
		//	"TABLE_TYPE varchar(50)," +
		//	"ENGINE varchar(64), " +
		//	"VERSION int," +
		//	"ROW_FORMAT varchar(50)," +
		//	"TABLE_ROWS bigint unsigned," +
		//	"AVG_ROW_LENGTH bigint unsigned," +
		//	"DATA_LENGTH bigint unsigned," +
		//	"MAX_DATA_LENGTH bigint unsigned," +
		//	"INDEX_LENGTH bigint unsigned," +
		//	"DATA_FREE bigint unsigned," +
		//	"`AUTO_INCREMENT` bigint unsigned," +
		//	"CREATE_TIME timestamp," +
		//	"UPDATE_TIME datetime," +
		//	"CHECK_TIME datetime," +
		//	"TABLE_COLLATION varchar(64)," +
		//	"CHECKSUM bigint," +
		//	"CREATE_OPTIONS varchar(256)," +
		//	"TABLE_COMMENT text" +
		//	");",

		"CREATE TABLE IF NOT EXISTS ENGINES (" +
			"ENGINE varchar(64)," +
			"SUPPORT varchar(8)," +
			"COMMENT varchar(160)," +
			"TRANSACTIONS varchar(3)," +
			"XA varchar(3)," +
			"SAVEPOINTS varchar(3)" +
			");",
		"CREATE TABLE IF NOT EXISTS ROUTINES (" +
			"SPECIFIC_NAME varchar(64)," +
			"ROUTINE_CATALOG varchar(64)," +
			"ROUTINE_SCHEMA varchar(64)," +
			"ROUTINE_NAME varchar(64)," +
			"ROUTINE_TYPE varchar(10)," +
			"DATA_TYPE  longtext," +
			"CHARACTER_MAXIMUM_LENGTH bigint," +
			"CHARACTER_OCTET_LENGTH bigint," +
			"NUMERIC_PRECISION int unsigned," +
			"NUMERIC_SCALE int unsigned," +
			"DATETIME_PRECISION int unsigned," +
			"CHARACTER_SET_NAME varchar(64)," +
			"COLLATION_NAME varchar(64)," +
			"DTD_IDENTIFIER longtext," +
			"ROUTINE_BODY varchar(3)," +
			"ROUTINE_DEFINITION longtext," +
			"EXTERNAL_NAME binary(0)," +
			"EXTERNAL_LANGUAGE varchar(64)," +
			"PARAMETER_STYLE varchar(3)," +
			"IS_DETERMINISTIC varchar(3)," +
			"SQL_DATA_ACCESS varchar(10)," +
			"SQL_PATH varchar(1000)," +
			"SECURITY_TYPE varchar(10)," +
			"CREATED timestamp," +
			"LAST_ALTERED timestamp," +
			"SQL_MODE varchar(1000)," +
			"ROUTINE_COMMENT text," +
			"DEFINER varchar(288)," +
			"CHARACTER_SET_CLIENT varchar(64)," +
			"COLLATION_CONNECTION varchar(64)," +
			"DATABASE_COLLATION  varchar(64)" +
			");",
		"CREATE TABLE IF NOT EXISTS PARAMETERS(" +
			"SPECIFIC_CATALOG varchar(64)," +
			"SPECIFIC_SCHEMA varchar(64)," +
			"SPECIFIC_NAME varchar(64)," +
			"ORDINAL_POSITION bigint unsigned," +
			"PARAMETER_MODE varchar(5)," +
			"PARAMETER_NAME varchar(64)," +
			"DATA_TYPE longtext," +
			"CHARACTER_MAXIMUM_LENGTH bigint," +
			"CHARACTER_OCTET_LENGTH bigint," +
			"NUMERIC_PRECISION int unsigned," +
			"NUMERIC_SCALE bigint," +
			"DATETIME_PRECISION int unsigned," +
			"CHARACTER_SET_NAME varchar(64)," +
			"COLLATION_NAME varchar(64)," +
			"DTD_IDENTIFIER mediumtext," +
			"ROUTINE_TYPE  varchar(64)" +
			");",
		"CREATE TABLE IF NOT EXISTS KEYWORDS (" +
			"WORD varchar(64)," +
			"RESERVED int unsigned" +
			");",
		informationSchemaKeywordsData,
	}
)

func InitSchema(ctx context.Context, ieFactory func() ie.InternalExecutor) error {
	initMysqlTables(ctx, ieFactory, trace.FileService)
	initInformationSchemaTables(ctx, ieFactory, trace.FileService)
	return nil
}

func initMysqlTables(ctx context.Context, ieFactory func() ie.InternalExecutor, batchProcessMode string) {
	exec := ieFactory()
	exec.ApplySessionOverride(ie.NewOptsBuilder().Database(MysqlDBConst).Internal(true).Finish())
	mustExec := func(sql string) {
		if err := exec.Exec(ctx, sql, ie.NewOptsBuilder().Finish()); err != nil {
			panic(fmt.Sprintf("[Mysql] init mysql tables error: %v, sql: %s", err, sql))
		}
	}
	mustExec(sqlCreateDBConst + MysqlDBConst)
	mustExec(sqlUseDbConst + MysqlDBConst)
	var createCost time.Duration
	defer func() {
		logutil.Debugf("[Mysql] init mysql tables: create cost %d ms", createCost.Milliseconds())
	}()
	instant := time.Now()

	for _, sql := range InitMysqlSysTables {
		mustExec(sql)
	}
	createCost = time.Since(instant)
}

func initInformationSchemaTables(ctx context.Context, ieFactory func() ie.InternalExecutor, batchProcessMode string) {
	exec := ieFactory()
	exec.ApplySessionOverride(ie.NewOptsBuilder().Database(InformationDBConst).Internal(true).Finish())
	mustExec := func(sql string) {
		if err := exec.Exec(ctx, sql, ie.NewOptsBuilder().Finish()); err != nil {
			panic(fmt.Sprintf("[information_schema] init information_schema tables error: %v, sql: %s", err, sql))
		}
	}
	mustExec(sqlCreateDBConst + InformationDBConst)
	mustExec(sqlUseDbConst + InformationDBConst)
	var createCost time.Duration
	defer func() {
		logutil.Debugf("[information_schema] init information_schema tables: create cost %d ms", createCost.Milliseconds())
	}()
	instant := time.Now()

	for _, sql := range InitInformationSchemaSysTables {
		mustExec(sql)
	}
	createCost = time.Since(instant)
}
