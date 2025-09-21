import React from 'react';

import SpentBudgetApp from "./SpentBudgetApp";
import {resetSearchParameters} from "./SearchCriteriaProvider";
import ComponentInitializer from "../components/ComponentInitializer";
import {authenticationChecker} from "../auth/Authenticator";

resetSearchParameters()
authenticationChecker()

ComponentInitializer(<SpentBudgetApp/>,)