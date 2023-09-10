import { BackendApiClient } from "./api-client";
import { SimDriver } from "./driver";
import { SimRider } from "./rider";
import { SimRunner, SimUser } from "./types";
import { urls } from "./util";

export interface SimEntity {
  abortController: AbortController;
  user: SimUser;
  runner: SimRunner;
}
class SimManager {
  private entities: SimEntity[] = [];

  public getEntities() {
    return this.entities;
  }
  loadUsers() {
    const simUsersStr = process.env.SIM_USERS ?? "[]";
    const simUsers = JSON.parse(simUsersStr) as SimUser[];
    this.entities = simUsers.map((user) => {
      const backendApiClient = new BackendApiClient(urls.backendApiBaseUrl);
      const abortController = new AbortController();
      const runner = user.isRider
        ? new SimRider(
            abortController,
            backendApiClient,
            user.email,
            user.password,
            "R",
            user.city
          )
        : new SimDriver(
            abortController,
            backendApiClient,
            user.email,
            user.password,
            "D"
          );
      return {
        abortController,
        user,
        runner,
      };
    });
  }

  startAll() {
    this.entities.forEach((entity) => entity.runner.run());
  }
  stopAll() {
    for (const entity of this.entities) {
      entity.runner.stop();

      const newAbortController = new AbortController();
      entity.abortController = newAbortController;
      entity.runner.setAbortController(newAbortController);
    }
  }
  public setTimeMultiplier(timeMultiplier: number) {
    this.entities.forEach((entity) =>
      entity.runner.setTimeMultiplier(timeMultiplier)
    );
  }
}

export const simManager = new SimManager();
