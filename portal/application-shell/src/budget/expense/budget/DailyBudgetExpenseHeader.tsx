import React from "react"
import { TableCell, TableRow } from "@mui/material";

type DailyBudgetExpenseHeaderProps = {
    key: string,
    date: string,
    total: string,
    totalLabel: string
}

const DailyBudgetExpenseHeader: React.FC<DailyBudgetExpenseHeaderProps> = ({ key, date, total, totalLabel }) => {
    return (
        <TableRow key={key} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
            <TableCell>{date}</TableCell>
            <TableCell></TableCell>
            <TableCell></TableCell>
            <TableCell></TableCell>
            <TableCell>{totalLabel} {total}</TableCell>
        </TableRow>)
}

export default DailyBudgetExpenseHeader 