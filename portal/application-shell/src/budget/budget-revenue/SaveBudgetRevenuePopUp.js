import React from "react"
import {Dialog, DialogActions, DialogContent, DialogTitle} from "@mui/material";
import {AddShoppingCart} from "@mui/icons-material";
import BudgetRevenueForm from "./BudgetRevenueForm";
import YesAndNoButtonGroup from "../../components/layout/YesAndNoButtonGroup";

export default ({budgetRevenue, handlers, saveCallback, modal, open, handleClose}) =>
    <Dialog onClose={handleClose} open={open} fullWidth scroll="paper">
        <DialogTitle>{modal.title}</DialogTitle>

        <DialogContent>
            <BudgetRevenueForm budgetRevenueData={budgetRevenue}
                               budgetRevenueHandlers={handlers}/>
        </DialogContent>
        <DialogActions>
            <YesAndNoButtonGroup yesIcon={<AddShoppingCart/>}
                                 yesFun={saveCallback}
                                 noFun={handleClose}
                                 buttonMessages={{
                                     "noLabel": modal.closeButtonLabel,
                                     "yesLabel": modal.saveButtonLabel
                                 }}/>
        </DialogActions>
    </Dialog>
