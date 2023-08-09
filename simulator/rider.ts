import { readFile } from "fs/promises";
import path from "path";
import { BackendApi } from "./api";
import { OSMAPI } from "./openstreetmap";
import {
  CityData,
  LatLng,
  NamedPoint,
  RideRequest,
  RideRequestState,
  cityDataFeatureToNamedPoints,
} from "./types";
import { randomIntFromInterval, wait } from "./util";
import { sample } from "lodash";

export class SimRider {
  private currentRideRequest: RideRequest | null = null;

  private currentFrom: NamedPoint | null = null;
  private currentTo: NamedPoint | null = null;

  private prevFrom: NamedPoint | null = null;
  private prevTo: NamedPoint | null = null;

  constructor(
    private api: BackendApi,
    private osm: OSMAPI,
    private userEmail: string,
    private userPassword: string,
    private city: string
  ) {}
  private async log(message?: any, ...optionalParams: any[]) {
    console.log(`${this.userEmail} [R] | ${message}`, ...optionalParams);
    await this.api.postLog({ tag: "R", message });
  }

  public async run() {
    this.selectRandomPoints();

    await this.api.signIn(this.userEmail, this.userPassword);
    const myAvailableRides = (await this.api.getMyRides()).filter(
      (x) =>
        x.state === RideRequestState.Available ||
        x.state === RideRequestState.Accepted ||
        x.state == RideRequestState.InProgress
    );
    if (myAvailableRides.length > 0) {
      this.currentRideRequest = myAvailableRides[0];
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

    while (true) {
      if (
        !this.currentRideRequest ||
        this.currentRideRequest.state === RideRequestState.Finished
      ) {
        const randomPoints = await this.selectRandomPoints(this.currentTo, 2);
        if (randomPoints.length !== 2) {
          console.log("failed to get 2 random points, waiting 10s...");
          await wait(10 * 1000);
          continue;
        }
        const [from, to] = randomPoints;
        this.prevFrom = this.currentFrom;
        this.prevTo = this.currentTo;
        this.currentFrom = from;
        this.currentTo = to;
        await this.requestRide(from, to);
      }
      const randomWait = randomIntFromInterval(5, 30);
      await wait(randomWait * 1000);
      await this.updateRide();
    }
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
      const resource = `${this.city}_clean.geojson`;
      const url = path.resolve(
        __dirname,
        `./data/random-city-data/${resource}`
      );
      const file = await readFile(url, { encoding: "utf-8" });
      const cityData = JSON.parse(file) as CityData;
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
