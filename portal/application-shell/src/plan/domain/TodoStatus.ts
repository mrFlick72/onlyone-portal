import { TodoStatus } from "./Plan";

const transitions: Record<TodoStatus, TodoStatus[]> = {
    TODO:        ["IN_PROGRESS", "ABORTED"],
    IN_PROGRESS: ["DONE", "ABORTED"],
    DONE:        [],
    ABORTED:     [],
};

export const allowedTransitionsFor = (status: TodoStatus): TodoStatus[] =>
    transitions[status] ?? [];

export const isTerminal = (status: TodoStatus): boolean =>
    allowedTransitionsFor(status).length === 0;

const labels: Record<TodoStatus, string> = {
    TODO:        "To Do",
    IN_PROGRESS: "In Progress",
    DONE:        "Done",
    ABORTED:     "Aborted",
};

export const labelFor = (status: TodoStatus): string => labels[status] ?? status;

type ChipColor = "default" | "primary" | "success" | "warning";
const colors: Record<TodoStatus, ChipColor> = {
    TODO:        "default",
    IN_PROGRESS: "primary",
    DONE:        "success",
    ABORTED:     "warning",
};

export const colorFor = (status: TodoStatus): ChipColor => colors[status] ?? "default";
