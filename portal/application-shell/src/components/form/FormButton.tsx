import React from "react";
import { Button, Grid } from "@mui/material";
import { commonStyle } from "../../theme/ThemeProvider";

type FormButtonProps = {
    labelPrefix?: React.ReactNode,
    label: string,
    type: "button" | "submit" | "reset",
    onClickHandler?: React.MouseEventHandler<HTMLButtonElement>,
    direction?: string
}


const FormButton: React.FC<FormButtonProps> = ({ labelPrefix, label, type, onClickHandler, direction }) => {
    return <div dir={direction || ""}>
        <Grid container sx={{ alignItems: "flex-end", ...commonStyle.formRow }}>
            <Grid size={{ xs: 12 }}>
                <Button type={type}
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