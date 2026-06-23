import React from "react"
import { Dialog, DialogActions, DialogContent, DialogTitle } from "@mui/material";
import { Check } from "@mui/icons-material";
import FormInputTextField from "../../components/form/FormInputTextField";
import YesAndNoButtonGroup from "../../components/layout/YesAndNoButtonGroup";
import { SaveModalMessageBundle, SearchTagFormMessageBundle } from "../../messages/MessageBundles";

type UpdateSearchTagPopUpProps = {
    open: boolean;
    handleClose: () => void;
    value: string;
    valueHandler: (event: React.ChangeEvent<HTMLInputElement>) => void;
    saveCallback: () => void;
    modal: SaveModalMessageBundle;
    formMessages: SearchTagFormMessageBundle;
}

// Edit is value-only: Key and Scope are immutable once a tag exists (see
// tag-api docs/adr/0008-tag-update-is-value-only.md), and Scope is already
// intrinsic to the tab the user is on, so neither is shown here.
const UpdateSearchTagPopUp: React.FC<UpdateSearchTagPopUpProps> = ({ open, handleClose, value, valueHandler, saveCallback, modal, formMessages }) => (
    <Dialog onClose={handleClose} open={open} fullWidth scroll="paper">
        <DialogTitle>{modal.title}</DialogTitle>
        <DialogContent>
            <FormInputTextField
                handler={valueHandler}
                autoFocus={true}
                value={value}
                id="updateSearchTagValue"
                label={formMessages.valueLabel} />
        </DialogContent>
        <DialogActions>
            <YesAndNoButtonGroup
                yesIcon={<Check />}
                yesFun={saveCallback}
                noFun={handleClose}
                buttonMessages={{
                    noLabel: modal.closeButtonLabel,
                    yesLabel: modal.saveButtonLabel,
                }} />
        </DialogActions>
    </Dialog>
);

export default UpdateSearchTagPopUp;
