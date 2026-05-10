import React from "react"
import { Button, ButtonGroup, TableCell, TableRow } from "@mui/material";
import { Delete, FormatListBulleted } from "@mui/icons-material";
import moment from "moment";
import Plan from "./domain/Plan";
import { ApiDateFormatPattern, FormDateFormatPattern } from "../components/form/FormDatePicker";

interface PlanListRowProps {
    plan: Plan;
    openDetail: () => void;
    openDelete: () => void;
}

const PlanListRow: React.FC<PlanListRowProps> = ({ plan, openDetail, openDelete }) => (
    <TableRow key={plan.id} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
        <TableCell>{moment(plan.date, ApiDateFormatPattern).format(FormDateFormatPattern)}</TableCell>
        <TableCell>{plan.title}</TableCell>
        <TableCell>
            <ButtonGroup variant="contained" aria-label="plan row actions">
                <Button onClick={openDetail}><FormatListBulleted /> Open</Button>
                <Button onClick={openDelete}><Delete /> Delete</Button>
            </ButtonGroup>
        </TableCell>
    </TableRow>
);

export default PlanListRow;
