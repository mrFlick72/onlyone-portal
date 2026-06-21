import React from "react"
import { Dialog, DialogActions, DialogContent, DialogTitle } from "@mui/material";
import BudgetRevenueForm from "./BudgetRevenueForm";
import YesAndNoButtonGroup from "../../components/layout/YesAndNoButtonGroup";
import selectUiAdapterFor from "../search-tags/SearchTagsUIAdapter";
import { SelectOption } from "../../components/form/FormSelect";
import { AddShoppingCart } from "@mui/icons-material";
import { BudgetRevenueFormMessageBundle, SaveModalMessageBundle } from "../../messages/MessageBundles";

interface SaveBudgetRevenuePopUpProps {
    budgetRevenue: {
        date: string;
        amount: string;
        note: string;
        searchTags: SelectOption[];
    };
    handlers: {
        date: (date: moment.Moment) => void;
        amount: (event: React.ChangeEvent<HTMLInputElement>) => void;
        note: (event: React.ChangeEvent<HTMLTextAreaElement>) => void;
        searchTag: (selectedOptions: SelectOption[]) => void;
    };
    searchTagRegistry: SearchTag[];
    saveCallback: () => void;
    modal: SaveModalMessageBundle;
    formMessages: BudgetRevenueFormMessageBundle;
    open: boolean;
    handleClose: () => void;
}

const SaveBudgetRevenuePopUp: React.FC<SaveBudgetRevenuePopUpProps> = ({ budgetRevenue, handlers, searchTagRegistry, saveCallback, modal, formMessages, open, handleClose }) =>
    <Dialog onClose={handleClose} open={open} fullWidth scroll="paper">
        <DialogTitle>{modal.title}</DialogTitle>

        <DialogContent>
            <BudgetRevenueForm budgetRevenueData={budgetRevenue}
                budgetRevenueHandlers={handlers}
                searchTagRegistry={selectUiAdapterFor(searchTagRegistry)}
                messages={formMessages} />
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