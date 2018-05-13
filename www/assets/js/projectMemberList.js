memberUpdate();

function memberUpdate() {
    let userList = document.querySelector("#projectMemberAddForm .usersList");
    userList.innerHTML = "";
    memberList.forEach(function (value) {
        let elem = ProjectUtils.createUserListLine(value.uuid, value.name, value.displayName);
        userList.appendChild(elem);
    });
}