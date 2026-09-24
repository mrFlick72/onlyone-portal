import React from "react";
import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import ScheduledExpense from "./domain/ScheduledExpense";
import ScheduledExpenseListRow from "./ScheduledExpenseListRow";
import { ScheduledExpenseListContentMessageBundle } from "../../messages/MessageBundles";

type ScheduledExpenseListContentProps = {
    scheduledExpenses: ScheduledExpense[];
    openDetail: (scheduledExpense: ScheduledExpense) => void;
    openDelete: (scheduledExpense: ScheduledExpense) => void;
    messages: ScheduledExpenseListContentMessageBundle;
}

const ScheduledExpenseListContent: React.FC<ScheduledExpenseListContentProps> = ({ scheduledExpenses, openDetail, openDelete, messages }) => (
    <TableContainer component={Paper}>
        <Table>
            <TableHead>
                <TableRow>
                    <TableCell>{messages.headers.description}</TableCell>
                    <TableCell>{messages.headers.day}</TableCell>
                    <TableCell>{messages.headers.month}</TableCell>
                    <TableCell>{messages.headers.amount}</TableCell>
                    <TableCell>{messages.headers.status}</TableCell>
                    <TableCell>{messages.headers.options}</TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {scheduledExpenses.map(scheduledExpense =>
                    <ScheduledExpenseListRow
                        key={scheduledExpense.id}
                        scheduledExpense={scheduledExpense}
                        openDetail={() => openDetail(scheduledExpense)}
                        openDelete={() => openDelete(scheduledExpense)}
                        messages={messages} />)}
            </TableBody>
        </Table>
    </TableContainer>
);

export default ScheduledExpenseListContent;
