import { ScheduledExpenseStatus } from "./ScheduledExpense";

type ChipColor = "success" | "warning" | "default";

const colors: Record<ScheduledExpenseStatus, ChipColor> = {
    ACTIVE: "success",
    PAUSED: "warning",
};

export const colorFor = (status?: ScheduledExpenseStatus): ChipColor =>
    status ? (colors[status] ?? "default") : "default";
