import ProjectUtils from "./projectUtils"

export default class ProjectMemberList {
    static update() {
        let userList = document.querySelector("#projectMemberAddForm .usersList");
        if(userList === undefined || userList == null) return;

        userList.innerHTML = "";
        memberList.forEach(function (value) {
            let elem = ProjectUtils.createUserListLine(value.uuid, value.name, value.displayName);
            userList.appendChild(elem);
        });
    }
}