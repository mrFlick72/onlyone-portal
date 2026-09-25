import React from "react"
import { Button, ButtonGroup, Chip, TableCell, TableRow } from "@mui/material";
import { Delete, Edit, Pause, PlayArrow } from "@mui/icons-material";
import ScheduledExpense from "./domain/ScheduledExpense";
import { colorFor } from "./domain/ScheduledExpenseStatusStyle";
import { ScheduledExpenseListContentMessageBundle } from "../../messages/MessageBundles";

type ScheduledExpenseListRowProps = {
    scheduledExpense: ScheduledExpense;
    openDetail: () => void;
    openDelete: () => void;
    toggleStatus: () => void;
    messages: ScheduledExpenseListContentMessageBundle;
}

const ScheduledExpenseListRow: React.FC<ScheduledExpenseListRowProps> = ({ scheduledExpense, openDetail, openDelete, toggleStatus, messages }) => (
    <TableRow key={scheduledExpense.id} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
        <TableCell>{scheduledExpense.description}</TableCell>
        <TableCell>{scheduledExpense.day}</TableCell>
        <TableCell>{scheduledExpense.month ?? messages.monthlyLabel}</TableCell>
        <TableCell>{scheduledExpense.amount}</TableCell>
        <TableCell>
            <Chip
                label={scheduledExpense.status ? messages.status[scheduledExpense.status] : ""}
                color={colorFor(scheduledExpense.status)}
                size="small" />
        </TableCell>
        <TableCell>
            <ButtonGroup variant="contained" aria-label="scheduled expense row actions">
                <Button onClick={openDetail}><Edit /> {messages.actions.open}</Button>
                {scheduledExpense.status === "PAUSED"
                    ? <Button onClick={toggleStatus}><PlayArrow /> {messages.actions.resume}</Button>
                    : <Button onClick={toggleStatus}><Pause /> {messages.actions.pause}</Button>}
                <Button onClick={openDelete}><Delete /> {messages.actions.delete}</Button>
            </ButtonGroup>
        </TableCell>
    </TableRow>
);

export default ScheduledExpenseListRow;
