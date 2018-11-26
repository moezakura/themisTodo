import LoginRequest from '@scripts/model/api/LoginRequest'
import LoginResult from "@scripts/model/api/LoginResult"
import BaseApiResult from "@scripts/model/api/BaseApiResult"
import BaseApi from "@scripts/api/BaseApi"

export default class AuthApi extends BaseApi {
    static login(loginObject: LoginRequest): Promise<LoginResult> {
        return fetch("/api/login", {
            method: 'POST',
            body: loginObject.toJson(),
            credentials: "same-origin"
        }).then(response => {
            return response.json()
        }).then(json => {
            let result = new LoginResult()
            result.message = json["message"]
            result.success = json["success"]
            return result
        })
    }

    static auth(): Promise<BaseApiResult> {
        return fetch("/api/auth", {
            method: 'OPTIONS',
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(response => {
            return response.json()
        }).then(json => {
            let result = new LoginResult()
            result.message = json["message"]
            result.success = json["success"]
            return result
        })
    }
}