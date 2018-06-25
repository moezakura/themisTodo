import ProjectUtils from "./projectUtils"

export default class UserSearchDialog {
    constructor(target, config) {
        this.selectUserIndex = -1;
        this.selectMode = false;
        this.oldSearchInput = "";
        this.searchResult = [];
        this.target = target;
        this.userSelectUserList = UserSearchDialog.createSearchBox();
        this.selectedUserId = -1;
        this.sendEvent = config.forceSubmit;
        this.userSelectEvent = config.userSelect;
        this.singleEnter = config.singleEnter === true;
        this.isIn = config.isIn === true;
        let that = this;

        target.addEventListener("keyup", function (e) {
            that.userSelectInputKeyUp(e, that);
        }, true);
        target.addEventListener("keydown", function (e) {
            that.userSelectInputKeyDown(e, that);
        }, true);
    }

    static createSearchBox() {
        let searchBox = document.createElement("ul");
        searchBox.id = "usersSearchedList-" + (new Date()).getTime();
        searchBox.classList.add("usersList");
        searchBox.classList.add("userSearchDialog");
        document.body.appendChild(searchBox);

        return searchBox;
    }

    get() {
        return getUuid(this);
    }

    hide() {
        this.userSelectUserList.style.display = "none";
    }

    userSelectInputKeyUp(e, that) {
        let inputText = that.target.value;
        if (inputText.length < 1) {
            that.userSelectUserList.style.display = "none";
            that.selectUserIndex = -1;
            return;
        }

        // up, down, escape, enter key down
        if (e.keyCode === 38 || e.keyCode === 40 | e.keyCode === 27 || e.keyCode === 13)
            return;

        if (inputText !== that.oldSearchInput) {
            that.selectUserIndex = -1;
        }

        let queryObject = {
            "name": inputText,
            "displayName": inputText,
            "project": projectId,
            "isInProject": that.isIn,
            "max": 20
        };

        let queryString = "?";
        for (let key in queryObject)
            queryString += key + "=" + encodeURIComponent(queryObject[key]) + "&";
        queryString = queryString.slice(0, -1);

        fetch("/account/search" + queryString, {
            method: 'GET'
        }).then(function (response) {
            return response.json();
        }).then(function (json) {
            that.searchResult = json;
            let targetPos = that.target.getBoundingClientRect();
            that.userSelectUserList.style.display = "block";
            that.userSelectUserList.style.left = (targetPos.left + window.pageXOffset) + "px";
            that.userSelectUserList.style.top = (targetPos.top + targetPos.height + window.pageYOffset) + "px";
            that.userSelectUserList.style.width = targetPos.width + "px";
            that.userSelectUserList.scrollTop = 0;
            that.selectMode = false;

            that.userSelectUserList.innerHTML = "";
            json.forEach(function (value) {
                let elem = ProjectUtils.createUserListLine(value.uuid, value.name, value.displayName);
                elem.addEventListener("click", function (e) {
                    that.userClick(e, this);
                });
                that.userSelectUserList.appendChild(elem);
            });
        });
    }

    userSelectInputKeyDown(e, that) {
        if (e.keyCode === 38) {
            //up key
            that.selectUserIndex--;
            if (that.selectUserIndex < 0) that.selectUserIndex = that.searchResult.length - 1;
        } else if (e.keyCode === 40) {
            //down key
            that.selectUserIndex++;
            if (that.selectUserIndex > that.searchResult.length - 1) that.selectUserIndex = 0;
        } else if (e.keyCode === 27) {
            //escape key
            that.userSelectUserList.style.display = "none";
            return;
        } else if (e.keyCode === 13) {
            //enter key
            this.projectMemberAddSubmit(e, that);
            return;
        } else return;

        if (that.selectUserIndex < 0 || that.selectUserIndex < 0 || that.selectUserIndex > that.searchResult.length - 1) return;
        that.selectMode = true;

        let all = that.userSelectUserList.querySelectorAll("li");
        all.forEach(function (value) {
            value.classList.remove("select");
        });

        let targetId = "#searchResult" + that.searchResult[that.selectUserIndex].uuid;
        let target = that.userSelectUserList.querySelector(targetId);
        let targetPos = target.getBoundingClientRect();
        let userSelectUserListPos = that.userSelectUserList.getBoundingClientRect();
        that.userSelectUserList.scrollTop = targetPos.y + that.userSelectUserList.scrollTop -
            userSelectUserListPos.height + targetPos.height - userSelectUserListPos.y;
        target.classList.add("select");
    }

    projectMemberAddSubmit(e, that) {
        e.preventDefault();
        if (that.userSelectEvent !== undefined || that.userSelectEvent != null)
            that.userSelectEvent(that.getUuid(that));

        this.selectedUserId = that.getUuid(that);

        if (!that.selectMode || this.singleEnter) {
            let sendUuid = that.getUuid(that);
            that.sendEvent(sendUuid);
            that.userSelectUserList.style.display = "none";
            that.setName(that);
        }
        else if (that.selectMode) {
            that.setName(that);
            that.userSelectUserList.style.display = "none";
            that.selectMode = false;
        }
    }

    userClick(e, _this) {
        let elements = [].slice.call(this.userSelectUserList.querySelectorAll("li"));
        this.selectUserIndex = elements.indexOf(_this);
        this.selectMode = true;
        this.projectMemberAddSubmit(e, this);
    }

    submit(){
        let sendUuid = this.getUuid(this);
        this.sendEvent(sendUuid);
    }

    setName(that) {
        if (!(that.selectUserIndex < 0 || that.selectUserIndex < 0 || that.selectUserIndex > that.searchResult.length - 1))
            that.target.value = that.searchResult[that.selectUserIndex].name;
    }

    getUuid(that) {
        let sendUuid = -1;
        if (that.selectUserIndex < 0 || that.selectUserIndex < 0 || that.selectUserIndex > that.searchResult.length - 1) {

        } else sendUuid = that.searchResult[that.selectUserIndex].uuid;

        return sendUuid;
    }

}