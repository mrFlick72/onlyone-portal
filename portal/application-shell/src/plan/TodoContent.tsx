import React from "react";
import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import { Todo } from "./domain/Plan";
import TodoRow from "./TodoRow";
import { TodoContentMessageBundle } from "../messages/MessageBundles";

interface TodoContentProps {
    todos: Todo[];
    openUpdate: (todo: Todo) => void;
    openDelete: (todo: Todo) => void;
    openChangeStatus: (todo: Todo) => void;
    messages: TodoContentMessageBundle;
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
