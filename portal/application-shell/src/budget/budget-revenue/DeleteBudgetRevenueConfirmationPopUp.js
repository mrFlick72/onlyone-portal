import React from "react"
import ConfirmationPopUp from "../../components/layout/ConfirmationPopUp";

export default ({saveCallback, modal, open, handleClose}) =>
    <ConfirmationPopUp confirmationHandler={saveCallback}
                       modalId={modal.id}
                       handleClose={handleClose}
                       open={open}
                       modalMessageBody={modal.message}
                       modalTitle={modal.title}/>
