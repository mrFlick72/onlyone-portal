import React, { ReactNode } from "react"
import { Box, Button } from "@mui/material";
import { Close } from "@mui/icons-material";

type YesAndNoButtonGroupProps = {
    buttonMessages: {[name:string]:string}
    yesIcon: ReactNode
    yesFun: () => void
    noFun: () => void

}

const YesAndNoButtonGroup: React.FC<YesAndNoButtonGroupProps> = ({ buttonMessages, yesIcon, yesFun, noFun }) => {
    return <Box>
        <Button variant="contained" onClick={noFun} color="primary">
            <Close /> {buttonMessages.noLabel}
        </Button>
        <Button variant="contained" onClick={yesFun} color="success">
            {yesIcon} {buttonMessages.yesLabel}
        </Button>
    </Box>
}


export default YesAndNoButtonGroup