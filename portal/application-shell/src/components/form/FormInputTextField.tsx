import { Grid2, TextField } from "@mui/material";
import React from "react";
import { commonStyle } from "../../theme/ThemeProvider";

interface FormInputTextFieldProps {
    id: string,
    label: string,
    type?: string,
    required?: boolean,
    autoFocus?: boolean,
    disabled?: boolean,
    prefix?: string,
    value: string,
    handler: (input: any) => void
}

const FormInputTextField: React.FC<FormInputTextFieldProps> = ({
    id,
    label,
    type,
    required,
    autoFocus,
    disabled,
    prefix,
    value,
    handler
}) => {
    return <Grid2 container spacing={8} alignItems="flex-end" style={commonStyle.formRow}>
        {prefix && <Grid2>

            {prefix}
        </Grid2>}
        <Grid2 size={{ xs: 12 }}>
            <TextField name={id} id={id} label={label} type={type || "text"} disabled={disabled}
                variant="outlined" fullWidth autoFocus={autoFocus} required={required || false}
                value={value}
                onChange={handler} />
        </Grid2>
    </Grid2>
}

export default FormInputTextField