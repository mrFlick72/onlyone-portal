import React from "react";
import { Button, Grid2 } from "@mui/material";
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
        <Grid2 container alignItems="flex-end" style={commonStyle.formRow}>
            <Grid2 size={{ xs: 12 }}>
                <Button type={type || "button"}
                    variant="outlined"
                    color="primary"
                    onClick={onClickHandler}
                    style={{ textTransform: "none" }}>
                    {labelPrefix} {label}
                </Button>
            </Grid2>
        </Grid2>
    </div>
}

export default FormButton