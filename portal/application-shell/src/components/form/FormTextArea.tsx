import React from "react";

import { Grid2, TextField } from "@mui/material";
import { commonStyle } from "../../theme/ThemeProvider";

type FormTextAreaProps = {
    id: string,
    label: string,
    type: string,
    required: boolean,
    autoFocus: boolean,
    disabled: boolean,
    value: string,
    onChangeHandler : React.ChangeEventHandler,
    row: number
}

const FormTextArea: React.FC<FormTextAreaProps> = ({ id, label, type, required, autoFocus, disabled, value, onChangeHandler, row }) => {
    return <Grid2 container spacing={8} alignItems="flex-end" style={commonStyle.formRow}>
        <Grid2 size={{ xs: 12 }}>
            <TextField name={id} id={id} label={label} type={type || "text"} disabled={disabled}
                variant="outlined" fullWidth autoFocus={autoFocus} required={required || false}
                value={value}
                multiline
                maxRows={row || 5}
                minRows={row || 5}
                onChange={onChangeHandler} />
        </Grid2>
    </Grid2>
}

export default FormTextArea