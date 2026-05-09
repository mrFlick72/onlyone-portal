export type MessageBundle = {
    [name: string]: string
}

export function getAllMessageRegistry(): MessageBundle {
    return {
        "common.title": "Only One Portal",
        "logout.label": "Log Out",
        "form.firstName.label": "First Name:",
        "form.firstName.placeholder": "First Name",
        "form.lastName.label": "Last Name:",
        "form.lastName.placeholder": "Last Name",
        "form.birthDate.label": "Birth Date:",
        "form.phone.label": "Phone:",
        "form.phone.placeholder": "Phone",
        "form.email.label": "Mail:",
        "form.email.placeholder": "Mail",
        "form.save.value": "Save changes",
        "menu.userProfile.label": "User Profile",
        "menu.budgetPage.label": "Budget",
        "menu.revenuePage.label": "Revenue",
        "menu.logOut.label": "Log Out",
        "menu.title": "Home",
        "common.button.save.label": "Save",
        "common.button.close.label": "Close",
        "common.button.upload.label": "Upload",
        "attachment.popup.chooseFile.label": "Choose file",
        "attachment.popup.noFileSelected.label": "No file selected",
        "attachment.popup.existingAttachments.label": "Existing attachments",
        "attachment.popup.noAttachments.label": "No attachments yet",
        "attachment.popup.download.label": "Download attachment",
        "attachment.popup.delete.label": "Delete attachment",
        "budgetExpensePage.menu.newBudgetExpense": "New Budget Expense",
        "budgetExpensePage.menu.search": "Search",
        "budgetExpensePage.menu.searchTags": "Search Tags",
        "budgetExpensePage.menu.diagrams": "Diagram Chart",
        "budgetExpensePage.search.popup.title": "Filter by tag",
        "budgetExpensePage.delete.popup.title": "Delete Budget Expense Attachments",
        "budgetExpensePage.delete.popup.message": "Are you sure of delete the Budget Expense from the list?",
        "budgetExpensePage.attachment.popup.title": "Load Attachments",
        "budgetExpensePage.newBudgetExpense.popup.title": "Save a Budget Expense",
        "budgetExpensePage.loadCsvFile.popup.title": "Load data by a csv file",
        "budgetExpensePage.loadCsvFile.popup.save": "Load File...",
        "budgetExpensePage.loadCsvFile.popup.formUploadFileLabel": "File:",
        "budgetExpensePage.cards.dailyDetails": "Daily details",
        "budgetExpensePage.cards.totalByCategories": "Total by categories",
        "budgetRevenuePage.newBudgetRevenueModal.title": "Save a Revenue Expense",
        "budgetRevenuePage.menu.changeYear": "Search",
        "budgetRevenuePage.changeYear.popup.title": "Search",
        "budgetRevenuePage.attachment.popup.title": "Load Attachments",
        "menu.planPage.label": "Plans",
        "planPage.menu.newPlan": "New Plan",
        "planPage.savePlan.popup.title": "Save a Plan",
        "planPage.deletePlan.popup.title": "Delete Plan",
        "planPage.deletePlan.popup.message": "Are you sure you want to delete this plan and all its todos?",
        "planDetailPage.menu.newTodo": "Add Todo",
        "planDetailPage.menu.backToList": "Back to plans",
        "planDetailPage.saveTodo.popup.title": "Add a Todo",
        "planDetailPage.updateTodo.popup.title": "Edit Todo",
        "planDetailPage.deleteTodo.popup.title": "Delete Todo",
        "planDetailPage.deleteTodo.popup.message": "Are you sure you want to delete this todo?"
    }
}

export function getMessageFor(bundle: MessageBundle, key: string): string {
    return bundle[key] || "";
}
