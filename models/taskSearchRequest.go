package models

type TaskSearchRequest struct {
	TaskId       int
	ProjectId    int
	Status       TaskStatus
	AssignUserId int
	CreateUserId int
}

func (t TaskSearchRequest) ToSqlQueryAndArgs(joinTableName, tableName string) (joinText, queryText string, queryArgs []interface{}) {
	whereSQL := ""
	joinSQL := ""
	joinArray := make([]interface{}, 0)
	whereArray := make([]interface{}, 0)

	if t.Status != TASK_STATUS_OTHER {
		joinSQL += joinTableName + ".status = ? "
		joinArray = append(joinArray, t.Status)
	}

	if t.TaskId > 0 {
		if len(whereSQL) > 0 {
			whereSQL += " AND "
		}
		whereSQL += tableName + ".id = ? "
		whereArray = append(whereArray, t.TaskId)
	}

	if t.ProjectId > 0 {
		if len(whereSQL) > 0 {
			whereSQL += " AND "
		}
		whereSQL += tableName + ".project = ? "
		whereArray = append(whereArray, t.ProjectId)
	}

	if t.AssignUserId > 0 {
		if len(joinSQL) > 0 {
			joinSQL += " AND "
		}
		joinSQL += joinTableName + ".assign = ? "
		joinArray = append(joinArray, t.AssignUserId)
	}

	if t.CreateUserId > 0 {
		if len(whereSQL) > 0 {
			whereSQL += " AND "
		}
		whereSQL += tableName + ".creator = ? "
		whereArray = append(whereArray, t.CreateUserId)
	}
	whereArray = append(joinArray, whereArray...)

	return joinSQL, whereSQL, whereArray
}
