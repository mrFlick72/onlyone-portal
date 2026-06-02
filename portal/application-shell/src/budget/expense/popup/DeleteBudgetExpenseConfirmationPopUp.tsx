import React from "react"
import ConfirmationPopUp from "../../../components/layout/ConfirmationPopUp";
type DeleteBudgetExpenseConfirmationPopUpProps = {
    deleteBudgetExpenseAction: () => void,
    modal: {
        title: string,
        message: string,
        yesLabel: string,
        noLabel: string
    },
    open: boolean,
    handleClose: () => void
}

const DeleteBudgetExpenseConfirmationPopUp: React.FC<DeleteBudgetExpenseConfirmationPopUpProps> = ({ deleteBudgetExpenseAction, modal, open, handleClose }) =>
    <ConfirmationPopUp confirmationHandler={deleteBudgetExpenseAction}
        handleClose={handleClose}
        open={open}
        modalMessageBody={modal.message}
        modalTitle={modal.title}
        yesLabel={modal.yesLabel}
        noLabel={modal.noLabel} />

export default DeleteBudgetExpenseConfirmationPopUp