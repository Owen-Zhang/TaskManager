package controllers

import (
	"app/jobs"
	"app/libs"
	"app/models"
	"app/models/response"
	"strconv"
	"strings"
	"time"
	"fmt"
	"os"
	"regexp"
	"github.com/astaxie/beego"
	"github.com/robfig/cron"
)

type TaskController struct {
	BaseController
}

// 任务列表
func (this *TaskController) List() {
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	groupId, _ := this.GetInt("groupid")
	filters := make([]interface{}, 0)
	if groupId > 0 {
		filters = append(filters, "group_id", groupId)
	}
	result, count := models.TaskGetList(page, this.pageSize, filters...)

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["name"] = v.TaskName
		row["cron_spec"] = v.CronSpec
		row["status"] = v.Status
		row["description"] = v.Description

		e := jobs.GetEntryById(v.Id)
		if e != nil {
			row["next_time"] = beego.Date(e.Next, "Y-m-d H:i:s")
			row["prev_time"] = "-"
			if e.Prev.Unix() > 0 {
				row["prev_time"] = beego.Date(e.Prev, "Y-m-d H:i:s")
			} else if v.PrevTime > 0 {
				row["prev_time"] = beego.Date(time.Unix(v.PrevTime, 0), "Y-m-d H:i:s")
			}
			row["running"] = 1
		} else {
			row["next_time"] = "-"
			if v.PrevTime > 0 {
				row["prev_time"] = beego.Date(time.Unix(v.PrevTime, 0), "Y-m-d H:i:s")
			} else {
				row["prev_time"] = "-"
			}
			row["running"] = 0
		}
		list[k] = row
	}

	// 分组列表
	groups, _ := models.TaskGroupGetList(1, 100)

	this.Data["pageTitle"] = "任务列表"
	this.Data["list"] = list
	this.Data["groups"] = groups
	this.Data["groupid"] = groupId
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("TaskController.List", "groupid", groupId), true).ToString()
	this.display()
}

// 添加任务
func (this *TaskController) Add() {
	groups, _ := models.TaskGroupGetList(1, 100)
	this.Data["groups"] = groups
	this.Data["pageTitle"] = "添加任务"
	this.display()
}

// 编辑任务
func (this *TaskController) Edit() {
	id, _ := this.GetInt("id")

	task, err := models.TaskGetById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	// 分组列表
	groups, _ := models.TaskGroupGetList(1, 100)
	this.Data["groups"] = groups
	this.Data["task"] = task
	this.Data["pageTitle"] = "编辑任务"
	this.display("task/add")
}

//查看任务
func (this *TaskController) View() {
	id, _ := this.GetInt("id")

	task, err := models.TaskGetById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	groups, _ := models.TaskGroupGetList(1, 100)
	this.Data["groups"] = groups
	this.Data["task"] = task
	this.Data["pageTitle"] = "编辑任务"
	this.Data["isview"] = 1
	this.display("task/add")
}

//保存上传的文件
func (this *TaskController) UploadRunFile() {
	f, h, err := this.GetFile("files[]")
	defer f.Close()

	uploadResult := &response.ResultData{IsSuccess: false,Msg: ""}

	if err != nil {
		uploadResult.Msg = "请选择要上传的文件,如:demo.tar.gz"
		this.jsonResult(uploadResult)

	} else {
		flag, _ := regexp.MatchString(`^.*\.tar\.gz$`, h.Filename)
		if !flag {
			uploadResult.Msg = "请上传正确的文件类型,如:demo.tar.gz"
			this.jsonResult(uploadResult)
		}
				
		fileTool := &libs.FileTool{Url: ""}
		uuidStr := fileTool.GenerateUuidStr()
		if uuidStr == "" {
			uploadResult.Msg = "文件保存出错，请重新选择文件"
			this.jsonResult(uploadResult)
		}
		
		uuidFileName := fmt.Sprintf("%s.tar.gz", uuidStr)
		filePath := fmt.Sprintf("%s/%s", models.TempDir , uuidFileName)
		os.MkdirAll(models.TempDir, 0777)
		this.SaveToFile("files[]", filePath)

		uploadResult.IsSuccess = true
		uploadResult.Data = &response.UploadFileInfo{
			OldFileName: h.Filename,
			NewFileName: uuidFileName,
		}
		this.jsonResult(uploadResult)
	}
}

