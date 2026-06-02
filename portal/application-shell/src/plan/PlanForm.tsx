import React from "react";
import { Box } from "@mui/material";
import moment from "moment";
import FormDatePicker, { FormDateFormatPattern } from "../components/form/FormDatePicker";
import FormInputTextField from "../components/form/FormInputTextField";

export type PlanFormData = {
    title: string;
    date: string;
};

interface PlanFormProps {
    data: PlanFormData;
    handlers: {
        title: (event: React.ChangeEvent<HTMLInputElement>) => void;
        date: (date: moment.Moment) => void;
    };
    messages: { title: string; date: string };
}

const PlanForm: React.FC<PlanFormProps> = ({ data, handlers, messages }) => (
    <Box>
        <FormInputTextField
            id="plan-title"
            label={messages.title}
            required
            autoFocus
            value={data.title}
            handler={handlers.title} />

        <FormDatePicker
            pattern={FormDateFormatPattern}
            label={messages.date}
            value={data.date}
            onClickHandler={handlers.date} />
    </Box>
);

export default PlanForm;
