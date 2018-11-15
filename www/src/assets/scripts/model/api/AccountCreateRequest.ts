export default class AccountCreateRequest {
    public name: string

    public toJson(): string {
        return JSON.stringify(this)
    }
}