//保存任务(此处还要处理当上传过程失败,在某一个时间段里不能删除上传的文件)
func (this *TaskController) SaveTask() {
	id, _ := this.GetInt("id", 0)
	isNew := true
	if id != 0 {
		isNew = false
	}

	task := new(models.Task)
	if !isNew {
		var err error
		task, err = models.TaskGetById(id)
		if err != nil {
			this.showMsg(err.Error()) //处理成ajax
		}
	} else {
		task.UserId = this.userId
	}

	task.TaskName = strings.TrimSpace(this.GetString("task_name"))
	task.Description = strings.TrimSpace(this.GetString("description"))
	task.GroupId, _ = this.GetInt("group_id")
	
	task.TaskType, _ = this.GetInt("task_type")
	task.ApiHeader = strings.TrimSpace(this.GetString("api_header"))
	task.ApiUrl = strings.TrimSpace(this.GetString("api_url"))
	task.ApiMethod = strings.TrimSpace(this.GetString("api_method"))
	task.PostBody = strings.TrimSpace(this.GetString("post_body"))
	task.OldGzipFile = strings.TrimSpace(this.GetString("old_gzip_file")) 
	
	task.Command = strings.TrimSpace(this.GetString("command"))
	
	task.Concurrent, _ = this.GetInt("concurrent")
	task.CronSpec = strings.TrimSpace(this.GetString("cron_spec"))
	task.Timeout, _ = this.GetInt("timeout")
	task.Notify, _ = this.GetInt("notify")
	task.NotifyEmail = strings.TrimSpace(this.GetString("notify_email"))
	new_temp_file := strings.TrimSpace(this.GetString("new_temp_file"))

	resultData := &response.ResultData{IsSuccess: false, Msg: ""}
	if task.TaskName == "" || task.CronSpec == "" || task.GroupId == 0  {
		resultData.Msg = "请填写完整信息,如下相关信息必填: 任务名称、cron表达式、分组"
		this.jsonResult(resultData)
	}
	
	//0:文件，1：API, 2:Shell脚本
	if (task.TaskType == 0 || task.TaskType == 2) && task.Command == "" {
			resultData.Msg = "请填入相关的运行命令"
			this.jsonResult(resultData)
	}
	
	if task.TaskType == 0 {
		if task.OldGzipFile == "" {
			resultData.Msg = "当选择类型为文件时,请上传要运行的文件"
			this.jsonResult(resultData)
		}
	} else if task.TaskType == 1 {
		if task.ApiUrl == "" || task.ApiMethod == "" {
			resultData.Msg = "当选择类型为API时,如下相关信息必填: API地址、调用方法"
			this.jsonResult(resultData)
		}
		if task.ApiMethod == "POST" && task.PostBody == "" {
			resultData.Msg = "当调用方式为POST时,请填写相关的POST Body"
			this.jsonResult(resultData)
		}
	} else if task.TaskType == 2{
		
	} else {
			resultData.Msg = "请正确选择类型"
			this.jsonResult(resultData)
	}
	
	if _, err := cron.Parse(task.CronSpec); err != nil {
		resultData.Msg = "cron表达式无效"
		this.jsonResult(resultData)
	}
	
	/*解压文件，生成文件夹*/
	if new_temp_file != "" {
		file := &libs.FileTool{Url: ""}
		tempfile := fmt.Sprintf("%s/%s", models.TempDir, new_temp_file)
		uuidStr := file.GenerateUuidStr()
		runDir := fmt.Sprintf("%s/%s", models.RunDir, uuidStr)
		runShell := fmt.Sprintf("%s/%s.sh", models.RunDir, uuidStr)
		
		if err := file.Ungzip(tempfile, runDir); err != nil {
			file.DeleteDir(runDir)
			resultData.Msg = "解压文件失败,文件打包是否正确"
			this.jsonResult(resultData)
		}
		
		//生成shell文件
		fileShell, errShell := os.OpenFile(runShell, os.O_RDWR|os.O_CREATE, 0766)
		if errShell != nil {
			beego.Info(errShell)
			resultData.Msg = "创建shell文件时出错"
			this.jsonResult(resultData)
		}
		fileShell.WriteString(
			fmt.Sprintf(
			`#!/bin/bash
			cd %s
		 	%s`, 
			fmt.Sprintf("%s/%s", models.RunDir, uuidStr), task.Command))
		
		fileShell.Close()
		os.Remove(tempfile)
		
		task.FileFolder = uuidStr
	}
	
	//保存数据库
	if isNew {
		if _, err := models.TaskAdd(task); err != nil {
			resultData.Msg = err.Error()
			this.jsonResult(resultData)
		}
	} else {
		if err := task.Update(); err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}
	}
	
	resultData.IsSuccess = true
	this.jsonResult(resultData)
}

