class ProjectUtils {
    static createUserListLine(uuid, name, displayName) {
        let parentLi = document.createElement("li");
        parentLi.dataset.uuid = uuid;
        parentLi.dataset.name = name;
        parentLi.dataset.displayName = displayName;
        parentLi.id = "searchResult" + uuid;

        let iconDiv = document.createElement("div");
        let nowTime = new Date();
        iconDiv.classList.add("icon");
        iconDiv.style.backgroundImage = "url(\"/account/icon/" + uuid + "?t=" + nowTime.getTime() + "\")";
        parentLi.appendChild(iconDiv);

        let nameLineDiv = document.createElement("div");
        nameLineDiv.classList.add("name");

        {
            let nameDiv = document.createElement("div");
            nameDiv.classList.add("nameId");
            nameDiv.innerText = name;

            let displayNameDiv = document.createElement("div");
            displayNameDiv.classList.add("displayName");
            displayNameDiv.innerText = displayName;

            nameLineDiv.appendChild(nameDiv);
            nameLineDiv.appendChild(displayNameDiv);
        }

        parentLi.appendChild(nameLineDiv);

        return parentLi;
    }
}