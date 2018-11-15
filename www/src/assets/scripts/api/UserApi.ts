import UserSearchRequest from "@scripts/model/api/UserSearchRequest"
import {UserSearchResult} from "@scripts/model/api/UserSearchResult"
import User from "@scripts/model/api/user/User"

export default class UserApi {
    static Search(searchRequest: UserSearchRequest): Promise<UserSearchResult> {
        const queryString = searchRequest.toQueryString()

        return fetch("/api/account/search" + queryString, {
            method: 'GET'
        }).then(res => {
            return res.json()
        }).then(json => {
            let res:UserSearchResult = []

            for(const _user of json){
                let user = new User()
                user.name = _user["name"]
                user.uuid = _user["uuid"]
                user.displayName = _user["displayName"]
                user.iconPath = _user["iconPath"]
                res.push(user)
            }

            return res
        })
    }
}