import React from "react"

import { v1 as uuidv1 } from 'uuid';
import { Box } from "@mui/material";
import FormDatePicker, { FormDateFormatPattern } from "../../../components/form/FormDatePicker";
import FormMoneyFormat from "../../../components/form/FormMoneyFormat";
import FormSelect from "../../../components/form/FormSelect";
import FormTextArea from "../../../components/form/FormTextArea";
import { SpentBudgetFormMessageBundle } from "../../../messages/MessageBundles";

type SpentBudgetFormProps = {
    spentBudgetData: {
        date: string,
        amount: number,
        note: string,
        searchTags: { value: string, label: string }[]
    },
    spentBudgetHandlers: {
        date: (args: any) => void,
        amount: (...args: any) => void,
        note: (...args: any) => void,
        searchTag: (...args: any) => void
    },
    searchTagRegistry: { value: string, label: string }[],
    messages: SpentBudgetFormMessageBundle
}

const SpentBudgetForm: React.FC<SpentBudgetFormProps> = ({ spentBudgetData, spentBudgetHandlers, searchTagRegistry, messages }) => {
    console.log("SpentBudgetForm render")
    console.log("SpentBudgetForm render with data: ", spentBudgetData)
    return <Box>
        <FormDatePicker
            pattern={FormDateFormatPattern}
            label={messages.date}
            value={spentBudgetData.date}
            onClickHandler={spentBudgetHandlers.date} />

        <FormMoneyFormat
            id={uuidv1()}
            label={messages.amount}
            required={true}
            handler={spentBudgetHandlers.amount}
            value={spentBudgetData.amount} />

        <FormSelect multi={true}
            id={uuidv1()}
            label={messages.searchTags}
            value={spentBudgetData.searchTags}
            onChangeHandler={spentBudgetHandlers.searchTag}
            options={searchTagRegistry} />

        <FormTextArea value={spentBudgetData.note}
            onChangeHandler={spentBudgetHandlers.note}
            id={uuidv1()}
            label={messages.note} />
    </Box>
}

export default SpentBudgetForm