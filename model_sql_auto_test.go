package db

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"
)

func TestNewTableRewardLootMissionRs(t *testing.T) {
	//NewTableRewardTask(db)
}

func TestFirstRewardTask(t *testing.T) {
	task := &TestTableTask{}
	err := FirstTestTableTask(db, task, "id=?", 11)
	if err != nil {
		t.Log(err)
	}
	t.Log(JsonString(task))
	t.Log(task.RewardedAt.Unix())
}

func TestSaveTestTableTask(t *testing.T) {
	NewTblTestTableTask(db)
	task := &TestTableTask{
		Id:         0,
		PlayerSn:   5566,
		Mission:    101,
		State:      2,
		Progress:   1,
		RewardedAt: time.Now(),
		PsX:        -1,
		PsY:        -1,
	}
	err := SaveTestTableTask(db, task)
	if err != nil {
		t.Log(err)
	}
	t.Logf("%+v", task)
}

func TestGenModelFunc(t *testing.T) {
	GenModelAutoFile("model_output_test.go", "db",
		&TestTableTask{},
		&RewardTask{},
	)
}

func TestFindTestTableTask(t *testing.T) {
	var task []*TestTableTask
	//SaveTestTableTask(db, &TestTableTask{})
	NewTblTestTableTask(db)
	for {

		SaveTestTableTask(db, &TestTableTask{PlayerSn: rand.Int31n(1000)})
		err := FindTestTableTask(db, &task, "player_sn=? or player_sn=?", rand.Int31n(1000), rand.Int31n(1000))
		if err != nil {
			t.Log(err)
		}
		t.Log(JsonString(task))
		if len(task) != 0 {
			break
		}
	}
}

func JsonString(v interface{}) string {
	bs, _ := json.MarshalIndent(v, "", "\t")
	return string(bs)
}
