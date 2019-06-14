export default class TaskTimerGetMyListRequest {
    public startDate: Date
    public endDate: Date

    public toQueryString(): string {
        let queryString = "?"

        const startDateText = (this.startDate.getFullYear() + "-")
            + (("0"+this.startDate.getMonth()).slice(-2) + "-")
            + (("0"+this.startDate.getDate()).slice(-2) + " ")
            + (("0"+this.startDate.getHours()).slice(-2) + ":")
            + (("0"+this.startDate.getMinutes()).slice(-2) + ":")
            + ("0"+this.startDate.getSeconds()).slice(-2)

        const endDateText = (this.endDate.getFullYear() + "-")
            + (("0"+this.endDate.getMonth()).slice(-2) + "-")
            + (("0"+this.endDate.getDate()).slice(-2) + " ")
            + (("0"+this.endDate.getHours()).slice(-2) + ":")
            + (("0"+this.endDate.getMinutes()).slice(-2) + ":")
            + ("0"+this.endDate.getSeconds()).slice(-2)

        queryString += "start=" + encodeURIComponent(startDateText) + "&end=" + encodeURIComponent(endDateText)

        return queryString
    }
}