import React from 'react';
import AccountDetailsPage from "./AccountDetailsPage";
import ComponentInitializer from "../components/ComponentInitializer";
import {authenticationChecker} from "../auth/Authenticator";

authenticationChecker()

ComponentInitializer(<AccountDetailsPage/>,)