import { createTheme } from "@mui/material";
import type { } from '@mui/x-date-pickers/themeAugmentation';

export const commonStyle = {
    formRow: {
        marginTop: "10px",
        paddingTop: "10px"
    }
}

const themeProvider = createTheme({
    palette: {
        primary: {
            main: '#252624',
            contrastText: '#fff',
        },
    },
});

export default themeProvider
