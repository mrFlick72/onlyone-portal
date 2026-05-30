import React from "react"
import DailyBudgetExpenseHeader from "./DailyBudgetExpenseHeader";
import DailyBudgetExpenseRow from "./DailyBudgetExpenseRow";
import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import { BudgetExpense, SavedBudgetExpense, SpentBudget } from "../domain/BudgetExpense";

type SpentBudgetContentProps = {
    spentBudget: SpentBudget
    openUpdateBudgetExpensePopUp: (budgetExpense: SavedBudgetExpense) => void,
    openDeleteBudgetExpensePopUp: (budgetExpense: BudgetExpense) => void,
    openUploadAttachmentPopUp: (budgetExpense: BudgetExpense) => void
}

const SpentBudgetContent: React.FC<SpentBudgetContentProps> = ({ spentBudget, openUpdateBudgetExpensePopUp, openDeleteBudgetExpensePopUp, openUploadAttachmentPopUp }) => {
    const tableContent: React.ReactNode[] = [];
    const dailyBudgetExpenseRepresentationList = spentBudget.dailyBudgetExpenseRepresentationList || []

    dailyBudgetExpenseRepresentationList.forEach((dailySpentBudget, dailyBudgetHeaderIndex) => {

        tableContent.push(<DailyBudgetExpenseHeader key={"H-" + dailyBudgetHeaderIndex}
            date={dailySpentBudget.date}
            total={dailySpentBudget.total} />)

        dailySpentBudget.budgetExpenseRepresentationList.forEach((budgetExpenseRepresentation, dailyBudgetColumnIndex) => {
            tableContent.push(<DailyBudgetExpenseRow key={"C-" + dailyBudgetHeaderIndex + "-" + dailyBudgetColumnIndex}
                dailyBudgetExpense={budgetExpenseRepresentation}
                openUpdateBudgetExpensePopUp={openUpdateBudgetExpensePopUp.bind({
                    id: budgetExpenseRepresentation.id,
                    date: budgetExpenseRepresentation.date,
                    amount: budgetExpenseRepresentation.amount,
                    note: budgetExpenseRepresentation.note,
                    searchTag: { value: budgetExpenseRepresentation.tag?.tagKey ?? "", label: budgetExpenseRepresentation.tag?.tagValue ?? "" }

                })}
                openDeleteBudgetExpensePopUp={openDeleteBudgetExpensePopUp.bind(budgetExpenseRepresentation)}
                openUploadAttachmentPopUp={openUploadAttachmentPopUp.bind(budgetExpenseRepresentation)} />)
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
                    <TableCell>Total:</TableCell><TableCell /><TableCell /><TableCell /><TableCell align="right">{spentBudget.total}</TableCell>
                </TableRow>
            </TableBody>
        </Table>
    </TableContainer>
}

export default SpentBudgetContent