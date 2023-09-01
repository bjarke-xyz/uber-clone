import { sample } from "lodash";
import { BackendApi } from "./api";
import {
  LatLng,
  NamedPoint,
  RideRequest,
  RideRequestState,
  cityDataFeatureToNamedPoints,
} from "./types";
import { getCityData, randomIntFromInterval, wait } from "./util";

export class SimRider {
  private currentRideRequest: RideRequest | null = null;

  private currentFrom: NamedPoint | null = null;
  private currentTo: NamedPoint | null = null;

  constructor(
    private api: BackendApi,
    private userEmail: string,
    private userPassword: string,
    private city: string
  ) {}
  private async log(message?: any, ...optionalParams: any[]) {
    console.log(`${this.userEmail} [R] | ${message}`, ...optionalParams);
    await this.api.postLog({ tag: "R", message });
  }

  public async run() {
    try {
      await this.api.signIn(this.userEmail, this.userPassword);
      const myAvailableRides = await this.getAvailableRides();
      if (myAvailableRides.length > 0) {
        await this.initializeExistingRideRequest(myAvailableRides);
      }

      while (true) {
        if (this.needsRide()) {
          const randomPoints = await this.selectRandomPoints(this.currentTo, 2);
          if (randomPoints.length !== 2) {
            console.log("failed to get 2 random points, waiting 10s...");
            await wait(10 * 1000);
            continue;
          }
          [this.currentFrom, this.currentTo] = randomPoints;
          await this.requestRide(this.currentFrom, this.currentTo);
        }
        const randomWait = randomIntFromInterval(5, 30);
        await wait(randomWait * 1000);
        await this.updateRide();
      }
    } catch (error) {
      console.error("Unexpected error in rider run", error);
    }
  }

  private needsRide(): boolean {
    return (
      !this.currentRideRequest ||
      this.currentRideRequest.state === RideRequestState.Finished
    );
  }

  private async initializeExistingRideRequest(
    rideRequests: RideRequest[]
  ): Promise<void> {
    if (rideRequests.length === 0) {
      return;
    }
    this.currentRideRequest = rideRequests[0];
    this.currentFrom = {
      location: new LatLng(
        this.currentRideRequest.fromLat,
        this.currentRideRequest.fromLng
      ),
      name: this.currentRideRequest.fromName,
    };
    this.currentTo = {
      location: new LatLng(
        this.currentRideRequest.toLat,
        this.currentRideRequest.toLng
      ),
      name: this.currentRideRequest.toName,
    };
    await this.log(
      `Found existing rider requested ride ${this.currentRideRequest.id}`
    );
  }

  private async getAvailableRides() {
    return (await this.api.getMyRides()).filter((x) =>
      [
        RideRequestState.Available,
        RideRequestState.Accepted,
        RideRequestState.InProgress,
      ].includes(x.state)
    );
  }

  private async updateRide() {
    if (!this.currentRideRequest) {
      return;
    }
    const updatedCurrentRide = await this.api.getRideRequest(
      this.currentRideRequest.id
    );
    if (updatedCurrentRide) {
      this.currentRideRequest = updatedCurrentRide;
    }
  }

  private async requestRide(from: NamedPoint, to: NamedPoint) {
    await this.log(`Requested ride ${from.name} -> ${to.name}`);
    this.currentRideRequest = await this.api.createRideRequest(
      from.location,
      from.name,
      to.location,
      to.name
    );
  }

  private async selectRandomPoints(
    startPoint: NamedPoint | null = null,
    numberOfPoints = 2
  ): Promise<NamedPoint[]> {
    if (numberOfPoints <= 0) return [];
    try {
      const cityData = await getCityData(this.city);
      if (cityData.features.length === 0) {
        return [];
      }
      if (!startPoint) {
        const p1 = sample(cityData.features);
        if (!p1) {
          return [];
        }
        startPoint = cityDataFeatureToNamedPoints(p1);
      }
      const points: NamedPoint[] = [startPoint];
      const minDistanceMeters = 1000;
      for (let i = 0; i < numberOfPoints - 1; i++) {
        for (let j = 0; j < 10; j++) {
          // Try 10 times finding the next point
          const nextPoint = sample(cityData.features);
          if (nextPoint) {
            const nextNamedPoint = cityDataFeatureToNamedPoints(nextPoint);
            const distanceToPrevPoint = points[
              points.length - 1
            ].location.distanceTo(nextNamedPoint.location);
            if (distanceToPrevPoint < minDistanceMeters) {
              continue;
            } else {
              points.push(nextNamedPoint);
              break;
            }
          }
        }
      }
      return points;
    } catch (error) {
      console.error("failed to select random points", error);
      return [];
    }
  }
}
