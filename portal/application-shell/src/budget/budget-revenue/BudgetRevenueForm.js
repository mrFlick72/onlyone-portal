import {v1 as uuidv1} from 'uuid';
import {Box} from "@mui/material";
import moment from "moment/moment";
import FormDatePicker, {FormDateFormatPattern} from "../../components/form/FormDatePicker";
import FormMoneyFormat from "../../components/form/FormMoneyFormat";
import FormTextArea from "../../components/form/FormTextArea";

export default ({budgetRevenueData, budgetRevenueHandlers}) => {
    return <Box>
        <FormDatePicker
            id={uuidv1()}
            label={"Date:"}
            value={moment(budgetRevenueData.date, FormDateFormatPattern)}
            onClickHandler={budgetRevenueHandlers.date}/>

        <FormMoneyFormat
            id={uuidv1()}
            label="Amount:"
            required={true}
            handler={budgetRevenueHandlers.amount}
            value={budgetRevenueData.amount}/>


        <FormTextArea id={uuidv1()}
                      value={budgetRevenueData.note}
                      onChangeHandler={budgetRevenueHandlers.note}
                      label="Note:"/>

    </Box>
}