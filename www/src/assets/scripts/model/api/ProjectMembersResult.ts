import BaseApiResult from "@scripts/model/api/BaseApiResult"
import User from "@scripts/model/api/user/User"

export default class ProjectMembersResult extends BaseApiResult {
    public members: Array<User>
}