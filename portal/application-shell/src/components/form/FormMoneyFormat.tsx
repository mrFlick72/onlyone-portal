import React from "react";
import { Grid2, TextField } from "@mui/material";
import { NumericFormat } from "react-number-format";
import { commonStyle } from "../../theme/ThemeProvider";

type FormMoneyFormatProps = {
    id: string
    label: string
    required: boolean
    autoFocus?: boolean
    disabled?: boolean
    value: any
    handler: any
}


const FormMoneyFormat: React.FC<FormMoneyFormatProps> = ({ id, label, required, autoFocus, disabled, value, handler }) => {
    return <Grid2 container spacing={8} alignItems="flex-end" style={commonStyle.formRow}>
        <Grid2 size={{ xs: 12 }}>
            <NumericFormat
                size={"medium"}
                label={label}
                variant="outlined"
                required={required}
                autoFocus={autoFocus}
                disabled={disabled}
                value={value}
                onChange={handler}
                fullWidth={true}
                id={id}
                customInput={TextField}
                decimalScale={2}
            />
        </Grid2>
    </Grid2>
}

export default FormMoneyFormat