package models

type TaskBulkUpdateRequestJson struct {
	Status   int      `json:"status" binding:"required"`
	Assign   int      `json:"assign" binding:"required"`
	Deadline string   `json:"deadline" binding:"required"`
	BulkList []string `json:"bulkList" binding:"required"`
}
