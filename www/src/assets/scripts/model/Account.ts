export default class Account {
    public name: string | null = null
    public displayName: string | null = null
    public password: string | null = null
    public currentPassword: string | null = null

    public toJson(): string {
        const obj = {
            name: this.name,
            displayName: this.displayName,
            password: this.password,
            currentPassword: this.currentPassword,
        }

        return JSON.stringify(obj)
    }
}