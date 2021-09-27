package db

import (
	"testing"
)


func TestRecreateTable(t *testing.T) {
	RecreateTable(db, TestTableTask{})
}

func TestGenModelAutoFile(t *testing.T) {
	GenModelAutoFile("model_output_test.go", "db",
		&TestTableTask{},
		&RewardTask{},
	)
}
//
//func FirstRewardTask(db *sql.DB, out *RewardTask, cond string, args ...interface{}) error {
//
//	sfmt := strings.Replace(cond, "?", "'%+v'", -1)
//	fmt.Printf("Exec Sql: SELECT id, player_sn FROM `reward_task` WHERE "+sfmt+" LIMIT 1", args...)
//	rows, err := db.Query("SELECT id, player_sn FROM `reward_task` WHERE "+cond+" LIMIT 1", args...)
//	if rows != nil {
//		defer rows.Close()
//	}
//	if err != nil {
//		fmt.Printf("db query error: %+v", err)
//		return err
//	}
//
//	if rows.Next() {
//		if err = rows.Scan(&out.Id, &out.Diamond); err != nil {
//			fmt.Printf("db query Scan error: %+v", err)
//			return err
//		}
//	}
//	return nil
//}
