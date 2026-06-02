import React from "react";
import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import Plan from "./domain/Plan";
import PlanListRow from "./PlanListRow";

interface PlanListContentMessages {
    headers: { date: string; title: string; todos: string; options: string };
    actions: { open: string; delete: string };
}

interface PlanListContentProps {
    plans: Plan[];
    openDetail: (plan: Plan) => void;
    openDelete: (plan: Plan) => void;
    messages: PlanListContentMessages;
}

const PlanListContent: React.FC<PlanListContentProps> = ({ plans, openDetail, openDelete, messages }) => (
    <TableContainer component={Paper}>
        <Table>
            <TableHead>
                <TableRow>
                    <TableCell>{messages.headers.date}</TableCell>
                    <TableCell>{messages.headers.title}</TableCell>
                    <TableCell>{messages.headers.todos}</TableCell>
                    <TableCell>{messages.headers.options}</TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {plans.map(plan =>
                    <PlanListRow
                        key={plan.id}
                        plan={plan}
                        openDetail={() => openDetail(plan)}
                        openDelete={() => openDelete(plan)}
                        actions={messages.actions} />)}
            </TableBody>
        </Table>
    </TableContainer>
);

export default PlanListContent;
