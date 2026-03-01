import React from "react";
import BudgetRevenueRow from "./BudgetRevenueRow";
import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import BudgetRevenue from "./domain/BudgetRevenue";

interface BudgetRevenueTableProps {
    revenues: BudgetRevenue[];
    openDeletePopUp: (revenue: BudgetRevenue) => void;
    openUpdatePopUp: (revenue: BudgetRevenue) => void;
}

const BudgetRevenueTable: React.FC<BudgetRevenueTableProps> = ({ revenues, openDeletePopUp, openUpdatePopUp }) => {
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
            </TableBody>
        </Table>
    </TableContainer>

}

export default BudgetRevenueTable;