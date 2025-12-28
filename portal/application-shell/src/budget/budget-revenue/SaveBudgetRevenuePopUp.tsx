import React from "react"
import { Dialog, DialogActions, DialogContent, DialogTitle } from "@mui/material";
import { AddShoppingCart } from "@mui/icons-material";
import BudgetRevenueForm from "./BudgetRevenueForm";
import YesAndNoButtonGroup from "../../components/layout/YesAndNoButtonGroup";
import BudgetRevenue from "./domain/BudbetRevenue";

interface SaveBudgetRevenuePopUpProps {
    budgetRevenue: BudgetRevenue;
    handlers: {
        date: (date: moment.Moment) => void;
        amount: (event: React.ChangeEvent<HTMLInputElement>) => void;
        note: (event: React.ChangeEvent<HTMLTextAreaElement>) => void;
    };
    saveCallback: () => void;
    modal: {
        title: string;
        saveButtonLabel: string;
        closeButtonLabel: string;
    };
    open: boolean;
    handleClose: () => void;
}

const SaveBudgetRevenuePopUp: React.FC<SaveBudgetRevenuePopUpProps> = ({ budgetRevenue, handlers, saveCallback, modal, open, handleClose }) =>
    <Dialog onClose={handleClose} open={open} fullWidth scroll="paper">
        <DialogTitle>{modal.title}</DialogTitle>

        <DialogContent>
            <BudgetRevenueForm budgetRevenueData={budgetRevenue}
                budgetRevenueHandlers={handlers} />
        </DialogContent>
        <DialogActions>
            <YesAndNoButtonGroup yesIcon={<AddShoppingCart />}
                yesFun={saveCallback}
                noFun={handleClose}
                buttonMessages={{
                    "noLabel": modal.closeButtonLabel,
                    "yesLabel": modal.saveButtonLabel
                }} />
        </DialogActions>
    </Dialog>


export default SaveBudgetRevenuePopUp;  