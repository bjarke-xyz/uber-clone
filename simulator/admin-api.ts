import express from "express";
import { getAuthToken } from "./auth-middleware";

export const adminRouter = express.Router();

adminRouter.get("/test", (req, res) => {
  const authToken = getAuthToken(req);
  res.json({ hello: "world", authToken });
});
