import { Position, along, lineDistance, lineString } from "@turf/turf";
import { BackendApi } from "./api";
import {
  BackendUser,
  LatLng,
  RideRequest,
  RideRequestState,
  Vehicle,
  getAngleInDegrees,
} from "./types";
import { decodePolyline, randomIntFromInterval, wait } from "./util";

export interface RouteStep {
  bearing: number;
  distance: number;
  duration: number;
  locations: LatLng[];
}

export class SimDriver {
  // TODO: Adjustable dynamically
  private timeMultiplier = 16;

  private user: BackendUser | null = null;

  private currentLocation: LatLng | null = null;

  constructor(
    private api: BackendApi,
    private userEmail: string,
    private userPassword: string
  ) {}

  private async log(message?: any, ...optionalParams: any[]) {
    console.log(`${this.userEmail} [D] | ${message}`, ...optionalParams);
    await this.api.postLog({ tag: "D", message });
  }

  public async run() {
    try {
      await this.api.signIn(this.userEmail, this.userPassword);
      const vehicle = await this.api.getVehicle();
      this.user = await this.api.getMyUser();
      if (!this.user) {
        await this.log("user not found");
        return;
      }
      if (!vehicle) {
        await this.log(`no vehicle found`);
        return;
      }
      if (vehicle.lastRecordedPosition) {
        this.currentLocation = new LatLng(
          vehicle.lastRecordedPosition.lat,
          vehicle.lastRecordedPosition.lng
        );
      }
      while (true) {
        const randomWait = randomIntFromInterval(5, 15);
        await wait(randomWait * 1000);
        await this.drive(vehicle);
      }
    } catch (error) {
      console.error("Unexpected error in driver run", error);
    }
  }

  private async getMyInProgressRides(): Promise<RideRequest[]> {
    const rides = (await this.api.getMyRides()).filter(
      (x) =>
        x.driverId === this.user?.id &&
        (x.state === RideRequestState.Accepted ||
          x.state === RideRequestState.InProgress)
    );
    return rides;
  }

  private async getRideRequest(): Promise<{
    rideRequest: RideRequest | null;
    claimed: boolean;
  }> {
    const inProgressRides = await this.getMyInProgressRides();
    let rideRequest: RideRequest | null = null;
    let claimed = false;
    if (inProgressRides.length > 0) {
      rideRequest = inProgressRides[0];
      await this.log(`Found in-progress driver ride request ${rideRequest.id}`);
    } else {
      await wait(randomIntFromInterval(1, 5) * 1000);
      const availableRideRequests = await this.api.getAvailableRideRequests();
      if (availableRideRequests.length > 0) {
        const potentialRideRequest = availableRideRequests[0];
        await this.log(`Claiming ${potentialRideRequest.id}`);
        claimed = await this.api.claimRideRequest(potentialRideRequest.id);
        if (claimed) {
          rideRequest = potentialRideRequest;
        } else {
          await this.log(
            `failed to claim ride request ${potentialRideRequest.id}, waiting 10s`
          );
          await wait(10 * 1000);
        }
      }
    }
    return { rideRequest, claimed };
  }

  private async drive(vehicle: Vehicle) {
    const { rideRequest, claimed } = await this.getRideRequest();
    if (rideRequest) {
      let _currentLocation: LatLng | null = null;
      if (claimed && this.currentLocation) {
        _currentLocation = this.currentLocation;
      }
      const steps = await this.getDirections(rideRequest.id, _currentLocation);
      await this.log(
        `Driving (${rideRequest.id}) ${
          _currentLocation ? `${_currentLocation.toString()} ->` : ""
        } ${rideRequest.fromName} -> ${rideRequest.toName}`
      );
      await this.move(vehicle, steps);
      await this.api.finishRideRequest(rideRequest.id);
      await this.log("Finished ride, sleeping 25 seconds");
      await wait(15 * 1000);
    } else {
      await this.log("no ride requests found, waiting 10s");
      await wait(10 * 1000);
    }
  }

