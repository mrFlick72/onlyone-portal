import React from "react"
import BudgetRevenueRow from "./BudgetRevenueRow";
import {Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow} from "@mui/material";

export default ({revenues, openDeletePopUp, openUpdatePopUp}) => {
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
                                                           openDeletePopUp={openDeletePopUp.bind(this, revenue)}/>)}
            </TableBody>
        </Table>
    </TableContainer>

}