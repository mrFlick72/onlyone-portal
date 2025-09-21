import React from "react"
import {Dialog, DialogActions, DialogContent, DialogTitle} from "@mui/material";
import {AddShoppingCart} from "@mui/icons-material";
import MonthsSelector from "../MonthsSelector";
import YearSelector from "../YearSelector";
import selectUiAdapterFor from "../../search-tags/SearchTagsUIAdapter";
import YesAndNoButtonGroup from "../../../components/layout/YesAndNoButtonGroup";
import FormSelect from "../../../components/form/FormSelect";

const SearchBudgetExpensePopUp = ({
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
            <FormSelect options={selectUiAdapterFor(searchTagRegistry)} value={searchTags} multi={true} onChangeHandler={handlers.searchTag}/>
            <MonthsSelector monthRegistry={monthRegistry} month={month} handler={handlers.month}/>
            <YearSelector year={year} handler={handlers.year}/>
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
}

export default SearchBudgetExpensePopUp