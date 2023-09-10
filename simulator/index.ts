import cors from "cors";
import dotenv from "dotenv";
import express, { Express, Request, Response } from "express";
import morgan from "morgan";
import "source-map-support/register";
import { SimDriver } from "./driver";
import { BackendApiClient } from "./api-client";
import { SimUser } from "./types";
import { SimRider } from "./rider";
import { AuthClient } from "./proto-gen/proto/auth";
import { credentials } from "@grpc/grpc-js";
import { adminRouter } from "./admin-api";
import { authMiddleware } from "./auth-middleware";

dotenv.config();

const backendApiBaseUrl = process.env.API_BASE_URL ?? "Missing API_BASE_URL";
const app: Express = express();
const port = process.env.PORT || 3000;

app.use(cors());
app.use(morgan("tiny"));

app.get("/", async (req: Request, res: Response) => {
  res.send("Express + TypeScript Server");
});

app.use("/api/admin", authMiddleware("ADMIN"), adminRouter);

app.listen(port, () => {
  console.log(`⚡️[server]: Server is running at http://localhost:${port}`);
  console.log(`using ${backendApiBaseUrl} as backend API`);
});

const simUsersStr = process.env.SIM_USERS ?? "[]";
const simUsers = JSON.parse(simUsersStr) as SimUser[];

const simDrivers = simUsers
  .filter((u) => !u.isRider)
  .map((user) => {
    const backendApiClient = new BackendApiClient(backendApiBaseUrl);
    return new SimDriver(backendApiClient, user.email, user.password);
  });
simDrivers.forEach((driver) => driver.run());

const simRiders = simUsers
  .filter((u) => u.isRider)
  .map((user) => {
    const backendApiClient = new BackendApiClient(backendApiBaseUrl);
    return new SimRider(backendApiClient, user.email, user.password, user.city);
  });
simRiders.forEach((rider) => rider.run());
