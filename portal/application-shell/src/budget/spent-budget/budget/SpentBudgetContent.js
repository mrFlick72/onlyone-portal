import React from "react"
import DailyBudgetExpenseHeader from "./DailyBudgetExpenseHeader";
import DailyBudgetExpenseRow from "./DailyBudgetExpenseRow";
import {Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow} from "@mui/material";

export default ({spentBudget, openUpdateBudgetExpensePopUp, openDeleteBudgetExpensePopUp}) => {
    let tableContent = [];
    let dailyBudgetExpenseRepresentationList = spentBudget.dailyBudgetExpenseRepresentationList || []
    dailyBudgetExpenseRepresentationList.forEach((dailySpentBudget, dailyBudgetHeaderIndex) => {
        tableContent.push(<DailyBudgetExpenseHeader key={"H-" + dailyBudgetHeaderIndex}
                                                    date={dailySpentBudget.date}
                                                    total={dailySpentBudget.total}/>)

        dailySpentBudget.budgetExpenseRepresentationList.forEach((budgetExpenseRepresentation, dailyBudgetColumnIndex) => {
            tableContent.push(<DailyBudgetExpenseRow key={"C-" + dailyBudgetHeaderIndex + "-" + dailyBudgetColumnIndex}
                                                     dailyBudgetExpense={budgetExpenseRepresentation}
                                                     openUpdateBudgetExpensePopUp={openUpdateBudgetExpensePopUp.bind(budgetExpenseRepresentation)}
                                                     openDeleteBudgetExpensePopUp={openDeleteBudgetExpensePopUp.bind(budgetExpenseRepresentation)}/>)
        })
    });

    return <TableContainer component={Paper}>
        <Table>
            <TableHead>
                <TableRow>
                    <TableCell>Date</TableCell>
                    <TableCell>Amount</TableCell>
                    <TableCell>Note</TableCell>
                    <TableCell>Type</TableCell>
                    <TableCell>Details</TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {tableContent}

                <TableRow>
                    <TableCell>Total:</TableCell><TableCell/><TableCell/><TableCell/><TableCell align="right">{spentBudget.total}</TableCell>
                </TableRow>
            </TableBody>
        </Table>
    </TableContainer>
}