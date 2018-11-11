import Task from "@scripts/model/api/task/Task"
import TaskListResult from "@scripts/model/api/TaskListResult"
import TaskResult from "@scripts/model/api/TaskResult"
import BaseApiResult from "@scripts/model/api/BaseApiResult"
import TaskCreateResult from "@scripts/model/api/TaskCreateResult"

export default class TaskApi {
    static GetTaskFromCreateDate(createDate: string): Promise<TaskResult> {
        return fetch("/api/tasks/view/" + createDate, {
            method: 'GET',
            credentials: "same-origin"
        }).then(function (response) {
            return response.json()
        })
    }

    static GetSearch(taskId, projectId): Promise<TaskListResult> {
        return fetch(`/api/tasks/search?taskId=${taskId}&projectId=${projectId}`, {
            method: 'GET',
            credentials: "same-origin"
        }).then(function (response) {
            return response.json()
        })
    }

    static Create(taskJson): Promise<TaskCreateResult> {
        return fetch("/api/tasks/create", {
            method: 'POST',
            body: JSON.stringify(taskJson),
            credentials: "same-origin"
        }).then(function (response) {
            return response.json()
        })
    }

    static Update(createDate, taskApi): Promise<BaseApiResult> {
        return fetch("/api/tasks/update/" + createDate, {
            method: 'POST',
            body: JSON.stringify(taskApi),
            credentials: "same-origin"
        }).then(function (response) {
            return response.json()
        })
    }

    static Delete(createDate): Promise<BaseApiResult> {
        return fetch("/api/tasks/delete/" + createDate, {
            method: 'POST',
            credentials: "same-origin"
        }).then(function (response) {
            return response.json()
        })
    }
}