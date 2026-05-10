import React, { useEffect, useMemo, useState } from "react"
import { createBrowserRouter, RouterProvider } from "react-router";
import { getAllMessageRegistry } from "../messages/MessageRepository";
import { isAuthenticated } from "../auth/Authenticator";
import PlanListPage from "./PlanListPage";
import PlanDetailPage from "./PlanDetailPage";

const PlanApp: React.FC = () => {
    const [messageRegistry, setMessageRegistry] = useState({})

    useEffect(() => {
        isAuthenticated().then()
        setMessageRegistry(getAllMessageRegistry())
    }, []);

    const router = useMemo(() => createBrowserRouter([
        { path: "/plan/index",  element: <PlanListPage messageRegistry={messageRegistry} /> },
        { path: "/plan/detail", element: <PlanDetailPage messageRegistry={messageRegistry} /> },
    ]), [messageRegistry]);

    return <RouterProvider router={router} />;
}

export default PlanApp;
