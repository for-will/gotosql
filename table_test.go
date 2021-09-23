package db

import (
	"testing"
	"time"
)

type RewardTask struct {
	ReturnCode  *int32 `db:"index"`
	Id          *int32 `db:"primary_key"`
	ExpAdd      *int32
	Gold        *int32
	Honor       *int32
	Achievement *int32
	Diamond     *int32
}

type TestTableTask struct {
	Id         int32 `db:"name:sn,type:int,index,primary_key"`
	PlayerSn   int32 `db:"unique:idx_player_mission"`
	Mission    int32 `db:"mission,unique:idx_player_mission"`
	State      int8  `db:"unique"`
	Progress   int32
	RewardedAt time.Time `db:"type:timestamp"`
	PsX        int32
	PsY        int32
}

func TestRecreateTable(t *testing.T) {
	RecreateTable(db, TestTableTask{})
}

func TestGenModelAutoFile(t *testing.T) {
	GenModelAutoFile("model_sql_auto.go",
		&TestTableTask{},
		&RewardTask{},
	)
}