// 任务执行日志列表
func (this *TaskController) Logs() {
	taskId, _ := this.GetInt("id")
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}

	task, err := models.TaskGetById(taskId)
	if err != nil {
		this.showMsg(err.Error())
	}

	result, count := models.TaskLogGetList(page, this.pageSize, "task_id", task.Id)

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["process_time"] = float64(v.ProcessTime) / 1000
		row["ouput_size"] = libs.SizeFormat(float64(len(v.Output)))
		row["status"] = v.Status
		list[k] = row
	}

	this.Data["pageTitle"] = "任务执行日志"
	this.Data["list"] = list
	this.Data["task"] = task
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("TaskController.Logs", "id", taskId), true).ToString()
	this.display()
}

// 查看日志详情
func (this *TaskController) ViewLog() {
	id, _ := this.GetInt("id")

	taskLog, err := models.TaskLogGetById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	task, err := models.TaskGetById(taskLog.TaskId)
	if err != nil {
		this.showMsg(err.Error())
	}

	data := make(map[string]interface{})
	data["id"] = taskLog.Id
	data["output"] = taskLog.Output
	data["error"] = taskLog.Error
	data["start_time"] = beego.Date(time.Unix(taskLog.CreateTime, 0), "Y-m-d H:i:s")
	data["process_time"] = float64(taskLog.ProcessTime) / 1000
	data["ouput_size"] = libs.SizeFormat(float64(len(taskLog.Output)))
	data["status"] = taskLog.Status

	this.Data["task"] = task
	this.Data["data"] = data
	this.Data["pageTitle"] = "查看日志"
	this.display()
}

// 批量操作日志
func (this *TaskController) LogBatch() {
	action := this.GetString("action")
	ids := this.GetStrings("ids")
	if len(ids) < 1 {
		this.ajaxMsg("请选择要操作的项目", MSG_ERR)
	}
	for _, v := range ids {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}
		switch action {
		case "delete":
			models.TaskLogDelById(id)
		}
	}

	this.ajaxMsg("", MSG_OK)
}

