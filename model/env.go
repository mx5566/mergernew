package model

import "gorm.io/gorm"

func SetEnv(isRemote bool, db *gorm.DB) {

	if !isRemote {
		GDB1.Exec("SET group_concat_max_len=15000;")
		GDB1.Exec("set global max_allowed_packet=2 * 1024 * 1024 * 128;")
		GDB1.Exec("set global innodb_buffer_pool_size=2 * 1024 * 1024 * 1024;")
		GDB1.Exec("set tmp_table_size = 2 * 1024 * 1024 * 1024;")
		GDB1.Exec("set GLOBAL innodb_log_buffer_size=16 * 1024 * 1024;")
		GDB1.Exec("set GLOBAL innodb_flush_log_at_trx_commit = 2;")
		GDB1.Exec("set session BULK_INSERT_BUFFER_SIZE=256217728;")
		GDB1.Exec("set GLOBAL max_connections=1000;")
		GDB1.Exec("set global innodb_thread_concurrency=32;")
		GDB1.Exec("set sql_log_bin=OFF;")
		return
	}

	if db != nil {
		db.Exec("SET group_concat_max_len=15000;")
		db.Exec("set global max_allowed_packet=2 * 1024 * 1024 * 128;")
		db.Exec("set global innodb_buffer_pool_size=2 * 1024 * 1024 * 1024;")
		db.Exec("set tmp_table_size = 2 * 1024 * 1024 * 1024;")
		db.Exec("set GLOBAL innodb_log_buffer_size=16 * 1024 * 1024;")
		db.Exec("set GLOBAL innodb_flush_log_at_trx_commit = 2;")
		db.Exec("set session BULK_INSERT_BUFFER_SIZE=256217728;")
		db.Exec("set GLOBAL max_connections=1000;")
		db.Exec("set global innodb_thread_concurrency=32;")
		db.Exec("set sql_log_bin=OFF;")

	}

}
