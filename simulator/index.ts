import cors from "cors";
import dotenv from "dotenv";
import express, { Express, Request, Response } from "express";
import morgan from "morgan";
// import "source-map-support/register";
import { adminRouter } from "./admin-api";
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

app.use("/api/admin", adminRouter);

app.listen(port, () => {
  console.log(`⚡️[server]: Server is running at http://localhost:${port}`);
  console.log(`using ${urls.backendApiBaseUrl} as backend API`);
});

simManager.loadUsers();
