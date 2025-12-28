import React, { useEffect, useState } from "react"
import { getAllMessageRegistry, MessageBundle } from "../messages/MessageRepository";
import BudgetExpensePage from "./spent-budget/BudgetExpensePage";
import { HashRouter, Route, Routes } from "react-router";
import BudgetRevenuePage from "./budget-revenue/BudgetRevenuePage";
import SearchTagsPage from "./search-tags/SearchTagsPage";
import { isAuthenticated } from "../auth/Authenticator";

interface SpentBudgetAppProps {
}

const SpentBudgetApp: React.FC<SpentBudgetAppProps> = () => {
    let [messageRegistry, setMessageRegistry] = useState({})

    useEffect(() => {
        isAuthenticated().then()
        setMessageRegistry(getAllMessageRegistry())
    }, []); // Or [] if effect doesn't need props or state



    return (
        <HashRouter>
            <Routes>
                <Route path="/"
                    element={<BudgetExpensePage messageRegistry={messageRegistry} />} />
                <Route path="/budget-revenue"
                    element={<BudgetRevenuePage messageRegistry={messageRegistry} />} />
                <Route path="/search-tags"
                    element={<SearchTagsPage messageRegistry={messageRegistry} />} />
            </Routes>
        </HashRouter>)
}

export default SpentBudgetApp;