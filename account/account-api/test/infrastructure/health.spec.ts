const request = require("supertest");
import {actuatorApp} from "../../app/app";

test('GET /health', async () => {
    let received = await request(actuatorApp)
        .get("/health")
        .send()
    expect(received.status).toBe(200);
});