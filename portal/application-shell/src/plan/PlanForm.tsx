import React from "react";
import { Box } from "@mui/material";
import moment from "moment";
import FormDatePicker, { FormDateFormatPattern } from "../components/form/FormDatePicker";
import FormInputTextField from "../components/form/FormInputTextField";
import { PlanFormMessageBundle } from "../messages/MessageBundles";

export type PlanFormData = {
    title: string;
    date: string;
};

type PlanFormProps = {
    data: PlanFormData;
    handlers: {
        title: (event: React.ChangeEvent<HTMLInputElement>) => void;
        date: (date: moment.Moment) => void;
    };
    messages: PlanFormMessageBundle;
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
