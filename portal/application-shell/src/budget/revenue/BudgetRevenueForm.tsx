import React from "react";

import { Box } from "@mui/material";
import { v1 as uuidv1 } from 'uuid';

import moment from "moment/moment";
import FormDatePicker, { FormDateFormatPattern } from "../../components/form/FormDatePicker";
import FormMoneyFormat from "../../components/form/FormMoneyFormat";
import FormSelect, { SelectOption } from "../../components/form/FormSelect";
import FormTextArea from "../../components/form/FormTextArea";
import { BudgetRevenueFormMessageBundle } from "../../messages/MessageBundles";

interface BudgetRevenueFormData {
    date: string;
    amount: string;
    note: string;
    searchTags: SelectOption[];
}

interface BudgetRevenueFormProps {
    budgetRevenueData: BudgetRevenueFormData;
    budgetRevenueHandlers: {
        date: (date: moment.Moment) => void;
        amount: (event: React.ChangeEvent<HTMLInputElement>) => void;
        note: (event: React.ChangeEvent<HTMLTextAreaElement>) => void;
        searchTag: (selectedOptions: SelectOption[]) => void;
    };
    searchTagRegistry: SelectOption[];
    messages: BudgetRevenueFormMessageBundle;
}
const BudgetRevenueForm: React.FC<BudgetRevenueFormProps> = ({ budgetRevenueData, budgetRevenueHandlers, searchTagRegistry, messages }) => {
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

        <FormSelect multi={true}
            id={uuidv1()}
            label={messages.searchTags}
            value={budgetRevenueData.searchTags}
            onChangeHandler={budgetRevenueHandlers.searchTag}
            options={searchTagRegistry} />

        <FormTextArea
            id={uuidv1()}
            value={budgetRevenueData.note}
            onChangeHandler={budgetRevenueHandlers.note}
            label={messages.note} />

    </Box>
}

export default BudgetRevenueForm;