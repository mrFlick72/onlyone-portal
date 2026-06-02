import React from "react"
import { Button, ButtonGroup, TableCell, TableRow } from "@mui/material";
import { AttachFile, Delete, Edit } from "@mui/icons-material";
import { v1 as uuidv1 } from "uuid";
import BudgetRevenue from "./domain/BudgetRevenue";

interface BudgetRevenueRowProps {
    revenue: BudgetRevenue
    openDeletePopUp: () => void;
    openUpdatePopUp: () => void;
    openUploadAttachmentPopUp: () => void;
    actions: { edit: string; attach: string; delete: string };
}

const BudgetRevenueRow: React.FC<BudgetRevenueRowProps> = ({ revenue, openDeletePopUp, openUpdatePopUp, openUploadAttachmentPopUp, actions }) => {
    return <TableRow key={uuidv1()} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
        <TableCell>{revenue.date}</TableCell>
        <TableCell>{revenue.amount}</TableCell>
        <TableCell>{revenue.note}</TableCell>
        <TableCell>
            <ButtonGroup variant="contained" aria-label="outlined primary button group">
                <Button onClick={openUpdatePopUp}><Edit /> {actions.edit}</Button>
                <Button onClick={openUploadAttachmentPopUp}><AttachFile /> {actions.attach}</Button>
                <Button onClick={openDeletePopUp}><Delete /> {actions.delete}</Button>
            </ButtonGroup>
        </TableCell>
    </TableRow>

}


export default BudgetRevenueRow;
