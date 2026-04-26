import React from "react";
import {Grid} from "@mui/material";
import {DesktopDatePicker} from "@mui/x-date-pickers"
import {AdapterMoment} from '@mui/x-date-pickers/AdapterMoment';
import {LocalizationProvider} from '@mui/x-date-pickers/LocalizationProvider';
import {commonStyle} from "../../theme/ThemeProvider";
import moment from "moment";

interface FormDatePickerProps {
    label: string
    value: string
    onClickHandler: (value: any) => void
    pattern: string

}


const FormDatePicker: React.FC<FormDatePickerProps> = ({
                                                           label, value, onClickHandler, pattern
                                                       }) => {
    let val = value && moment(value, pattern)

    return <Grid container sx={{ alignItems: "flex-end", ...commonStyle.formRow }}>
        <Grid size={{xs: 12}}>
            <LocalizationProvider dateAdapter={AdapterMoment}>
                <DesktopDatePicker
                    slotProps={{textField: {fullWidth: true}}}
                    label={label}
                    format={FormDateFormatPattern}
                    onChange={onClickHandler || {}}
                    value={val || moment()}
                />
            </LocalizationProvider>
        </Grid>
    </Grid>

}

export const ApiDateFormatPattern: string = 'YYYY-MM-DD'
export const FormDateFormatPattern = 'DD/MM/YYYY'
export default FormDatePicker

