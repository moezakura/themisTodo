import Task from "@scripts/model/api/task/Task"
import TaskListResult from "@scripts/model/api/TaskListResult"
import TaskResult from "@scripts/model/api/TaskResult"
import BaseApiResult from "@scripts/model/api/BaseApiResult"
import TaskCreateResult from "@scripts/model/api/TaskCreateResult"
import TaskAddRequest from "@scripts/model/api/TaskAddRequest"

export default class TaskApi {
    static GetTaskFromCreateDate(createDate: string): Promise<TaskResult> {
        return fetch("/api/tasks/view/" + createDate, {
            method: 'GET',
            credentials: "same-origin"
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

    static GetSearch(taskId: number, projectId: number): Promise<TaskListResult> {
        return fetch(`/api/tasks/search?taskId=${taskId}&projectId=${projectId}`, {
            method: 'GET',
            credentials: "same-origin"
        }).then(res => {
            return res.json()
        }).then(json => {
            let taskList = new TaskListResult()

            taskList.message = json["message"]
            taskList.success = json["success"]

            taskList.task = []
            for (let i of <Array<any>>json["task"]) {
                let task = new Task()
                task.fromAny(i)

                taskList.task.push(task)
            }

            return taskList
        })
    }

    static Create(taskAddRequest: TaskAddRequest): Promise<TaskCreateResult> {
        return fetch("/api/tasks/create", {
            method: 'POST',
            body: taskAddRequest.toJson(),
            credentials: "same-origin"
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

    static Update(createDate: string, task: Task): Promise<BaseApiResult> {
        return fetch("/api/tasks/update/" + createDate, {
            method: 'POST',
            body: task.toJson(),
            credentials: "same-origin"
        }).then(res => {
            return res.json()
        }).then(json => {
            let res = new BaseApiResult()

            res.success = json["success"]
            res.message = json["message"]

            return res
        })
    }

    static Delete(createDate: string): Promise<BaseApiResult> {
        return fetch("/api/tasks/delete/" + createDate, {
            method: 'POST',
            credentials: "same-origin"
        }).then(res => {
            return res.json()
        }).then(json => {
            let res = new BaseApiResult()

            res.success = json["success"]
            res.message = json["message"]

            return res
        })
    }
}