import {TaskStatus} from "@scripts/enums/TaskStatus"
import ProjectListResult from "@scripts/model/api/ProjectListResult"
import Project from "@scripts/model/api/project/Project"
import TaskListResult from "@scripts/model/api/TaskListResult"
import Task from "@scripts/model/api/task/Task"

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
                task.taskId = i["taskId"]
                task.projectId = i["projectId"]
                task.name = i["name"]
                task.creator = i["creator"]
                task.creatorName = i["creatorName"]
                task.assign = i["assign"]
                task.assignName = i["assignName"]
                task.status = i["status"]
                task.deadline = i["deadline"]
                task.limitDate = i["limitDate"]
                task.deadlineMD = i["deadlineMD"]
                task.description = i["description"]
                task.createDate = i["createDate"]

                taskList.task.push(task)
            }

            return taskList
        })
    }

    static getProject(): Promise<ProjectListResult> {
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
}