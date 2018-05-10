package jobs

import (
	"bytes"
	"os/exec"
	//"app/mail"
	//"html/template"
	
	"fmt"
	"errors"
	"io/ioutil"	
	"app/models"
	//"runtime/debug"
	"strings"
	"time"
	"net/http"
	"github.com/axgle/mahonia"
	"github.com/imroc/req"
	"github.com/astaxie/beego"
)

type Job struct {
	id         int                                               // 任务ID
	logId      int64                                             // 日志记录ID
	name       string                                            // 任务名称
	task       *models.Task                                      // 任务对象
	runFunc    func(time.Duration) (string, string, error, bool) // 执行函数
	status     int                                               // 任务状态，大于0表示正在执行中
	Concurrent bool                                              // 同一个任务是否允许并行执行
}

func NewJobFromTask(task *models.Task) (*Job, error) {
	if task.Id < 1 {
		return nil, fmt.Errorf("ToJob: 缺少id")
	}
	job := NewCommandJob(task)
	job.task = task
	job.Concurrent = task.Concurrent == 1
	return job, nil
}

func NewCommandJob(task *models.Task) *Job {
	job := &Job{
		id:   task.Id,
		name: task.TaskName,
	}
	job.runFunc = func(timeout time.Duration) (string, string, error, bool) {
		//TaskType 0:文件, 1: API, 2:Shell脚本
		if task.TaskType == 0 {
			bufOut := new(bytes.Buffer)
			bufErr := new(bytes.Buffer)
			runShell := fmt.Sprintf("%s/%s.sh", models.RunDir, task.FileFolder)
			cmd := exec.Command("/bin/bash", "-c", runShell)
			
			cmd.Stdout = bufOut
			cmd.Stderr = bufErr
			cmd.Start()
			err, isTimeout := runCmdWithTimeout(cmd, timeout)
	
			encoder := mahonia.NewDecoder("gbk")
			return encoder.ConvertString(bufOut.String()), encoder.ConvertString(bufErr.String()), err, isTimeout
		
		} else if task.TaskType == 1 {
			header := make(http.Header)
			if task.ApiHeader != "" && strings.TrimSpace(task.ApiHeader) != "" {
				headers := strings.Split(task.ApiHeader, "\n")
				for _,val := range headers {
					keyval := strings.Split(val, "=")
					if len(keyval) > 0 {
						v := strings.TrimSpace(keyval[0])
						v1 := strings.TrimSpace(keyval[1])
						if v != "" && v1 != "" {
							header.Set(v, v1)
						} else {
							continue
						}
					}
				}
			}
			responsestr := ""
			var err error
			var res *req.Resp
			
			//要支持content-type:urlencode
			req.SetTimeout(time.Second * time.Duration(task.Timeout))
			if task.ApiMethod == "POST" {
				if task.PostBody != "" {
					contenttype := header.Get("Content-Type")			
					//如果没有设置就用json方式提交
					if contenttype == "" || contenttype == "application/json" {
						res, err = req.Post(task.ApiUrl, header, req.BodyJSON(task.PostBody))
					} else if contenttype == "application/xml" {
						res, err = req.Post(task.ApiUrl, header, req.BodyXML(task.PostBody))
					} else  {
						// application/x-www-form-urlencoded
						res, err = req.Post(task.ApiUrl, header, task.PostBody)
					}
				} else {
					res, err = req.Post(task.ApiUrl, header)
				}
				
			} else {
				res, err = req.Get(task.ApiUrl, header)
			}
			
			if err == nil {
				bodystr, _ := ioutil.ReadAll(res.Response().Body)
				defer res.Response().Body.Close()
	
				responsestr = string(bodystr)
				
				if res.Response().StatusCode != 200 {
					return responsestr, "", errors.New(fmt.Sprintf("返回的状态码为：%s", res.Response().StatusCode)), false
				}
				
				return responsestr, "", nil, false
			} else {
				return "", "", err, false
			}		
		} else {
			//shell 脚本 
			bufOut := new(bytes.Buffer)
			bufErr := new(bytes.Buffer)
			cmd := exec.Command("/bin/bash", "-c", task.Command)
			
			cmd.Stdout = bufOut
			cmd.Stderr = bufErr
			cmd.Start()
			err, isTimeout := runCmdWithTimeout(cmd, timeout)
	
			encoder := mahonia.NewDecoder("gbk")
			return encoder.ConvertString(bufOut.String()), encoder.ConvertString(bufErr.String()), err, isTimeout
		}	
	}
	return job
}

func (j *Job) Status() int {
	return j.status
}

func (j *Job) GetName() string {
	return j.name
}

func (j *Job) GetId() int {
	return j.id
}

func (j *Job) GetLogId() int64 {
	return j.logId
}

func (j *Job) Run() {
	if !j.Concurrent && j.status > 0 {
		beego.Warn(fmt.Sprintf("任务[%d]上一次执行尚未结束，本次被忽略。\n", j.id))
		return
	}

	defer func() {
		if err := recover(); err != nil {
			beego.Error(err, "\n")
		}
	}()

	if workPool != nil {
		workPool <- true
		defer func() {
			<-workPool
		}()
	}

	beego.Debug(fmt.Sprintf("开始执行任务: %d\n", j.id))

	j.status++
	defer func() {
		j.status--
	}()

	t := time.Now()
	timeout := time.Duration(time.Hour * 24)
	if j.task.Timeout > 0 {
		timeout = time.Second * time.Duration(j.task.Timeout)
	}

	cmdOut, cmdErr, err, isTimeout := j.runFunc(timeout)

	ut := time.Now().Sub(t) / time.Millisecond

	// 插入日志
	log := new(models.TaskLog)
	log.TaskId = j.id
	log.Output = cmdOut
	log.Error = cmdErr
	log.ProcessTime = int(ut)
	log.CreateTime = t.Unix()

	if isTimeout {
		log.Status = models.TASK_TIMEOUT
		log.Error = fmt.Sprintf("任务执行超过 %d 秒\n----------------------\n%s\n", int(timeout/time.Second), cmdErr)
	} else if err != nil {
		log.Status = models.TASK_ERROR
		log.Error = err.Error() + ":" + cmdErr
	}
	
	j.logId, _ = models.TaskLogAdd(log)

	// 更新上次执行时间
	j.task.PrevTime = t.Unix()
	j.task.ExecuteTimes++
	j.task.Update("PrevTime", "ExecuteTimes")
}
