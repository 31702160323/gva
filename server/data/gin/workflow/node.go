package data

import (
	"github.com/flipped-aurora/gva/global"
	model "github.com/flipped-aurora/gva/model/gin/workflow"
	"time"

	"gorm.io/gorm"
)

var Node = new(node)

type node struct{}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@description: workflow_nodes 表数据初始化
func (n *node) Init() error {
	var nodes = []model.WorkflowNode{
		{ID: "end1603681358043", CreatedAt: time.Now(), UpdatedAt: time.Now(), WorkflowProcessID: "leaveFlow", Clazz: "end", Label: global.I18n.T("LeaveFail"), Type: "end-node", Shape: "end-node", Description: "", View: "view/exa_wf_leave/exa_wf_leaveFrom.vue", X: 302, Y: 545.5, HideIcon: false, AssignType: "", AssignValue: "", Success: false},
		{ID: "end1603681360882", CreatedAt: time.Now(), UpdatedAt: time.Now(), WorkflowProcessID: "leaveFlow", Clazz: "end", Label: global.I18n.T("LeaveSuccess"), Type: "end-node", Shape: "end-node", Description: global.I18n.T("LeaveSuccessDesc"), View: "view/exa_wf_leave/exa_wf_leaveFrom.vue", X: 83.5, Y: 546, HideIcon: false, AssignType: "", AssignValue: "", Success: true},
		{ID: "start1603681292875", CreatedAt: time.Now(), UpdatedAt: time.Now(), WorkflowProcessID: "leaveFlow", Clazz: "start", Label: global.I18n.T("InitiateRequestLeave"), Type: "start-node", Shape: "start-node", Description: global.I18n.T("InitiateRequestLeaveDesc"), View: "view/exa_wf_leave/exa_wf_leaveFrom.vue", X: 201, Y: 109, HideIcon: false, AssignType: "", AssignValue: "", Success: false},
		{ID: "userTask1603681299962", CreatedAt: time.Now(), UpdatedAt: time.Now(), WorkflowProcessID: "leaveFlow", Clazz: "userTask", Label: global.I18n.T("ExaminationApproval"), Type: "user-task-node", Shape: "user-task-node", Description: global.I18n.T("ExaminationApprovalDesc"), View: "view/exa_wf_leave/exa_wf_leaveFrom.vue", X: 202, Y: 320.5, HideIcon: false, AssignType: "user", AssignValue: ",1,2,", Success: false},
	}
	return global.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&nodes).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		return nil
	})
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@description: 定义表名
func (n *node) TableName() string {
	return "workflow_nodes"
}
