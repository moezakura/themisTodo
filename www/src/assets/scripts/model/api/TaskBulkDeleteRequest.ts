export default class TaskBulkDeleteRequest {
    public bulkList: Array<string>

    public toJson(): string {
        return JSON.stringify(this)
    }
}