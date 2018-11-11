import BaseApiResult from "@scripts/model/api/BaseApiResult"
import Project from "@scripts/model/api/project/Project"

export default class ProjectListResult extends BaseApiResult {
    public project: Array<Project>
}