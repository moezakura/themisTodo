export default class AddMemberRequest {
    public uuid: number

    public toJson(): string {
        return JSON.stringify(this)
    }
}