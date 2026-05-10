import React from "react";
import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import { Todo } from "./domain/Plan";
import TodoRow from "./TodoRow";

interface TodoContentProps {
    todos: Todo[];
    openUpdate: (todo: Todo) => void;
    openDelete: (todo: Todo) => void;
}

const TodoContent: React.FC<TodoContentProps> = ({ todos, openUpdate, openDelete }) => (
    <TableContainer component={Paper}>
        <Table>
            <TableHead>
                <TableRow>
                    <TableCell>Date</TableCell>
                    <TableCell>Content</TableCell>
                    <TableCell>Options</TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {todos.map(todo =>
                    <TodoRow
                        key={todo.id}
                        todo={todo}
                        openUpdate={() => openUpdate(todo)}
                        openDelete={() => openDelete(todo)} />)}
            </TableBody>
        </Table>
    </TableContainer>
);

export default TodoContent;
