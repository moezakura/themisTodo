import {TaskStatus} from "@scripts/enums/TaskStatus"
import ProjectListResult from "@scripts/model/api/ProjectListResult"
import Project from "@scripts/model/api/project/Project"
import TaskListResult from "@scripts/model/api/TaskListResult"
import Task from "@scripts/model/api/task/Task"
import ProjectResult from "@scripts/model/api/ProjectResult"
import BaseApiResult from "@scripts/model/api/BaseApiResult"
import ProjectUpdateRequest from "@scripts/model/ProjectUpdateRequest"

export default class ProjectApi {
    static getList(taskStatus: TaskStatus): Promise<TaskListResult> {
        const status = taskStatus.toString()
        return fetch(`/api/tasks/my?status=${status}`, {
            method: 'GET',
            credentials: 'same-origin'
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

    static getProjects(): Promise<ProjectListResult> {
        return fetch("/api/project/my", {
            method: 'GET',
            credentials: 'same-origin'
        }).then(res => {
            return res.json()
        }).then(json => {
            let projectList = new ProjectListResult()

            projectList.message = json["message"]
            projectList.success = json["success"]

            projectList.project = []
            for (let i of <Array<any>>json["project"]) {
                let project = new Project()
                project.uuid = i["uuid"]
                project.name = i["name"]
                project.description = i["description"]

                projectList.project.push(project)
            }

            return projectList
        })
    }

    static deleteProject(projectId: number): Promise<BaseApiResult> {
        return fetch(`/api/project/delete/${projectId}`, {
            method: 'POST',
            body: "",
            credentials: "same-origin"
        }).then(res => {
            return res.json()
        }).then(json => {
            let res = new BaseApiResult()

            res.message = json["message"]
            res.success = json["success"]

            return res
        })
    }

    static updateProject(projectId: number, updateRequest: ProjectUpdateRequest): Promise<BaseApiResult> {
        return fetch(`/api/project/update/${projectId}`, {
            method: 'POST',
            body: updateRequest.toJson(),
            credentials: "same-origin"
        }).then(res => {
            return res.json();
        }).then(json => {
            let res = new BaseApiResult()

            res.message = json["message"]
            res.success = json["success"]

            return res
        });
    }

    static getProject(projectId: number): Promise<ProjectResult> {
        return fetch(`/api/project/info/${projectId}`, {
            method: 'GET',
            credentials: 'same-origin'
        }).then(res => {
            return res.json()
        }).then(json => {
            let project = new ProjectResult()

            project.message = json["message"]
            project.success = json["success"]

            project.project = new Project()
            project.project.uuid = json["uuid"]
            project.project.name = json["name"]
            project.project.description = json["description"]

            return project
        })
    }

    static getTasks(projectId: number): Promise<TaskListResult> {
        return fetch(`/api/project/tasks/${projectId}`, {
            method: 'GET',
            credentials: 'same-origin'
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
}