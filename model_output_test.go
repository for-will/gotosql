//
// Code generated auto. DO NOT EDIT.

package db

import (
	"database/sql"
	"strings"
)

func NewTblTestTableTask(db *sql.DB) {

	querySql := `create or replace table test_table_task
(
    id          int       not null auto_increment
        primary key,
    player_sn   int       not null,
    mission     int       not null,
    state       tinyint   not null,
    progress    int       not null,
    rewarded_at timestamp not null default from_unixtime(1),
    ps_x        int       not null default -1,
    ps_y        int       not null default -1,
	index idx_state_progress(state, progress) using btree,
	unique index uni_player_mission(player_sn, mission) using btree
)`
	LogSql(querySql)
	_, err := db.Exec(querySql)
	if err != nil {
		LogError("%+v", err)
	}
}

func SaveTestTableTask(db *sql.DB, obj *TestTableTask) error {

	LogSql("INSERT INTO test_table_task(player_sn, mission, state, progress, rewarded_at, ps_x, ps_y) VALUE ('%v', '%v', '%v', '%v', '%v', '%v', '%v')",
		obj.PlayerSn, obj.Mission, obj.State, obj.Progress, obj.RewardedAt, obj.PsX, obj.PsY)
	result, err := db.Exec("INSERT INTO test_table_task(player_sn, mission, state, progress, rewarded_at, ps_x, ps_y) VALUE (?, ?, ?, ?, ?, ?, ?)",
		obj.PlayerSn, obj.Mission, obj.State, obj.Progress, obj.RewardedAt, obj.PsX, obj.PsY)

	if err != nil {
		LogError("db query error: %+v", err)
		return err
	}

	if id, err := result.LastInsertId(); err != nil {
		LogError("get last insert id error: %+v", err)
		return err
	} else {
		obj.Id = int32(id)
	}
	return nil
}

func FirstTestTableTask(db *sql.DB, out *TestTableTask, cond string, args ...interface{}) error {

	sfmt := strings.Replace(cond, "?", "'%+v'", -1)
	LogSql("SELECT id, player_sn, mission, state, progress, rewarded_at, ps_x, ps_y FROM test_table_task WHERE "+sfmt+" LIMIT 1", args...)
	rows, err := db.Query("SELECT id, player_sn, mission, state, progress, rewarded_at, ps_x, ps_y FROM test_table_task WHERE "+cond+" LIMIT 1", args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		LogError("db query error: %+v", err)
		return err
	}

	if rows.Next() {
		if err = rows.Scan(&out.Id, &out.PlayerSn, &out.Mission, &out.State, &out.Progress, &out.RewardedAt, &out.PsX, &out.PsY); err != nil {
			LogError("db query Scan error: %+v", err)
			return err
		}
	} else {
		return sql.ErrNoRows
	}
	return nil
}

func FindTestTableTask(db *sql.DB, out *[]*TestTableTask, cond string, args ...interface{}) error {

	sfmt := strings.Replace(cond, "?", "'%+v'", -1)
	LogSql("SELECT id, player_sn, mission, state, progress, rewarded_at, ps_x, ps_y FROM test_table_task WHERE "+sfmt, args...)
	rows, err := db.Query("SELECT id, player_sn, mission, state, progress, rewarded_at, ps_x, ps_y FROM test_table_task WHERE "+cond, args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		LogError("db query error: %+v", err)
		return err
	}

	for rows.Next() {
		obj := &TestTableTask{}
		if err = rows.Scan(&obj.Id, &obj.PlayerSn, &obj.Mission, &obj.State, &obj.Progress, &obj.RewardedAt, &obj.PsX, &obj.PsY); err != nil {
			LogError("db query Scan error: %+v", err)
			return err
		}
		*out = append(*out, obj)
	}
	return nil
}

func NewTblRewardTask(db *sql.DB) {

	querySql := `create or replace table reward_task
(
    return_code int not null,
    id          int not null auto_increment
        primary key,
    exp_add     int not null,
    gold        int not null,
    honor       int not null,
    achievement int not null,
    diamond     int not null,
	index idx_reward_task_return_code(return_code) using btree
)`
	LogSql(querySql)
	_, err := db.Exec(querySql)
	if err != nil {
		LogError("%+v", err)
	}
}

func SaveRewardTask(db *sql.DB, obj *RewardTask) error {

	LogSql("INSERT INTO reward_task(return_code, exp_add, gold, honor, achievement, diamond) VALUE ('%v', '%v', '%v', '%v', '%v', '%v')",
		obj.ReturnCode, obj.ExpAdd, obj.Gold, obj.Honor, obj.Achievement, obj.Diamond)
	result, err := db.Exec("INSERT INTO reward_task(return_code, exp_add, gold, honor, achievement, diamond) VALUE (?, ?, ?, ?, ?, ?)",
		obj.ReturnCode, obj.ExpAdd, obj.Gold, obj.Honor, obj.Achievement, obj.Diamond)

	if err != nil {
		LogError("db query error: %+v", err)
		return err
	}

	if id, err := result.LastInsertId(); err != nil {
		LogError("get last insert id error: %+v", err)
		return err
	} else {
		obj.Id = int32(id)
	}
	return nil
}

func FirstRewardTask(db *sql.DB, out *RewardTask, cond string, args ...interface{}) error {

	sfmt := strings.Replace(cond, "?", "'%+v'", -1)
	LogSql("SELECT return_code, id, exp_add, gold, honor, achievement, diamond FROM reward_task WHERE "+sfmt+" LIMIT 1", args...)
	rows, err := db.Query("SELECT return_code, id, exp_add, gold, honor, achievement, diamond FROM reward_task WHERE "+cond+" LIMIT 1", args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		LogError("db query error: %+v", err)
		return err
	}

	if rows.Next() {
		if err = rows.Scan(&out.ReturnCode, &out.Id, &out.ExpAdd, &out.Gold, &out.Honor, &out.Achievement, &out.Diamond); err != nil {
			LogError("db query Scan error: %+v", err)
			return err
		}
	} else {
		return sql.ErrNoRows
	}
	return nil
}

func FindRewardTask(db *sql.DB, out *[]*RewardTask, cond string, args ...interface{}) error {

	sfmt := strings.Replace(cond, "?", "'%+v'", -1)
	LogSql("SELECT return_code, id, exp_add, gold, honor, achievement, diamond FROM reward_task WHERE "+sfmt, args...)
	rows, err := db.Query("SELECT return_code, id, exp_add, gold, honor, achievement, diamond FROM reward_task WHERE "+cond, args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		LogError("db query error: %+v", err)
		return err
	}

	for rows.Next() {
		obj := &RewardTask{}
		if err = rows.Scan(&obj.ReturnCode, &obj.Id, &obj.ExpAdd, &obj.Gold, &obj.Honor, &obj.Achievement, &obj.Diamond); err != nil {
			LogError("db query Scan error: %+v", err)
			return err
		}
		*out = append(*out, obj)
	}
	return nil
}