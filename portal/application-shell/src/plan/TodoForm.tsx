import React from "react";
import { Box } from "@mui/material";
import { v1 as uuidv1 } from "uuid";
import moment from "moment";
import FormDatePicker, { FormDateFormatPattern } from "../components/form/FormDatePicker";
import FormTextArea from "../components/form/FormTextArea";

export type TodoFormData = {
    content: string;
    date: string;
};

interface TodoFormProps {
    data: TodoFormData;
    handlers: {
        content: (event: React.ChangeEvent<HTMLTextAreaElement>) => void;
        date: (date: moment.Moment) => void;
    };
}

const TodoForm: React.FC<TodoFormProps> = ({ data, handlers }) => (
    <Box>
        <FormDatePicker
            pattern={FormDateFormatPattern}
            label="Date:"
            value={data.date}
            onClickHandler={handlers.date} />

        <FormTextArea
            id={uuidv1()}
            label="Content:"
            value={data.content}
            onChangeHandler={handlers.content} />
    </Box>
);

export default TodoForm;
