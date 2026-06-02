import React from "react";
import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import { Todo, TodoStatus } from "./domain/Plan";
import TodoRow from "./TodoRow";

interface TodoContentMessages {
    headers: { date: string; content: string; status: string; options: string };
    actions: { changeStatus: string; edit: string; delete: string };
    status: Record<TodoStatus, string>;
}

interface TodoContentProps {
    todos: Todo[];
    openUpdate: (todo: Todo) => void;
    openDelete: (todo: Todo) => void;
    openChangeStatus: (todo: Todo) => void;
    messages: TodoContentMessages;
}

const TodoContent: React.FC<TodoContentProps> = ({ todos, openUpdate, openDelete, openChangeStatus, messages }) => (
    <TableContainer component={Paper}>
        <Table>
            <TableHead>
                <TableRow>
                    <TableCell>{messages.headers.date}</TableCell>
                    <TableCell>{messages.headers.content}</TableCell>
                    <TableCell>{messages.headers.status}</TableCell>
                    <TableCell>{messages.headers.options}</TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {todos.map(todo =>
                    <TodoRow
                        key={todo.id}
                        todo={todo}
                        openUpdate={() => openUpdate(todo)}
                        openDelete={() => openDelete(todo)}
                        openChangeStatus={() => openChangeStatus(todo)}
                        actions={messages.actions}
                        statusMessages={messages.status} />)}
            </TableBody>
        </Table>
    </TableContainer>
);

export default TodoContent;
