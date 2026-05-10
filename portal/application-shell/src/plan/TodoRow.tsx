import React from "react"
import { Button, ButtonGroup, TableCell, TableRow } from "@mui/material";
import { Delete, Edit } from "@mui/icons-material";
import moment from "moment";
import { Todo } from "./domain/Plan";
import { ApiDateFormatPattern, FormDateFormatPattern } from "../components/form/FormDatePicker";

interface TodoRowProps {
    todo: Todo;
    openUpdate: () => void;
    openDelete: () => void;
}

const TodoRow: React.FC<TodoRowProps> = ({ todo, openUpdate, openDelete }) => (
    <TableRow key={todo.id} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
        <TableCell>{moment(todo.date, ApiDateFormatPattern).format(FormDateFormatPattern)}</TableCell>
        <TableCell>{todo.content}</TableCell>
        <TableCell>
            <ButtonGroup variant="contained" aria-label="todo row actions">
                <Button onClick={openUpdate}><Edit /> Edit</Button>
                <Button onClick={openDelete}><Delete /> Delete</Button>
            </ButtonGroup>
        </TableCell>
    </TableRow>
);

export default TodoRow;
