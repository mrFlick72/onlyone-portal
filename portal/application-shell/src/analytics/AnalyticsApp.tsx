import React, { useEffect, useMemo, useState } from "react"
import { getAllMessageRegistry } from "../messages/MessageRepository";
import { createBrowserRouter, RouterProvider } from "react-router";
import { isAuthenticated } from "../auth/Authenticator";
import AnalyticsDashboardPage from "./AnalyticsDashboardPage";

const AnalyticsApp: React.FC = () => {
    const [messageRegistry, setMessageRegistry] = useState({})

    useEffect(() => {
        isAuthenticated().then()
        setMessageRegistry(getAllMessageRegistry())
    }, []);

    const router = useMemo(() => createBrowserRouter([
        { path: "/analytics/index", element: <AnalyticsDashboardPage messageRegistry={messageRegistry} /> },
    ]), [messageRegistry]);

    return <RouterProvider router={router} />;
}

export default AnalyticsApp;
