import React from "react"
import { Dialog, DialogActions, DialogContent, DialogTitle } from "@mui/material"
import { Search } from "@mui/icons-material"
import YearSelector from "../expense/YearSelector"
import YesAndNoButtonGroup from "../../components/layout/YesAndNoButtonGroup"

type RevenueSearchPopUpProps = {
    open: boolean
    handleClose: () => void
    modal: { title: string; saveButtonLabel: string; closeButtonLabel: string }
    year: string
    handler: (event: React.ChangeEvent<HTMLInputElement>) => void
    saveCallback: () => void
}

const RevenueSearchPopUp: React.FC<RevenueSearchPopUpProps> = ({
    open, handleClose, modal, year, handler, saveCallback
}) => (
    <Dialog onClose={handleClose} open={open} fullWidth scroll="paper">
        <DialogTitle>{modal.title}</DialogTitle>
        <DialogContent>
            <YearSelector year={year} handler={handler} />
        </DialogContent>
        <DialogActions>
            <YesAndNoButtonGroup
                yesIcon={<Search />}
                yesFun={saveCallback}
                noFun={handleClose}
                buttonMessages={{ noLabel: modal.closeButtonLabel, yesLabel: modal.saveButtonLabel }} />
        </DialogActions>
    </Dialog>
)

export default RevenueSearchPopUp
