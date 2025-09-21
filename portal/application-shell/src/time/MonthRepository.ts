import MONTHS, { Months } from "./months";

export async function getMonthRegistry(): Promise<Months> {
    return MONTHS
}