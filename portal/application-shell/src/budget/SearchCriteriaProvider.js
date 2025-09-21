export function getMonthSearchCriteria() {
    let month = window.sessionStorage.getItem("month");
    if (!month) {
        month = String(new Date().getMonth() + 1);
        window.sessionStorage.setItem("month", month)
    }
    return month
}

export function setMonthSearchCriteria(month) {
    window.sessionStorage.setItem("month", month)
}

export function getYearSearchCriteria() {
    let year = window.sessionStorage.getItem("year");
    if (!year) {
        year = String(new Date().getFullYear());
        window.sessionStorage.setItem("year", year)
    }
    return year
}

export function setYearSearchCriteria(year) {
    window.sessionStorage.setItem("year", year)
}

export function getSearchTagsSearchCriteria() {
    let searchTags = window.sessionStorage.getItem("searchTags");
    if (!searchTags) {
        searchTags = ""
        window.sessionStorage.setItem("searchTags", searchTags)
    }
    return searchTags.split(",")
}

export function setSearchTagsSearchCriteria(searchTags) {
    window.sessionStorage.setItem("searchTags", searchTags)
}

export function resetSearchParameters() {
    window.sessionStorage.removeItem("month")
    window.sessionStorage.removeItem("year")
    window.sessionStorage.removeItem("searchTags")

    getMonthSearchCriteria()
    getYearSearchCriteria()
    getSearchTagsSearchCriteria()
}