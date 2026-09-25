import React, { useEffect, useMemo, useState } from "react"
import { getAllMessageRegistry } from "../messages/MessageRepository";
import BudgetExpensePage from "./expense/BudgetExpensePage";
import { createBrowserRouter, RouterProvider } from "react-router";
import BudgetRevenuePage from "./revenue/BudgetRevenuePage";
import SearchTagsPage from "./search-tags/SearchTagsPage";
import ScheduledExpenseListPage from "./scheduled-expense/ScheduledExpenseListPage";
import ScheduledExpenseDetailPage from "./scheduled-expense/ScheduledExpenseDetailPage";
import { isAuthenticated } from "../auth/Authenticator";

const SpentBudgetApp: React.FC = () => {
    const [messageRegistry, setMessageRegistry] = useState({})

    useEffect(() => {
        isAuthenticated().then()
        setMessageRegistry(getAllMessageRegistry())
    }, []);

    const router = useMemo(() => createBrowserRouter([
        { path: "/budget/expense/index",     element: <BudgetExpensePage messageRegistry={messageRegistry} /> },
        { path: "/budget/revenue/index",     element: <BudgetRevenuePage messageRegistry={messageRegistry} /> },
        { path: "/budget/search-tags/index", element: <SearchTagsPage messageRegistry={messageRegistry} /> },
        { path: "/budget/scheduled-expense/index", element: <ScheduledExpenseListPage messageRegistry={messageRegistry} /> },
        { path: "/budget/scheduled-expense/detail", element: <ScheduledExpenseDetailPage messageRegistry={messageRegistry} /> },
    ]), [messageRegistry]);

    return <RouterProvider router={router} />;
}

export default SpentBudgetApp;