// 批量操作
func (this *TaskController) Batch() {
	action := this.GetString("action")
	ids := this.GetStrings("ids")
	if len(ids) < 1 {
		this.ajaxMsg("请选择要操作的项目", MSG_ERR)
	}

	for _, v := range ids {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}
		switch action {
		case "active":
			if task, err := models.TaskGetById(id); err == nil {
				job, err := jobs.NewJobFromTask(task)
				if err == nil {
					jobs.AddJob(task.CronSpec, job)
					task.Status = 1
					task.Update()
				}
			}
		case "pause":
			jobs.RemoveJob(id)
			if task, err := models.TaskGetById(id); err == nil {
				task.Status = 0
				task.Update()
			}
		case "delete":
			jobs.RemoveJob(id)
			this.deleteFileAndFolder(id)
			models.TaskDel(id)
			models.TaskLogDelByTaskId(id)
			
		}
	}

	this.ajaxMsg("", MSG_OK)
}

// 启动任务
func (this *TaskController) Start() {
	id, _ := this.GetInt("id")

	task, err := models.TaskGetById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	job, err := jobs.NewJobFromTask(task)
	if err != nil {
		this.showMsg(err.Error())
	}

	if jobs.AddJob(task.CronSpec, job) {
		task.Status = 1
		task.Update()
	}

	prevTimeStr := "-"
	if task.PrevTime > 0 {
		prevTimeStr = time.Unix(task.PrevTime, 0).Format(models.CNTimeFormat)
	}
	
	startJob := jobs.GetEntryById(id)
	this.Data["json"] = &response.ResultData{
		IsSuccess: true,
		Msg:       "",
		Data: &response.JobInfo{
			Status: 1,
			Prev:   prevTimeStr,
			Next:   beego.Date(startJob.Next, "Y-m-d H:i:s"),
		},
	}
	this.ServeJSON()
}

// 暂停任务
func (this *TaskController) Pause() {
	id, _ := this.GetInt("id")

	task, err := models.TaskGetById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	jobs.RemoveJob(id)
	task.Status = 0
	task.Update()

	this.Data["json"] = &response.ResultData{
		IsSuccess: true,
		Msg:       "",
		Data: &response.JobInfo{
			Status: 0,
			Prev:   time.Unix(task.PrevTime, 0).Format(models.CNTimeFormat),
			Next:   "-",
		},
	}
	this.ServeJSON()
}

// 立即执行
func (this *TaskController) Run() {
	id, _ := this.GetInt("id")

	entry := jobs.GetEntryById(id)
	if entry == nil {
		this.showMsg("没有找到相关的任务。")
	}
	entry.Job.Run()

	startJob := jobs.GetEntryById(id)
	task, _ := models.TaskGetById(id)

	prevTimeStr := "-"
	if task.PrevTime > 0 {
		prevTimeStr = time.Unix(task.PrevTime, 0).Format(models.CNTimeFormat)
	}
	
	this.Data["json"] = &response.ResultData{
		IsSuccess: true,
		Msg:       "",
		Data: &response.JobInfo{
			Status: 1,
			Prev:   prevTimeStr,
			Next:   beego.Date(startJob.Next, "Y-m-d H:i:s"),
		},
	}
	this.ServeJSON()
}

// 删除任务，同时删除数据库中的Task
func (this *TaskController) Delete() {
	id, _ := this.GetInt("id")
	jobs.RemoveJob(id)
	this.deleteFileAndFolder(id)
	models.TaskDel(id)
	models.TaskLogDelByTaskId(id)

	this.Data["json"] = &response.ResultData{
		IsSuccess: true,
		Msg:       "",
		Data:      true,
	}
	this.ServeJSON()
}

// 当删除任务后如果TaskType为0应该要删除相关的文件夹和shell文件
func (this *TaskController) deleteFileAndFolder(taskId int) error {
	task, err := models.TaskGetById(taskId)
	if err != nil {
		return err
	}
	if task.TaskType == 0 {
		taskDir := task.FileFolder
		if taskDir == "" {
			return nil
		}
		
		dataFileDir := fmt.Sprintf("%s/%s", models.RunDir, taskDir)
		shellFiel := fmt.Sprintf("%s.sh", dataFileDir)
		os.Remove(shellFiel)
		os.RemoveAll(dataFileDir)
	}
	return nil
}
