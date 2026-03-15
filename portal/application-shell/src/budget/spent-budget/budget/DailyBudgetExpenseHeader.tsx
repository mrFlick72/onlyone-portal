import React from "react"
import { TableCell, TableRow } from "@mui/material";

type DailyBudgetExpenseHeaderProps = {
    key: string,
    date: string,
    total: number
}

const DailyBudgetExpenseHeader: React.FC<DailyBudgetExpenseHeaderProps> = ({ key, date, total }) => {
    return (
        <TableRow key={key} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
            <TableCell>{date}</TableCell>
            <TableCell></TableCell>
            <TableCell></TableCell>
            <TableCell></TableCell>
            <TableCell>Total: {total}</TableCell>
        </TableRow>)
}

export default DailyBudgetExpenseHeader 