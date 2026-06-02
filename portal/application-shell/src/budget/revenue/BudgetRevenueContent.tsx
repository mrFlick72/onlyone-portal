import React from "react";
import BudgetRevenueRow from "./BudgetRevenueRow";
import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import BudgetRevenue from "./domain/BudgetRevenue";

interface BudgetRevenueTableMessages {
    headers: { date: string; amount: string; note: string; options: string };
    totalLabel: string;
    actions: { edit: string; attach: string; delete: string };
}

interface BudgetRevenueTableProps {
    revenues: BudgetRevenue[];
    total: string;
    openDeletePopUp: (revenue: BudgetRevenue) => void;
    openUpdatePopUp: (revenue: BudgetRevenue) => void;
    openUploadAttachmentPopUp: (revenue: BudgetRevenue) => void;
    messages: BudgetRevenueTableMessages;
}

const BudgetRevenueTable: React.FC<BudgetRevenueTableProps> = ({ revenues, total, openDeletePopUp, openUpdatePopUp, openUploadAttachmentPopUp, messages }) => {

    return <TableContainer component={Paper}>
        <Table>
            <TableHead>
                <TableRow>
                    <TableCell>{messages.headers.date}</TableCell>
                    <TableCell>{messages.headers.amount}</TableCell>
                    <TableCell>{messages.headers.note}</TableCell>
                    <TableCell>{messages.headers.options}</TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {revenues.map(revenue => <BudgetRevenueRow revenue={revenue}
                    openUpdatePopUp={openUpdatePopUp.bind(this, revenue)}
                    openDeletePopUp={openDeletePopUp.bind(this, revenue)}
                    openUploadAttachmentPopUp={openUploadAttachmentPopUp.bind(this, revenue)}
                    actions={messages.actions} />)}
                <TableRow>
                    <TableCell sx={{ fontWeight: 'bold' }}>{messages.totalLabel}</TableCell>
                    <TableCell></TableCell>
                    <TableCell></TableCell>
                    <TableCell sx={{ fontWeight: 'bold' }}>{total}</TableCell>
                    <TableCell colSpan={2} />
                </TableRow>
            </TableBody>
        </Table>
    </TableContainer>

}

export default BudgetRevenueTable;