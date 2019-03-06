import TaskHistoryItem from "@scripts/model/api/task/TaskHistoryItem"

export default class TaskHistory {
    public createDate: string
    public updateDate: string
    public task: TaskHistoryItem
    public updateDateFormat: string


    public fromAny(data: any) {
        this.createDate = data["create_date"]
        this.updateDate = data["update_date"]
        this.task = new TaskHistoryItem()
        this.task.fromAny(data["task"])
    }

    public toJson(): string {
        return JSON.stringify(this)
    }
}
