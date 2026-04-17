import React from "react"
import ConfirmationPopUp from "../../components/layout/ConfirmationPopUp";

interface DeleteBudgetRevenueConfirmationPopUpProps {
    saveCallback: () => void;
    modal: {        
        title: string;
        message: string;
    };
    open: boolean;
    handleClose: () => void;
}

const DeleteBudgetRevenueConfirmationPopUp: React.FC<DeleteBudgetRevenueConfirmationPopUpProps> = ({ saveCallback, modal, open, handleClose }) =>
    <ConfirmationPopUp confirmationHandler={saveCallback}
        handleClose={handleClose}
        open={open}
        modalMessageBody={modal.message}
        modalTitle={modal.title} />

export default DeleteBudgetRevenueConfirmationPopUp;