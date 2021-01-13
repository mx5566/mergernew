package model

import (
	"fmt"
	"gorm.io/gorm"
)

func HandleRoleNameOrigin(db1, db2, db3 *gorm.DB) error {
	// 修改role_name_origin字段
	var target_database = db1.Migrator().CurrentDatabase()
	var target_database_copy = db2.Migrator().CurrentDatabase()

	var target_table_name = "role_data"
	var target_table_name_copy = "role_data"
	var column_name = "role_name_origin"

	/*
			select count(*) into @cnt1 FROM information_schema.columns WHERE table_schema = target_database AND table_name = target_table_name AND column_name = target_column_name;
		if @cnt1 = 0 then
				set @st1 = CONCAT('ALTER TABLE ', target_table_name, ' ADD COLUMN ', target_column_name, ' VARCHAR(32) default NULL');
				PREPARE STMT1 FROM @st1;
				EXECUTE STMT1;
				DEALLOCATE PREPARE STMT1;
		end if;

		select count(*) into @cnt2 FROM information_schema.columns WHERE table_schema = target_database AND table_name = target_table_name_copy AND column_name = target_column_name;
		if @cnt2 = 0 then
				set @st2 = CONCAT('ALTER TABLE ', target_table_name_copy, ' ADD COLUMN ', target_column_name, ' VARCHAR(32) default NULL');
				PREPARE STMT2 FROM @st2;
				EXECUTE STMT2;
				DEALLOCATE PREPARE STMT2;
		end if;


	*/

	var c int64 = 0

	// 计算个数
	err := db3.Table("columns").Where("table_schema = ? and table_name = ? and column_name = ?", target_database, target_table_name, column_name).Select("count(*) as count").Count(&c).Error
	if err != nil {
		return err
	}
	if c == 0 {
		err = db1.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s VARCHAR(32) default NULL;", target_table_name, column_name)).Error
		if err != nil {
			return err
		}
	}

	err = db3.Table("columns").Where("table_schema = ? and table_name = ? and column_name = ?", target_database_copy, target_table_name_copy, column_name).Select("count(*) as count").Count(&c).Error
	if err != nil {
		return err
	}
	if c == 0 {
		err = db2.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s VARCHAR(32) default NULL;", target_table_name_copy, column_name)).Error
		if err != nil {
			return err
		}
	}

	// 更新两个role_data表
	err = db1.Exec("update role_data set role_name_origin = role_name where ISNULL(role_name_origin);").Error
	if err != nil {
		return err
	}

	err = db2.Exec("update role_data set role_name_origin = role_name where ISNULL(role_name_origin);").Error
	if err != nil {
		return err
	}

	return nil
}
