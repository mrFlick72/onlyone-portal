type Month = {
    monthLabel: string
    monthValue: number
}
type Months = Month[]


const MONTHS =
    [
        {
            monthLabel: "January",
            monthValue: 1
        },
        {
            monthLabel: "February",
            monthValue: 2
        },
        {
            monthLabel: "March",
            monthValue: 3
        },
        {
            monthLabel: "April",
            monthValue: 4
        },
        {
            monthLabel: "May",
            monthValue: 5
        },
        {
            monthLabel: "June",
            monthValue: 6
        },
        {
            monthLabel: "July",
            monthValue: 7
        },
        {
            monthLabel: "August",
            monthValue: 8
        },
        {
            monthLabel: "September",
            monthValue: 9
        },
        {
            monthLabel: "October",
            monthValue: 10
        },
        {
            monthLabel: "November",
            monthValue: 11
        },
        {
            monthLabel: "December",
            monthValue: 12
        }
    ] as Months

export default MONTHS
export type { Months, Month }