import React from "react";

import { Box } from "@mui/material";
import { v1 as uuidv1 } from 'uuid';

import moment from "moment/moment";
import FormDatePicker, { FormDateFormatPattern } from "../../components/form/FormDatePicker";
import FormMoneyFormat from "../../components/form/FormMoneyFormat";
import FormTextArea from "../../components/form/FormTextArea";
import BudgetRevenue from "./domain/BudgetRevenue";

interface BudgetRevenueFormProps {
    budgetRevenueData: BudgetRevenue;
    budgetRevenueHandlers: {
        date: (date: moment.Moment) => void;
        amount: (event: React.ChangeEvent<HTMLInputElement>) => void;
        note: (event: React.ChangeEvent<HTMLTextAreaElement>) => void;
    };
    messages: { date: string; amount: string; note: string };
}
const BudgetRevenueForm: React.FC<BudgetRevenueFormProps> = ({ budgetRevenueData, budgetRevenueHandlers, messages }) => {
    return <Box>
        <FormDatePicker
            pattern={FormDateFormatPattern}
            label={messages.date}
            value={budgetRevenueData.date}
            onClickHandler={budgetRevenueHandlers.date} />

        <FormMoneyFormat
            id="amount"
            label={messages.amount}
            required={true}
            handler={budgetRevenueHandlers.amount}
            value={budgetRevenueData.amount} />


        <FormTextArea
            id={uuidv1()}
            value={budgetRevenueData.note}
            onChangeHandler={budgetRevenueHandlers.note}
            label={messages.note} />

    </Box>
}

export default BudgetRevenueForm;