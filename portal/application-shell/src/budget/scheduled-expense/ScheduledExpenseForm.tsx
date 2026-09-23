import React from "react";
import { Box, Checkbox, FormControlLabel } from "@mui/material";
import moment from "moment/moment";
import { v1 as uuidv1 } from 'uuid';
import FormDatePicker, { FormDateFormatPattern } from "../../components/form/FormDatePicker";
import FormInputTextField from "../../components/form/FormInputTextField";
import FormMoneyFormat from "../../components/form/FormMoneyFormat";
import FormSelect, { SelectOption } from "../../components/form/FormSelect";
import FormTextArea from "../../components/form/FormTextArea";
import { ScheduledExpenseFormMessageBundle } from "../../messages/MessageBundles";

// day/month are plain strings here (not numbers) so the form can represent
// "not yet typed" / "left empty" without a 0 that would look like a real
// value — day is required before submit, month stays optional (empty means
// monthly recurrence, see domain/ScheduledExpense.ts).
export type ScheduledExpenseFormData = {
    description: string;
    amount: string;
    note: string;
    searchTags: SelectOption[];
    day: string;
    month: string;
    hasEndDate: boolean;
    endDate: string;
}

type ScheduledExpenseFormProps = {
    data: ScheduledExpenseFormData;
    handlers: {
        description: (event: React.ChangeEvent<HTMLInputElement>) => void;
        amount: (event: React.ChangeEvent<HTMLInputElement>) => void;
        note: (event: React.ChangeEvent<HTMLTextAreaElement>) => void;
        searchTag: (selectedOptions: SelectOption[]) => void;
        day: (event: React.ChangeEvent<HTMLInputElement>) => void;
        month: (event: React.ChangeEvent<HTMLInputElement>) => void;
        toggleEndDate: (event: React.ChangeEvent<HTMLInputElement>) => void;
        endDate: (date: moment.Moment) => void;
    };
    searchTagRegistry: SelectOption[];
    messages: ScheduledExpenseFormMessageBundle;
}

const ScheduledExpenseForm: React.FC<ScheduledExpenseFormProps> = ({ data, handlers, searchTagRegistry, messages }) => (
    <Box>
        <FormInputTextField
            id="scheduled-expense-description"
            label={messages.description}
            required
            autoFocus
            value={data.description}
            handler={handlers.description} />

        <FormMoneyFormat
            id="scheduled-expense-amount"
            label={messages.amount}
            required={true}
            handler={handlers.amount}
            value={data.amount} />

        <FormSelect multi={true}
            id={uuidv1()}
            label={messages.searchTags}
            value={data.searchTags}
            onChangeHandler={handlers.searchTag}
            options={searchTagRegistry} />

        <FormInputTextField
            id="scheduled-expense-day"
            label={messages.day}
            type="number"
            required
            value={data.day}
            handler={handlers.day} />

        <FormInputTextField
            id="scheduled-expense-month"
            label={messages.month}
            type="number"
            value={data.month}
            handler={handlers.month} />

        <FormControlLabel
            control={<Checkbox checked={data.hasEndDate} onChange={handlers.toggleEndDate} />}
            label={messages.setEndDate} />

        {data.hasEndDate && <FormDatePicker
            pattern={FormDateFormatPattern}
            label={messages.endDate}
            value={data.endDate}
            onClickHandler={handlers.endDate} />}

        <FormTextArea
            id={uuidv1()}
            value={data.note}
            onChangeHandler={handlers.note}
            label={messages.note} />
    </Box>
);

export default ScheduledExpenseForm;
