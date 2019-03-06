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

export default class TaskApi extends BaseApi {
    static getTaskFromCreateDate(createDate: string): Promise<TaskResult> {
        return fetch("/api/tasks/view/" + createDate, {
            method: 'GET',
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            let task = new TaskResult()

            task.success = json["success"]
            task.message = json["message"]
            task.task = new Task()
            task.task.fromAny(json["task"])

            return task
        })
    }

    static search(searchRequest: TaskSearchRequest): Promise<TaskListResult> {
        return fetch(`/api/tasks/search${searchRequest.toQueryString()}`, {
            method: 'GET',
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            let taskList = new TaskListResult()

            taskList.message = json["message"]
            taskList.success = json["success"]

            taskList.task = []
            for (let i of <Array<any>>json["tasks"]) {
                let task = new Task()
                task.fromAny(i)

                taskList.task.push(task)
            }

            return taskList
        })
    }

    static create(taskAddRequest: TaskAddRequest): Promise<TaskCreateResult> {
        return fetch("/api/tasks/create", {
            method: 'POST',
            body: taskAddRequest.toJson(),
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            let res = new TaskCreateResult()

            res.message = json["message"]
            res.success = json["success"]
            res.createDate = json["createDate"]

            return res
        })
    }

    static update(createDate: string, task: Task): Promise<BaseApiResult> {
        return fetch("/api/tasks/update/" + createDate, {
            method: 'POST',
            body: task.toJson(),
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            let res = new BaseApiResult()

            res.success = json["success"]
            res.message = json["message"]

            return res
        })
    }

    static bulkUpdate(updateRequest: TaskBulkUpdateRequest): Promise<BaseApiResult> {
        return fetch("/api/tasks/bulkUpdate", {
            method: 'POST',
            body: updateRequest.toJson(),
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            let res = new BaseApiResult()

            res.success = json["success"]
            res.message = json["message"]

            return res
        })
    }

    static delete(createDate: string): Promise<BaseApiResult> {
        return fetch(`/api/tasks/delete/${createDate}`, {
            method: 'POST',
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            let res = new BaseApiResult()

            res.success = json["success"]
            res.message = json["message"]

            return res
        })
    }

    static bulkDelete(deleteRequest: TaskBulkDeleteRequest): Promise<BaseApiResult> {
        return fetch(`/api/tasks/bulkDelete`, {
            method: 'DELETE',
            body: deleteRequest.toJson(),
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            let res = new BaseApiResult()

            res.success = json["success"]
            res.message = json["message"]

            return res
        })
    }

    static getHistory(createDate: string): Promise<Array<TaskHistory> | undefined> {
        return fetch(`/api/tasks/history/${createDate}`, {
            method: 'GET',
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            if (!json["success"]) {
                return undefined
            }

            let res: Array<TaskHistory> = []
            for (const th of json["payload"]) {
                const t = new TaskHistory()
                t.fromAny(th)
                res.push(t)
            }

            return res
        })
    }
}