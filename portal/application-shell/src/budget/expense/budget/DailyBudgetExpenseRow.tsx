import React from "react"
import { Button, ButtonGroup, TableCell, TableRow } from "@mui/material";
import { AttachFile, Delete, Edit } from "@mui/icons-material";
import { BudgetExpense, SavedBudgetExpense } from "../domain/BudgetExpense";

type DailyBudgetExpenseRowProps = {
    key: string,
    dailyBudgetExpense: BudgetExpense
    openUpdateBudgetExpensePopUp: (budgetExpense: SavedBudgetExpense) => void,
    openDeleteBudgetExpensePopUp: (budgetExpense: BudgetExpense) => void,
    openUploadAttachmentPopUp: (budgetExpense: BudgetExpense) => void
}

const DailyBudgetExpenseRow: React.FC<DailyBudgetExpenseRowProps> = ({ key, dailyBudgetExpense, openUpdateBudgetExpensePopUp, openDeleteBudgetExpensePopUp, openUploadAttachmentPopUp }: DailyBudgetExpenseRowProps) => {
    return <TableRow key={key} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
        <TableCell></TableCell>
        <TableCell>{dailyBudgetExpense.amount}</TableCell>
        <TableCell>{dailyBudgetExpense.note}</TableCell>
        <TableCell>{dailyBudgetExpense.tag.tagValue}</TableCell>
        <TableCell>
            <ButtonGroup variant="contained" aria-label="outlined primary button group">
                <Button onClick={() => openUpdateBudgetExpensePopUp({
                    id: dailyBudgetExpense.id,
                    amount: dailyBudgetExpense.amount,
                    date: dailyBudgetExpense.date,
                    note: dailyBudgetExpense.note,
                    searchTag: { value: dailyBudgetExpense.tag.tagKey, label: dailyBudgetExpense.tag.tagValue }
                })}><Edit /> Edit</Button>
                <Button onClick={() => openUploadAttachmentPopUp(dailyBudgetExpense)}><AttachFile /> Attach</Button>
                <Button onClick={() => openDeleteBudgetExpensePopUp(dailyBudgetExpense)}><Delete />Delete </Button>
            </ButtonGroup>
        </TableCell>
    </TableRow>

}


export default DailyBudgetExpenseRow
