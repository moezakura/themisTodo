export default class Task {
    public taskId: number
    public projectId: number
    public name: string
    public creator: number
    public creatorName: string
    public creatorIconPath: string
    public assign: number
    public assignName: string
    public assignIconPath: string
    public status: number
    public deadline: string
    public limitDate: number
    public deadlineMD: string
    public description: string
    public createDate: string
    public updateDate: string
    public isDoing: boolean

    public fromAny(data: any) {
        this.taskId = data["taskId"]
        this.projectId = data["projectId"]
        this.name = data["name"]
        this.creator = data["creator"]
        this.creatorName = data["creatorName"]
        this.creatorIconPath = data["creatorIconPath"]
        this.assign = data["assign"]
        this.assignName = data["assignName"]
        this.assignIconPath = data["assignIconPath"]
        this.status = data["status"]
        this.deadline = data["deadline"]
        this.limitDate = data["limitDate"]
        this.deadlineMD = data["deadlineMD"]
        this.description = data["description"]
        this.createDate = data["createDate"]
        this.updateDate = data["updateDate"]
        this.isDoing = data["is_doing"]
    }

    public toJson(): string {
        return JSON.stringify(this)
    }
}