  private async move(vehicle: Vehicle, steps: RouteStep[]): Promise<void> {
    let prevLocation: LatLng | null = null;
    let bearing = 0;
    for (let i = 0; i < steps.length; i++) {
      const step = steps[i];
      bearing = step.bearing;
      const secondsPerMeter = step.duration / step.distance;
      for (let j = 0; j < step.locations.length; j++) {
        const location = step.locations[j];
        if (prevLocation) {
          const distance = prevLocation.distanceTo(location);
          const timeSeconds = secondsPerMeter * distance;
          const waitSeconds = timeSeconds / this.timeMultiplier;
          const waitMs = waitSeconds * 1000;
          const nextLocation = step.locations[j + 1];
          if (nextLocation) {
            bearing = getAngleInDegrees(location, nextLocation);
          }
          const speedKmh = (distance / timeSeconds) * 3.6;
          await this.api.updateLocation(
            vehicle.ID,
            location.lat,
            location.lng,
            bearing,
            speedKmh
          );
          await wait(waitMs);
        } else {
          await this.api.updateLocation(
            vehicle.ID,
            location.lat,
            location.lng,
            bearing,
            0
          );
        }
        this.currentLocation = location;
        prevLocation = location;
      }
    }
  }

  private async getDirections(
    rideRequestId: number,
    startPoint: LatLng | null
  ): Promise<RouteStep[]> {
    let directions = await this.api.getDirectionsV1(rideRequestId, startPoint);
    while (!directions) {
      const sleepSeconds = randomIntFromInterval(30 * 1000, 90 * 1000);
      await this.log(`failed to get directions, sleeping ${sleepSeconds}s`);
      await wait(sleepSeconds);
      directions = await this.api.getDirectionsV1(rideRequestId, startPoint);
    }
    const steps: RouteStep[] = [];
    if (directions) {
      for (const route of directions.routes) {
        const geometryPolyline = decodePolyline(route.geometry);
        for (const segment of route.segments) {
          for (const step of segment.steps) {
            const stepInfo: RouteStep = {
              bearing: step.maneuver.bearing_after,
              distance: step.distance,
              duration: step.duration,
              locations: [],
            };
            // stepInfo.locations.push(
            //   new LatLng(step.maneuver.location[1], step.maneuver.location[0])
            // );
            steps.push(stepInfo);
            // const coordinates = geometryPolyline.filter((x, i) =>
            //   step.way_points.includes(i)
            // );
            let coordinates: number[][] = [];
            if (step.way_points.length === 2) {
              const [startIndex, endIndex] = step.way_points;
              for (let i = 0; i < geometryPolyline.length; i++) {
                const coord = geometryPolyline[i];
                if (i >= startIndex && i <= endIndex) {
                  coordinates.push(coord);
                }
              }
            }
            for (const coord of coordinates) {
              const position = new LatLng(coord[0], coord[1]);
              if (!stepInfo.locations.some((x) => x.equals(position))) {
                stepInfo.locations.push(position);
              }
            }
          }
        }
      }
    }

    generatePoints(steps);
    return steps;
  }
}

function generatePoints(steps: RouteStep[]) {
  for (const step of steps) {
    if (step.locations.length >= 2) {
      const positions = step.locations.map((x) => x.toPosition());
      const locationsLineString = lineString(positions);
      const length = lineDistance(locationsLineString, { units: "meters" });
      const distance = Math.floor(length);
      const newPoints: Position[] = [positions[0]];
      let metersInterval = 100;
      if (distance < 100) {
        metersInterval = 10;
      }
      for (
        let step = metersInterval;
        step < distance + 1;
        step = step + metersInterval
      ) {
        const newPoint = along(locationsLineString, step, { units: "meters" });
        newPoints.push(newPoint.geometry.coordinates);
      }
      if (newPoints.length > 0) {
        step.locations = newPoints.map((x) => new LatLng(x[0], x[1]));
      }
    } else {
      // TODO: grab end of prev/ start of next to increase number of locations
    }
  }
}
