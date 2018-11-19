export default class TaskBulkUpdateRequest {
    public status: number
    public assign: number
    public deadline: string
    public bulkList: Array<string>

    public toJson(): string {
        return JSON.stringify(this)
    }
}