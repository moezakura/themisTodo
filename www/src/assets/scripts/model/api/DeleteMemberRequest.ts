export default class DeleteMemberRequest {
    public uuid: number

    public toJson(): string {
        return JSON.stringify(this)
    }
}