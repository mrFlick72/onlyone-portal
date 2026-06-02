import React from "react"
import SpentBudgetForm from "../budget/SpentBudgetForm";
import { Dialog, DialogActions, DialogContent, DialogTitle } from "@mui/material";
import { AddShoppingCart } from "@mui/icons-material";
import selectUiAdapterFor from "../../search-tags/SearchTagsUIAdapter";
import YesAndNoButtonGroup from "../../../components/layout/YesAndNoButtonGroup";
import { SavedBudgetExpense } from "../domain/BudgetExpense";
import { SaveModalMessageBundle, SpentBudgetFormMessageBundle } from "../../../messages/MessageBundles";

type SaveBudgetExpensePopUpProps = {
    open: boolean,
    handleClose: () => void,
    modal: SaveModalMessageBundle,
    searchTagRegistry: SearchTag[],
    spentBudgetHandlers: {
        date: any,
        amount: any,
        note: any,
        searchTag: (selectedOption: any) => void

    },
    budgetExpense: SavedBudgetExpense,
    saveCallback: () => void,
    formMessages: SpentBudgetFormMessageBundle
}

const SaveBudgetExpensePopUp: React.FC<SaveBudgetExpensePopUpProps> = ({
    open,
    handleClose,
    modal,
    searchTagRegistry,
    spentBudgetHandlers,
    budgetExpense,
    saveCallback,
    formMessages
}) => {
    console.log("SaveBudgetExpensePopUp render with data: ", budgetExpense)
    return <Dialog onClose={handleClose} open={open} fullWidth scroll="paper">
        <DialogTitle>{modal.title}</DialogTitle>

        <DialogContent>
            <SpentBudgetForm spentBudgetData={{
                date: budgetExpense.date,
                amount: budgetExpense.amount,
                note: budgetExpense.note,
                searchTags: budgetExpense.searchTags
            }}
                spentBudgetHandlers={spentBudgetHandlers}
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
}

export default SaveBudgetExpensePopUp