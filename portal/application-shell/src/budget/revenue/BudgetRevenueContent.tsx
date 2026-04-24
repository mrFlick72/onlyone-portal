import React from "react";
import BudgetRevenueRow from "./BudgetRevenueRow";
import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import BudgetRevenue from "./domain/BudgetRevenue";

interface BudgetRevenueTableProps {
    revenues: BudgetRevenue[];
    total: string;
    openDeletePopUp: (revenue: BudgetRevenue) => void;
    openUpdatePopUp: (revenue: BudgetRevenue) => void;
}

const BudgetRevenueTable: React.FC<BudgetRevenueTableProps> = ({ revenues, total, openDeletePopUp, openUpdatePopUp }) => {

    return <TableContainer component={Paper}>
        <Table>
            <TableHead>
                <TableRow>
                    <TableCell>Date</TableCell>
                    <TableCell>Amount</TableCell>
                    <TableCell>Note</TableCell>
                    <TableCell>Options</TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {revenues.map(revenue => <BudgetRevenueRow revenue={revenue}
                    openUpdatePopUp={openUpdatePopUp.bind(this, revenue)}
                    openDeletePopUp={openDeletePopUp.bind(this, revenue)} />)}
                <TableRow>
                    <TableCell sx={{ fontWeight: 'bold' }}>Total</TableCell>
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