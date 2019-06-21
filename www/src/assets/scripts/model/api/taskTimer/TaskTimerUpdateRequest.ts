export default class TaskTimerUpdateRequest {
    public startDateHMS: string
    public endDateHMS: string
    public note: string

    public toJson(): string {
        return JSON.stringify({
            start_date_hms: this.startDateHMS,
            end_date_hms: this.endDateHMS,
            note: this.note
        })
    }
}