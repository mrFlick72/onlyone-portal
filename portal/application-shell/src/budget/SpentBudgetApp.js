import React, {useEffect, useState} from "react"
import {getAllMessageRegistry} from "../messages/MessageRepository";
import BudgetExpensePage from "./spent-budget/BudgetExpensePage";
import {HashRouter, Route, Routes} from "react-router";
import BudgetRevenuePage from "./budget-revenue/BudgetRevenuePage";
import SearchTagsPage from "./search-tags/SearchTagsPage";
import {isAuthenticated} from "../auth/Authenticator";


export default () => {
    let [messageRegistry, setMessageRegistry] = useState({})

    useEffect(() => {
        isAuthenticated().then()
        setMessageRegistry(getAllMessageRegistry())
    }, []); // Or [] if effect doesn't need props or state



    return (
        <HashRouter>
            <Routes>
                <Route exact={true} path="/"
                       element={<BudgetExpensePage messageRegistry={messageRegistry}/>}/>
                <Route exact={true} path="/budget-revenue"
                       element={<BudgetRevenuePage  messageRegistry={messageRegistry}/>}/>
                <Route path="/search-tags" exact={true}
                       element={<SearchTagsPage messageRegistry={messageRegistry}/>}/>
            </Routes>
        </HashRouter>)
}
