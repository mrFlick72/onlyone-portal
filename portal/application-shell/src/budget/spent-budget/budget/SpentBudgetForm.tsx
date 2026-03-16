import React from "react"

import { v1 as uuidv1 } from 'uuid';
import { Box } from "@mui/material";
import FormDatePicker, { FormDateFormatPattern } from "../../../components/form/FormDatePicker";
import FormMoneyFormat from "../../../components/form/FormMoneyFormat";
import FormSelect from "../../../components/form/FormSelect";
import FormTextArea from "../../../components/form/FormTextArea";

type SpentBudgetFormProps = {
    spentBudgetData: {
        date: string,
        amount: number,
        note: string,
        searchTag: { value: string, label: string }
    },
    spentBudgetHandlers: {
        date: (args: any) => void,
        amount: (...args: any) => void,
        note: (...args: any) => void,
        searchTag: (...args: any) => void
    },
    searchTagRegistry: { value: string, label: string }[]
}

const SpentBudgetForm: React.FC<SpentBudgetFormProps> = ({ spentBudgetData, spentBudgetHandlers, searchTagRegistry }) => {
    console.log("SpentBudgetForm render")
    console.log("SpentBudgetForm render with data: ", spentBudgetData)
    return <Box>
        <FormDatePicker
            pattern={FormDateFormatPattern}
            label={"Date:"}
            value={spentBudgetData.date}
            onClickHandler={spentBudgetHandlers.date} />

        <FormMoneyFormat
            id={uuidv1()}
            label="Amount:"
            required={true}
            handler={spentBudgetHandlers.amount}
            value={spentBudgetData.amount} />

        <FormSelect multi={false}
            id={uuidv1()}
            label="Search Tag:"
            value={spentBudgetData.searchTag}
            onChangeHandler={spentBudgetHandlers.searchTag}
            options={searchTagRegistry} />

        <FormTextArea value={spentBudgetData.note}
            onChangeHandler={spentBudgetHandlers.note}
            id={uuidv1()}
            label="Note:" />
    </Box>
}

export default SpentBudgetForm