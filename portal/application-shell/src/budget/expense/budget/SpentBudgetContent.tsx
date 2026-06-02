import React from "react"
import DailyBudgetExpenseHeader from "./DailyBudgetExpenseHeader";
import DailyBudgetExpenseRow from "./DailyBudgetExpenseRow";
import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import { BudgetExpense, SavedBudgetExpense, SpentBudget } from "../domain/BudgetExpense";
import { SpentBudgetContentMessageBundle } from "../../../messages/MessageBundles";

type SpentBudgetContentProps = {
    spentBudget: SpentBudget
    openUpdateBudgetExpensePopUp: (budgetExpense: SavedBudgetExpense) => void,
    openDeleteBudgetExpensePopUp: (budgetExpense: BudgetExpense) => void,
    openUploadAttachmentPopUp: (budgetExpense: BudgetExpense) => void,
    messages: SpentBudgetContentMessageBundle
}

const SpentBudgetContent: React.FC<SpentBudgetContentProps> = ({ spentBudget, openUpdateBudgetExpensePopUp, openDeleteBudgetExpensePopUp, openUploadAttachmentPopUp, messages }) => {
    const tableContent: React.ReactNode[] = [];
    const dailyBudgetExpenseRepresentationList = spentBudget.dailyBudgetExpenseRepresentationList || []

    dailyBudgetExpenseRepresentationList.forEach((dailySpentBudget, dailyBudgetHeaderIndex) => {

        tableContent.push(<DailyBudgetExpenseHeader key={"H-" + dailyBudgetHeaderIndex}
            date={dailySpentBudget.date}
            totalLabel={messages.totalLabel}
            total={dailySpentBudget.total} />)

        dailySpentBudget.budgetExpenseRepresentationList.forEach((budgetExpenseRepresentation, dailyBudgetColumnIndex) => {
            tableContent.push(<DailyBudgetExpenseRow key={"C-" + dailyBudgetHeaderIndex + "-" + dailyBudgetColumnIndex}
                dailyBudgetExpense={budgetExpenseRepresentation}
                openUpdateBudgetExpensePopUp={openUpdateBudgetExpensePopUp.bind({
                    id: budgetExpenseRepresentation.id,
                    date: budgetExpenseRepresentation.date,
                    amount: budgetExpenseRepresentation.amount,
                    note: budgetExpenseRepresentation.note,
                    searchTags: budgetExpenseRepresentation.tags.map((tag) => ({ value: tag.tagKey, label: tag.tagValue }))

                })}
                openDeleteBudgetExpensePopUp={openDeleteBudgetExpensePopUp.bind(budgetExpenseRepresentation)}
                openUploadAttachmentPopUp={openUploadAttachmentPopUp.bind(budgetExpenseRepresentation)}
                actions={messages.actions} />)
        })
    });

    return <TableContainer component={Paper}>
        <Table>
            <TableHead>
                <TableRow>
                    <TableCell>{messages.headers.date}</TableCell>
                    <TableCell>{messages.headers.amount}</TableCell>
                    <TableCell>{messages.headers.note}</TableCell>
                    <TableCell>{messages.headers.type}</TableCell>
                    <TableCell>{messages.headers.details}</TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {tableContent}

                <TableRow>
                    <TableCell>{messages.totalLabel}</TableCell><TableCell /><TableCell /><TableCell /><TableCell align="right">{spentBudget.total}</TableCell>
                </TableRow>
            </TableBody>
        </Table>
    </TableContainer>
}

export default SpentBudgetContent