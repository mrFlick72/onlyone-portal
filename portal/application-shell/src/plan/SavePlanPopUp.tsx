import React from "react"
import { Dialog, DialogActions, DialogContent, DialogTitle } from "@mui/material";
import { PlaylistAdd } from "@mui/icons-material";
import moment from "moment";
import PlanForm, { PlanFormData } from "./PlanForm";
import YesAndNoButtonGroup from "../components/layout/YesAndNoButtonGroup";

interface SavePlanPopUpProps {
    open: boolean;
    handleClose: () => void;
    plan: PlanFormData;
    handlers: {
        title: (event: React.ChangeEvent<HTMLInputElement>) => void;
        date: (date: moment.Moment) => void;
    };
    saveCallback: () => void;
    modal: {
        title: string;
        saveButtonLabel: string;
        closeButtonLabel: string;
    };
}

const SavePlanPopUp: React.FC<SavePlanPopUpProps> = ({ open, handleClose, plan, handlers, saveCallback, modal }) => (
    <Dialog onClose={handleClose} open={open} fullWidth scroll="paper">
        <DialogTitle>{modal.title}</DialogTitle>
        <DialogContent>
            <PlanForm data={plan} handlers={handlers} />
        </DialogContent>
        <DialogActions>
            <YesAndNoButtonGroup
                yesIcon={<PlaylistAdd />}
                yesFun={saveCallback}
                noFun={handleClose}
                buttonMessages={{
                    noLabel: modal.closeButtonLabel,
                    yesLabel: modal.saveButtonLabel,
                }} />
        </DialogActions>
    </Dialog>
);

export default SavePlanPopUp;
