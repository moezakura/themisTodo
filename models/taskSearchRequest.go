package models

type TaskSearchRequest struct {
	TaskId       int
	ProjectId    int
	Status       TaskStatus
	AssignUserId int
	CreateUserId int
}

func (t TaskSearchRequest) ToSqlQueryAndArgs(tableName string) (queryText string, queryArgs []interface{}) {
	whereSQL := ""
	whereArray := make([]interface{}, 0)

	if t.Status != TASK_STATUS_OTHER {
		whereSQL += tableName + ".status = ? "
		whereArray = append(whereArray, t.Status)
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
		if len(whereSQL) > 0 {
			whereSQL += " AND "
		}
		whereSQL += tableName + ".assign = ? "
		whereArray = append(whereArray, t.AssignUserId)
	}

	if t.CreateUserId > 0 {
		if len(whereSQL) > 0 {
			whereSQL += " AND "
		}
		whereSQL += tableName + ".creator = ? "
		whereArray = append(whereArray, t.CreateUserId)
	}

	return whereSQL, whereArray
}
