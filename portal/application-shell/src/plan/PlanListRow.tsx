import React from "react"
import { Button, ButtonGroup, TableCell, TableRow } from "@mui/material";
import { Delete, FormatListBulleted } from "@mui/icons-material";
import moment from "moment";
import Plan from "./domain/Plan";
import { ApiDateFormatPattern, FormDateFormatPattern } from "../components/form/FormDatePicker";
import { PlanListRowActionsMessageBundle } from "../messages/MessageBundles";

type PlanListRowProps = {
    plan: Plan;
    openDetail: () => void;
    openDelete: () => void;
    actions: PlanListRowActionsMessageBundle;
}

const PlanListRow: React.FC<PlanListRowProps> = ({ plan, openDetail, openDelete, actions }) => (
    <TableRow key={plan.id} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
        <TableCell>{moment(plan.date, ApiDateFormatPattern).format(FormDateFormatPattern)}</TableCell>
        <TableCell>{plan.title}</TableCell>
        <TableCell>{plan.todo_count}</TableCell>
        <TableCell>
            <ButtonGroup variant="contained" aria-label="plan row actions">
                <Button onClick={openDetail}><FormatListBulleted /> {actions.open}</Button>
                <Button onClick={openDelete}><Delete /> {actions.delete}</Button>
            </ButtonGroup>
        </TableCell>
    </TableRow>
);

export default PlanListRow;
