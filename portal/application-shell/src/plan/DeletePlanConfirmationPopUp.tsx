import React from "react"
import ConfirmationPopUp from "../components/layout/ConfirmationPopUp";

interface DeletePlanConfirmationPopUpProps {
    open: boolean;
    handleClose: () => void;
    saveCallback: () => void;
    modal: {
        title: string;
        message: string;
        yesLabel: string;
        noLabel: string;
    };
}

const DeletePlanConfirmationPopUp: React.FC<DeletePlanConfirmationPopUpProps> = ({ open, handleClose, saveCallback, modal }) => (
    <ConfirmationPopUp
        open={open}
        handleClose={handleClose}
        confirmationHandler={saveCallback}
        modalTitle={modal.title}
        modalMessageBody={modal.message}
        yesLabel={modal.yesLabel}
        noLabel={modal.noLabel} />
);

export default DeletePlanConfirmationPopUp;
