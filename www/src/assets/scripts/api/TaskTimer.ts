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
import TaskTimerGetResult from "@scripts/model/api/taskTimer/TaskTimerGetResult"
import TaskTimerGetMyListRequest from "@scripts/model/api/taskTimer/TaskTimerGetMyListRequest"
import TaskTimerGetMyListResult from "@scripts/model/api/taskTimer/TaskTimerGetMyListResult"
import TaskTimer from "@scripts/model/api/taskTimer/TaskTimer"
import TaskTimerGetStatusResult from "@scripts/model/api/taskTimer/TaskTimerGetStatusResult"
import TaskTimerUpdateRequest from "@scripts/model/api/taskTimer/TaskTimerUpdateRequest"
import TaskTimerGetMyDoingResult from "@scripts/model/api/taskTimer/TaskTimerGetMyDoingListResult"

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

    static getTaskTimer(createDate: string): Promise<TaskTimerGetResult> {
        return fetch("/api/tasks/timer/view/" + createDate, {
            method: 'GET',
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            let task = new TaskTimerGetResult()

            task.success = json["success"]
            task.message = json["message"]
            task.start = json["start"]
            task.LastStartTime = json["last_start_time"]
            task.LastEndTime = json["last_end_time"]
            task.TotalTime = json["total_time"]
            task.TodayTime = json["today_time"]

            return task
        })
    }

    static deleteTaskTimer(id: number): Promise<BaseApiResult> {
        return fetch("/api/tasks/timer/delete/" + id, {
            method: 'DELETE',
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            let task = new BaseApiResult()

            task.success = json["success"]
            task.message = json["message"]

            return task
        })
    }

    static updateTaskTimer(id: number, updateRequest: TaskTimerUpdateRequest): Promise<BaseApiResult> {
        return fetch("/api/tasks/timer/update/" + id, {
            method: 'POST',
            body: updateRequest.toJson(),
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            let task = new BaseApiResult()

            task.success = json["success"]
            task.message = json["message"]

            return task
        })
    }

    static getTaskTimerStatus(createDate: string): Promise<TaskTimerGetStatusResult> {
        return fetch("/api/tasks/timer/status/" + createDate, {
            method: 'GET',
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            let task = new TaskTimerGetStatusResult()

            task.success = json["success"]
            task.message = json["message"]
            task.start = json["start"]

            return task
        })
    }

    static getMyList(projectId: number, searchRequest: TaskTimerGetMyListRequest): Promise<TaskTimerGetMyListResult> {
        return fetch(`/api/tasks/timer/myList/${projectId}${searchRequest.toQueryString()}`, {
            method: 'GET',
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            let res = new TaskTimerGetMyListResult()

            res.success = json["success"]
            res.message = json["message"]
            res.list = new Array<TaskTimer>();
            for (const item of json["list"]) {
                const t = new TaskTimer()
                t.fromAny(item)
                res.list.push(t)
            }

            return res
        })
    }

    static getDoingTimerByProjectId(projectId: number){
        return fetch(`/api/tasks/timer/project/${projectId}`, {
            method: 'GET',
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            let res = new TaskTimerGetMyListResult()

            res.success = json["success"]
            res.message = json["message"]
            res.list = new Array<TaskTimer>();
            for (const item of json["list"]) {
                const t = new TaskTimer()
                t.fromAny(item)
                res.list.push(t)
            }

            return res
        })
    }

    static getMyDoingList(): Promise<TaskTimerGetMyDoingResult> {
        return fetch(`/api/tasks/timer/myDoing`, {
            method: 'GET',
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            let res = new TaskTimerGetMyDoingResult()

            res.success = json["success"]
            res.message = json["message"]
            res.list = new Array<TaskTimer>();
            for (const item of json["list"]) {
                const t = new TaskTimer()
                t.fromAny(item)
                res.list.push(t)
            }

            return res
        })
    }
}