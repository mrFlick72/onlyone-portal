import React from 'react';

import PlanApp from "./PlanApp";
import ComponentInitializer from "../components/ComponentInitializer";
import { authenticationChecker } from "../auth/Authenticator";

authenticationChecker()

ComponentInitializer(<PlanApp />)
