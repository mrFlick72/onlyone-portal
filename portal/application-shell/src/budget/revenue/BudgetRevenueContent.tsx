import React from "react";
import BudgetRevenueRow from "./BudgetRevenueRow";
import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import BudgetRevenue from "./domain/BudgetRevenue";
import { BudgetRevenueTableMessageBundle } from "../../messages/MessageBundles";

type BudgetRevenueTableProps = {
    revenues: BudgetRevenue[];
    total: string;
    openDeletePopUp: (revenue: BudgetRevenue) => void;
    openUpdatePopUp: (revenue: BudgetRevenue) => void;
    openUploadAttachmentPopUp: (revenue: BudgetRevenue) => void;
    messages: BudgetRevenueTableMessageBundle;
}

const BudgetRevenueTable: React.FC<BudgetRevenueTableProps> = ({ revenues, total, openDeletePopUp, openUpdatePopUp, openUploadAttachmentPopUp, messages }) => {

    return <TableContainer component={Paper}>
        <Table>
            <TableHead>
                <TableRow>
                    <TableCell>{messages.headers.date}</TableCell>
                    <TableCell>{messages.headers.amount}</TableCell>
                    <TableCell>{messages.headers.note}</TableCell>
                    <TableCell>{messages.headers.tags}</TableCell>
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
                    <TableCell sx={{ fontWeight: 'bold' }}>{total}</TableCell>
                    <TableCell></TableCell>
                    <TableCell></TableCell>
                    <TableCell></TableCell>
                </TableRow>
            </TableBody>
        </Table>
    </TableContainer>

}

export default BudgetRevenueTable;