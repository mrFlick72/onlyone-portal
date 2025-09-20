import { actuatorApp, actuatorPort, app, port } from "./app";

app.listen(port, "0.0.0.0", () => {
    console.log(`⚡️[server]: Server is running at http://localhost:${port}`);
});

actuatorApp.listen(actuatorPort, "0.0.0.0", () => {
    console.log(`⚡️[server]: Server is running at http://localhost:${actuatorPort}`);
});