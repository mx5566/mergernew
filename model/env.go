package model

func SetEnv() {
	GDB1.Exec("SET group_concat_max_len=15000;")
	GDB1.Exec("set global max_allowed_packet=2 * 1024 * 1024 * 64;")
	GDB1.Exec("set global innodb_buffer_pool_size=2 * 1024 * 1024 * 1024;")
	GDB1.Exec("set tmp_table_size = 2 * 1024 * 1024 * 1024;")
	GDB1.Exec("set GLOBAL innodb_log_buffer_size=16 * 1024 * 1024;")
	GDB1.Exec("set GLOBAL innodb_flush_log_at_trx_commit = 2;")
	GDB1.Exec("set session BULK_INSERT_BUFFER_SIZE=256217728;")
	GDB1.Exec("SET group_concat_max_len=15000;")
	GDB1.Exec("set GLOBAL max_connections=1000;")
	GDB1.Exec("set global innodb_thread_concurrency=32;")

	GDB2.Exec("SET group_concat_max_len=15000;")
	GDB2.Exec("set global max_allowed_packet=2 * 1024 * 1024 * 64;")
	GDB2.Exec("set global innodb_buffer_pool_size=2 * 1024 * 1024 * 1024;")
	GDB2.Exec("set tmp_table_size = 2 * 1024 * 1024 * 1024;")
	GDB2.Exec("set GLOBAL innodb_log_buffer_size=16 * 1024 * 1024;")
	GDB2.Exec("set GLOBAL innodb_flush_log_at_trx_commit = 2;")
	GDB2.Exec("set session BULK_INSERT_BUFFER_SIZE=256217728;")
	GDB2.Exec("SET group_concat_max_len=15000;")
	GDB2.Exec("set GLOBAL max_connections=1000;")
	GDB2.Exec("set global innodb_thread_concurrency=32;")

	//runtime.GOMAXPROCS(runtime.NumCPU())
}
