import {TaskStatus} from "@scripts/enums/TaskStatus"

export default class TaskHistoryItem {
    public name: string
    public editor: string
    public editorName: string
    public editorIconPath: string
    public status: TaskStatus
    public assign: number
    public assignName: string
    public assignIconPath: string
    public deadline: string
    public limitDate: number
    public deadlineMD: string
    public description: string
    public createDate: string
    public updateDate: string

    public fromAny(data: any) {
        this.name = data["name"]
        this.editor = data["name"]
        this.editorName = data["editor_name"]
        this.editorIconPath = data["editor_icon_path"]
        this.status = data["status"]
        this.assign = data["assign"]
        this.assignName = data["assign_name"]
        this.assignIconPath = data["assign_icon_path"]
        this.deadline = data["deadline"]
        this.limitDate = data["limit_date"]
        this.deadlineMD = data["deadline_md"]
        this.description = data["description"]
        this.createDate = data["create_date"]
        this.updateDate = data["update_date"]
    }
}