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
import { simManager } from "./sim-manager";
import { urls } from "./util";

dotenv.config();
urls.load();

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
  console.log(`using ${urls.backendApiBaseUrl} as backend API`);
});

simManager.loadUsers();
simManager.startAll();
