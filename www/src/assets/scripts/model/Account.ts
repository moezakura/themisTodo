import BaseJsonStruct from "./BaseJsonStruct"

export default class Account extends BaseJsonStruct {
    public name: string | null = null
    public displayName: string | null = null
    public password: string | null = null
    public currentPassword: string | null = null
}