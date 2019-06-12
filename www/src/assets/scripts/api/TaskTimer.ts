import Task from "@scripts/model/api/task/Task"
import TaskListResult from "@scripts/model/api/TaskListResult"
import TaskResult from "@scripts/model/api/TaskResult"
import BaseApiResult from "@scripts/model/api/BaseApiResult"
import TaskCreateResult from "@scripts/model/api/TaskCreateResult"
import TaskAddRequest from "@scripts/model/api/TaskAddRequest"
import TaskSearchRequest from "@scripts/model/api/TaskSearchRequest"
import TaskBulkUpdateRequest from "@scripts/model/api/TaskBulkUpdateRequest"
import TaskBulkDeleteRequest from "@scripts/model/api/TaskBulkDeleteRequest"
import BaseApi from "@scripts/api/BaseApi"
import TaskHistory from "@scripts/model/api/task/TaskHistory"
import TaskTimerToggleResult from "@scripts/model/api/taskTimer/TaskTimerToggleResult"

export default class TaskTimerApi extends BaseApi {
    static toggleTimer(createDate: string): Promise<TaskTimerToggleResult> {
        return fetch("/api/tasks/timer/toggle/" + createDate, {
            method: 'PATCH',
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            let task = new TaskTimerToggleResult()

            task.success = json["success"]
            task.message = json["message"]
            task.start = json["start"]

            return task
        })
    }

}