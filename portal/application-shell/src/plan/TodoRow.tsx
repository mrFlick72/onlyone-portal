import React from "react"
import { Button, ButtonGroup, Chip, TableCell, TableRow } from "@mui/material";
import { Delete, Edit, SwapHoriz } from "@mui/icons-material";
import moment from "moment";
import { Todo, TodoStatus } from "./domain/Plan";
import { ApiDateFormatPattern, FormDateFormatPattern } from "../components/form/FormDatePicker";
import { colorFor, isTerminal } from "./domain/TodoStatus";

interface TodoRowProps {
    todo: Todo;
    openUpdate: () => void;
    openDelete: () => void;
    openChangeStatus: () => void;
    actions: { changeStatus: string; edit: string; delete: string };
    statusMessages: Record<TodoStatus, string>;
}

const TodoRow: React.FC<TodoRowProps> = ({ todo, openUpdate, openDelete, openChangeStatus, actions, statusMessages }) => (
    <TableRow key={todo.id} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
        <TableCell>{moment(todo.date, ApiDateFormatPattern).format(FormDateFormatPattern)}</TableCell>
        <TableCell>{todo.content}</TableCell>
        <TableCell>
            <Chip label={statusMessages[todo.status]} color={colorFor(todo.status)} size="small" />
        </TableCell>
        <TableCell>
            <ButtonGroup variant="contained" aria-label="todo row actions">
                <Button
                    onClick={openChangeStatus}
                    disabled={isTerminal(todo.status)}
                    aria-label="change status">
                    <SwapHoriz /> {actions.changeStatus}
                </Button>
                <Button onClick={openUpdate}><Edit /> {actions.edit}</Button>
                <Button onClick={openDelete}><Delete /> {actions.delete}</Button>
            </ButtonGroup>
        </TableCell>
    </TableRow>
);

export default TodoRow;
