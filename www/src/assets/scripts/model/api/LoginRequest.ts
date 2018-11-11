import BaseJsonStruct from "../BaseJsonStruct"

export default class LoginRequest extends BaseJsonStruct {
    public id: string
    public password: string
}