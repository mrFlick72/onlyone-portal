import React from "react";
import { Button, Grid } from "@mui/material";
import { commonStyle } from "../../theme/ThemeProvider";

interface FormButtonProps {
    labelPrefix?: any,
    label: string,
    type: any,
    onClickHandler?: React.MouseEventHandler<HTMLButtonElement>,
    direction?: string
}


const FormButton: React.FC<FormButtonProps> = ({ labelPrefix, label, type, onClickHandler, direction }) => {
    return <div dir={direction || ""}>
        <Grid container alignItems="flex-end" style={commonStyle.formRow}>
            <Grid size={{ xs: 12 }}>
                <Button type={type || "button"}
                    variant="outlined"
                    color="primary"
                    onClick={onClickHandler}
                    style={{ textTransform: "none" }}>
                    {labelPrefix} {label}
                </Button>
            </Grid>
        </Grid>
    </div>
}

export default FormButton