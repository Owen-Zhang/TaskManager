package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	TASK_SUCCESS = 0  // 任务执行成功
	TASK_ERROR   = -1 // 任务执行出错
	TASK_TIMEOUT = -2 // 任务执行超时
)

/*
  header例子：
  aaa=123
  bb=sdfasdasdf
*/
type Task struct {
	Id           int
	UserId       int
	GroupId      int
	TaskName     string
	TaskType     int        //0:文件，1:API, 2:Shell脚本
	FileFolder   string     //当takstype为文件时，上传文件保存的文件夹名，也是shell的文件名
	OldGzipFile  string     //用户前一次上传的文件名
	Command      string     //运行的命令
	ApiHeader    string     //调用接口的header 
	ApiUrl       string     //调用的API地址
	ApiMethod    string     //提交的Method，现只支持GET, POST
	PostBody     string     //Post方式提交的body
	Description  string
	CronSpec     string
	Concurrent   int
	Status       int
	RunStatus    int        //运行状态
	Notify       int
	NotifyEmail  string
	Timeout      int
	ExecuteTimes int
	PrevTime     int64
	CreateTime   int64
}

func (t *Task) TableName() string {
	return TableName("task")
}

func (t *Task) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func TaskAdd(task *Task) (int64, error) {
	if task.CreateTime == 0 {
		task.CreateTime = time.Now().Unix()
	}
	return orm.NewOrm().Insert(task)
}

func TaskGetList(page, pageSize int, filters ...interface{}) ([]*Task, int64) {
	offset := (page - 1) * pageSize

	tasks := make([]*Task, 0)

	query := orm.NewOrm().QueryTable(TableName("task"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&tasks)

	return tasks, total
}

func TaskResetGroupId(groupId int) (int64, error) {
	return orm.NewOrm().QueryTable(TableName("task")).Filter("group_id", groupId).Update(orm.Params{
		"group_id": 0,
	})
}

func TaskGetById(id int) (*Task, error) {
	task := &Task{
		Id: id,
	}

	err := orm.NewOrm().Read(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func TaskDel(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("task")).Filter("id", id).Delete()
	return err
}
