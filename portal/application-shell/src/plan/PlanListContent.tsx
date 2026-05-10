import React from "react";
import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import Plan from "./domain/Plan";
import PlanListRow from "./PlanListRow";

interface PlanListContentProps {
    plans: Plan[];
    openDetail: (plan: Plan) => void;
    openDelete: (plan: Plan) => void;
}

const PlanListContent: React.FC<PlanListContentProps> = ({ plans, openDetail, openDelete }) => (
    <TableContainer component={Paper}>
        <Table>
            <TableHead>
                <TableRow>
                    <TableCell>Date</TableCell>
                    <TableCell>Title</TableCell>
                    <TableCell>Todos</TableCell>
                    <TableCell>Options</TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {plans.map(plan =>
                    <PlanListRow
                        key={plan.id}
                        plan={plan}
                        openDetail={() => openDetail(plan)}
                        openDelete={() => openDelete(plan)} />)}
            </TableBody>
        </Table>
    </TableContainer>
);

export default PlanListContent;
