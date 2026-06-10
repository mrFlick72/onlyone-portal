import React from 'react';

import AnalyticsApp from "./AnalyticsApp";
import ComponentInitializer from "../components/ComponentInitializer";
import { authenticationChecker } from "../auth/Authenticator";

authenticationChecker()

ComponentInitializer(<AnalyticsApp />,)
