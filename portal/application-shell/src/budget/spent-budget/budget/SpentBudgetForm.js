import React from "react"

import {v1 as uuidv1} from 'uuid';
import {Box} from "@mui/material";
import moment from "moment";
import FormDatePicker, {FormDateFormatPattern} from "../../../components/form/FormDatePicker";
import FormMoneyFormat from "../../../components/form/FormMoneyFormat";
import FormSelect from "../../../components/form/FormSelect";
import FormTextArea from "../../../components/form/FormTextArea";

export default ({spentBudgetComponentId, spentBudgetData, spentBudgetHandlers, searchTagRegistry}) => {
    let spentBudgetComponentIdAux = spentBudgetComponentId || {}
    return <Box>
        <FormDatePicker
            label={"Date:"}
            value={moment(spentBudgetData.date, FormDateFormatPattern, true)}
            onClickHandler={spentBudgetHandlers.date}/>

        <FormMoneyFormat
            id={spentBudgetComponentIdAux.amount || uuidv1()}
            label="Amount:"
            required={true}
            handler={spentBudgetHandlers.amount}
            value={spentBudgetData.amount}/>

        <FormSelect multi={false}
                    componentId={spentBudgetComponentIdAux.searchTag || uuidv1()}
                    componentLabel="Search Tag:"
                    value={spentBudgetData.searchTag}
                    onChangeHandler={spentBudgetHandlers.searchTag}
                    options={searchTagRegistry}/>

        <FormTextArea value={spentBudgetData.note}
                      onChangeHandler={spentBudgetHandlers.note}
                      id={spentBudgetComponentIdAux.note || uuidv1()}
                      label="Note:"/>
    </Box>
}