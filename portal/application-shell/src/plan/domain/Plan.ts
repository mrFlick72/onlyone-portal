export type TodoStatus = "TODO" | "IN_PROGRESS" | "DONE" | "ABORTED";

export type Todo = {
    id: string;
    user_name: string;
    date: string;
    content: string;
    status: TodoStatus;
};

export type Plan = {
    id: string;
    user_name: string;
    title: string;
    date: string;
    todos: Todo[];
    todo_count: number;
};

export type NewPlan = {
    title: string;
    date: string;
};

export type TodoPayload = {
    date: string;
    content: string;
};

export default Plan;
