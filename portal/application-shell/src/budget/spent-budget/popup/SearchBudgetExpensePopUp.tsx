import React from "react"
import { v1 as uuidv1 } from 'uuid';    
import SpentBudgetForm from "../budget/SpentBudgetForm";       
import { Dialog, DialogActions, DialogContent, DialogTitle } from "@mui/material";
import { AddShoppingCart } from "@mui/icons-material";
import MonthsSelector from "../MonthsSelector";
import YearSelector from "../YearSelector";
import selectUiAdapterFor from "../../search-tags/SearchTagsUIAdapter";
import YesAndNoButtonGroup from "../../../components/layout/YesAndNoButtonGroup";
import FormSelect from "../../../components/form/FormSelect";
import { Month } from "../../../time/months";

type SearchBudgetExpensePopUpProps = {
    open: boolean,
    handleClose: () => void,
    modal: {
        title: string,
        saveButtonLabel: string,
        closeButtonLabel: string
    },
    month: string,
    year: string,
    searchTags: { value: string, label: string }[],
    searchTagRegistry: SearchTag[],
    monthRegistry: Month[],
    handlers: {
        month: (selectedOption: any) => void,
        year: (selectedOption: any) => void,
        searchTag: (...args: any) => void
    },
    saveCallback: () => void
}

const SearchBudgetExpensePopUp: React.FC<SearchBudgetExpensePopUpProps> = ({
    open,
    handleClose,
    modal,
    month,
    year,
    searchTags,
    searchTagRegistry,
    monthRegistry,
    handlers,
    saveCallback
}) => {
    return <Dialog onClose={handleClose} open={open} fullWidth scroll="paper">
        <DialogTitle>{modal.title}</DialogTitle>

        <DialogContent>
            <FormSelect label="" id={uuidv1()} options={selectUiAdapterFor(searchTagRegistry)} value={searchTags} multi={true} onChangeHandler={handlers.searchTag} />
            <MonthsSelector monthRegistry={monthRegistry} month={month} handler={handlers.month} />
            <YearSelector year={year} handler={handlers.year} />
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

export default SearchBudgetExpensePopUp