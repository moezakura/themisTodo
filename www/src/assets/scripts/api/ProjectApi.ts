import {TaskStatus} from "@scripts/enums/TaskStatus"
import ProjectListResult from "@scripts/model/api/ProjectListResult"
import Project from "@scripts/model/api/project/Project"

export default class ProjectApi {
    static getList(taskStatus: TaskStatus) {
        const status = taskStatus.toString()
        return fetch(`/api/tasks/my?status=${status}`, {
            method: 'GET',
            credentials: 'same-origin'
        }).then(res => {
            return res.json()
        }).then(json => {

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