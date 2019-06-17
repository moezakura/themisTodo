import Task from "@scripts/model/api/task/Task"

export default class TaskTimer {
    public id: number
    public createDate: string
    public assign: number
    public startDate: Date
    public startDateUnix: number
    public endDate: Date
    public endDateUnix: number
    public note: string
    public task: Task

    public startDateHM: string
    public endDateHM: string
    public totalSec: number
    public totalHM: string

    public fromAny(data: any) {
        this.id = data["id"]
        this.createDate = data["create_date"]
        this.assign = data["assign"]
        this.startDate = new Date(data["start_date_unix"] * 1000)
        this.startDateUnix = data["start_date_unix"]
        this.endDateUnix = data["end_date_unix"]
        if (this.endDateUnix > 0) {
            this.endDate = new Date(data["end_date_unix"] * 1000)
        } else {
            this.endDate = new Date()
        }
        this.note = data["note"]

        this.startDateHM = ("0" + this.startDate.getHours()).slice(-2) + ":" + ("0" + this.startDate.getMinutes()).slice(-2)
        this.endDateHM = ("0" + this.endDate.getHours()).slice(-2) + ":" + ("0" + this.endDate.getMinutes()).slice(-2)

        this.totalSec = Math.floor((this.endDate.getTime() - this.startDate.getTime()) / 1000)
        const totalHour = Math.floor(this.totalSec / 3600)
        const totalMin = Math.floor(this.totalSec / 60) - totalHour * 60
        this.totalHM = ("0" + totalHour).slice(-2) + ":" + ("0" + totalMin).slice(-2)

        let task = data["task"]
        if (typeof task !== "undefined" || task !== null) {
            this.task = new Task()
            this.task.fromAny(task)
        }
    }

    public toJson(): string {
        return JSON.stringify(this)
    }
}