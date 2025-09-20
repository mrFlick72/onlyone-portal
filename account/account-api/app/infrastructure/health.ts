import { Express, Request, Response } from "express";

export function registerHealthEndPointFor(app: Express) {
    app.get('/health', (req: Request, res: Response) => {
        res.status(200).end()
    });
}