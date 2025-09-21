import React from "react"
import {Button, ButtonGroup, TableCell, TableRow} from "@mui/material";
import {Delete, Edit} from "@mui/icons-material";

export default ({key, dailyBudgetExpense, openUpdateBudgetExpensePopUp, openDeleteBudgetExpensePopUp}) => {
    let budgetExpense = {
        "id": dailyBudgetExpense.id,
        "date": dailyBudgetExpense.date,
        "amount": dailyBudgetExpense.amount,
        "note": dailyBudgetExpense.note,
        "searchTag": {"value": dailyBudgetExpense.tagKey, "label": dailyBudgetExpense.tagValue}
    };

    return <TableRow key={key} sx={{'&:last-child td, &:last-child th': {border: 0}}}>
        <TableCell></TableCell>
        <TableCell>{dailyBudgetExpense.amount}</TableCell>
        <TableCell>{dailyBudgetExpense.note}</TableCell>
        <TableCell>{dailyBudgetExpense.tagValue}</TableCell>
        <TableCell>
            <ButtonGroup variant="contained" aria-label="outlined primary button group">
                <Button onClick={openUpdateBudgetExpensePopUp.bind(this, budgetExpense)}><Edit/> Edit</Button>
                <Button onClick={openDeleteBudgetExpensePopUp.bind(this, dailyBudgetExpense)}><Delete/>Delete </Button>
            </ButtonGroup>
        </TableCell>
    </TableRow>

}