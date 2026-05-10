import { getPlanApiBaseUrl } from "../../config/ConfigLoader";
import Plan, { NewPlan, TodoPayload, TodoStatus } from "./Plan";

const PLAN_URI = (baseUrl: string, planId?: string) => planId ?
    `${baseUrl}/api/plan/${planId}` :
    `${baseUrl}/api/plan`;

const TODO_URI = (baseUrl: string, planId: string, todoId?: string) => todoId ?
    `${baseUrl}/api/plan/${planId}/todo/${todoId}` :
    `${baseUrl}/api/plan/${planId}/todo`;

const authHeaders = (extra: Record<string, string> = {}) => ({
    'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
    'Accept': 'application/json',
    ...extra,
});

export async function getAllPlans(): Promise<Plan[]> {
    const baseUrl = await getPlanApiBaseUrl();
    const response = await fetch(PLAN_URI(baseUrl), {
        method: "GET",
        credentials: 'include',
        headers: authHeaders(),
    });
    if (!response.ok) {
        return [];
    }
    return response.json();
}

export async function getPlan(planId: string): Promise<Plan | null> {
    const baseUrl = await getPlanApiBaseUrl();
    const response = await fetch(PLAN_URI(baseUrl, planId), {
        method: "GET",
        credentials: 'include',
        headers: authHeaders(),
    });
    if (!response.ok) {
        return null;
    }
    return response.json();
}

export async function createPlan(plan: NewPlan): Promise<{ id: string } | null> {
    const baseUrl = await getPlanApiBaseUrl();
    const response = await fetch(PLAN_URI(baseUrl), {
        method: "POST",
        credentials: 'include',
        headers: authHeaders({ 'Content-Type': 'application/json' }),
        body: JSON.stringify(plan),
    });
    if (!response.ok) {
        return null;
    }
    return response.json();
}

export async function deletePlan(planId: string) {
    const baseUrl = await getPlanApiBaseUrl();
    return fetch(PLAN_URI(baseUrl, planId), {
        method: "DELETE",
        credentials: 'include',
        headers: authHeaders(),
    });
}

export async function addTodo(planId: string, todo: TodoPayload) {
    const baseUrl = await getPlanApiBaseUrl();
    return fetch(TODO_URI(baseUrl, planId), {
        method: "PUT",
        credentials: 'include',
        headers: authHeaders({ 'Content-Type': 'application/json' }),
        body: JSON.stringify(todo),
    });
}

export async function updateTodo(planId: string, todoId: string, todo: TodoPayload) {
    const baseUrl = await getPlanApiBaseUrl();
    return fetch(TODO_URI(baseUrl, planId, todoId), {
        method: "PUT",
        credentials: 'include',
        headers: authHeaders({ 'Content-Type': 'application/json' }),
        body: JSON.stringify(todo),
    });
}

export async function removeTodo(planId: string, todoId: string) {
    const baseUrl = await getPlanApiBaseUrl();
    return fetch(TODO_URI(baseUrl, planId, todoId), {
        method: "DELETE",
        credentials: 'include',
        headers: authHeaders(),
    });
}

export async function changeTodoStatus(planId: string, todoId: string, status: TodoStatus) {
    const baseUrl = await getPlanApiBaseUrl();
    return fetch(`${TODO_URI(baseUrl, planId, todoId)}/status`, {
        method: "PUT",
        credentials: 'include',
        headers: authHeaders({ 'Content-Type': 'application/json' }),
        body: JSON.stringify({ status }),
    });
}
