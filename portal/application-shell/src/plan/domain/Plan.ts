export type Todo = {
    id: string;
    user_name: string;
    date: string;
    content: string;
};

export type Plan = {
    id: string;
    user_name: string;
    title: string;
    date: string;
    todos: Todo[];
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
