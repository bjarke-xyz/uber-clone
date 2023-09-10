import express from "express";
import { getAuthToken } from "./auth-middleware";
import { simManager } from "./sim-manager";

export const adminRouter = express.Router();

adminRouter.get("/test", (req, res) => {
  const authToken = getAuthToken(req);
  res.json({ hello: "world", authToken });
});

function getStatus() {
  const entities = simManager.getEntities();
  const result = entities.map((entity) => {
    return {
      email: entity.user.email,
      isRider: entity.user.isRider ?? false,
      state: entity.runner.getState(),
    };
  });
  return result;
}

adminRouter.get("/status", (req, res) => {
  const status = getStatus();
  res.json(status);
});

adminRouter.post("/stop", (req, res) => {
  simManager.stopAll();
  const status = getStatus();
  res.json(status);
});

adminRouter.post("/start", (req, res) => {
  simManager.startAll();
  const status = getStatus();
  res.json(status);
});
