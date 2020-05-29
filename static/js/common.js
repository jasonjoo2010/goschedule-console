/**
 * Format given timestamp in milliseconds as 'yyyy-MM-dd HH:mm:ss'
 */
function formatMillis(millis) {
    var created = new Date(millis)
    var month = (created.getMonth() + 1) + "";
    var day = created.getDate() + "";
    var hour = created.getHours() + "";
    var minute = created.getMinutes() + "";
    var second = created.getSeconds() + "";
    return created.getFullYear() + "-" + month.padStart(2, "0") + "-" + day.padStart(2, "0") + " " + hour.padStart(2, "0") + ":" + minute.padStart(2, "0") + ":" + second.padStart(2, "0");
}