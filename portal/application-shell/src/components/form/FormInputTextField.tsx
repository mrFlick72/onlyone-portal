import { Grid, TextField } from "@mui/material";
import React from "react";
import { commonStyle } from "../../theme/ThemeProvider";

type FormInputTextFieldProps = {
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
    return <Grid container spacing={8} sx={{ alignItems: "flex-end", ...commonStyle.formRow }}>
        {prefix && <Grid>

            {prefix}
        </Grid>}
        <Grid size={{ xs: 12 }}>
            <TextField name={id} id={id} label={label} type={type || "text"} disabled={disabled}
                variant="outlined" fullWidth autoFocus={autoFocus} required={required || false}
                value={value}
                onChange={handler} />
        </Grid>
    </Grid>
}

export default FormInputTextField