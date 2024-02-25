import express from "express";
import { authMiddleware } from "./auth-middleware";
import { simManager } from "./sim-manager";

export const adminRouter = express.Router();

function getStatus() {
  const entities = simManager.getEntities();
  const result = entities.map((entity) => {
    return {
      user: entity.runner.getUser(),
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

adminRouter.post("/stop", authMiddleware("ADMIN"), (req, res) => {
  simManager.stopAll();
  const status = getStatus();
  res.json(status);
});

adminRouter.post("/start", (req, res) => {
  simManager.startAll();
  const status = getStatus();
  res.json(status);
});
