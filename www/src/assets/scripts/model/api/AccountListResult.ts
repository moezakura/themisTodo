import BaseApiResult from "@scripts/model/api/BaseApiResult"
import User from "@scripts/model/api/user/User"

export default class AccountListResult extends BaseApiResult {
    public users: Array<User>
}