import { Position, along, lineDistance, lineString } from "@turf/turf";
import { BackendApiClient } from "./api-client";
import {
  BackendUser,
  LatLng,
  RideRequest,
  RideRequestState,
  SimRunner,
  Vehicle,
  getAngleInDegrees,
  isAbortError,
} from "./types";
import { decodePolyline, randomIntFromInterval } from "./util";

export interface RouteStep {
  bearing: number;
  distance: number;
  duration: number;
  locations: LatLng[];
}

export class SimDriver extends SimRunner {
  private user: BackendUser | null = null;
  private currentLocation: LatLng | null = null;

  public async run() {
    if (this.running) {
      return;
    }
    try {
      await this.apiClient.signIn(this.userEmail, this.userPassword);
      this.starting();
      const vehicle = await this.apiClient.getVehicle();
      this.user = await this.apiClient.getMyUser();
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
      this.started();
      while (this.running) {
        const randomWait = randomIntFromInterval(5, 15);
        await this.wait(randomWait * 1000);
        await this.drive(vehicle);
      }
    } catch (error) {
      if (!isAbortError(error)) {
        console.error("Unexpected error in driver run", error);
      }
    } finally {
      await this.stopped();
    }
  }

  private async getMyInProgressRides(): Promise<RideRequest[]> {
    const rides = (await this.apiClient.getMyRides()).filter(
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
      await this.wait(randomIntFromInterval(1, 5) * 1000);
      const availableRideRequests =
        await this.apiClient.getAvailableRideRequests();
      if (availableRideRequests.length > 0) {
        const potentialRideRequest = availableRideRequests[0];
        await this.log(`Claiming ${potentialRideRequest.id}`);
        claimed = await this.apiClient.claimRideRequest(
          potentialRideRequest.id
        );
        if (claimed) {
          rideRequest = potentialRideRequest;
        } else {
          await this.log(
            `failed to claim ride request ${potentialRideRequest.id}, waiting 10s`
          );
          await this.wait(10 * 1000);
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
      const finished = await this.move(vehicle, steps);
      if (finished) {
        await this.apiClient.finishRideRequest(rideRequest.id);
        await this.log("Finished ride, sleeping 25 seconds");
        await this.wait(15 * 1000);
      }
    } else {
      await this.log("no ride requests found, waiting 10s");
      await this.wait(10 * 1000);
    }
  }

  private async move(vehicle: Vehicle, steps: RouteStep[]): Promise<boolean> {
    let prevLocation: LatLng | null = null;
    let bearing = 0;
    for (let i = 0; i < steps.length; i++) {
      const step = steps[i];
      bearing = step.bearing;
      const secondsPerMeter = step.duration / step.distance;
      for (let j = 0; j < step.locations.length; j++) {
        if (!this.running) {
          return false;
        }
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
          await this.apiClient.updateLocation(
            vehicle.ID,
            location.lat,
            location.lng,
            bearing,
            speedKmh
          );
          await this.wait(waitMs);
        } else {
          await this.apiClient.updateLocation(
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
    return true;
  }

  private async getDirections(
    rideRequestId: number,
    startPoint: LatLng | null
  ): Promise<RouteStep[]> {
    let directions = await this.apiClient.getDirectionsV1(
      rideRequestId,
      startPoint
    );
    while (!directions) {
      const sleepSeconds = randomIntFromInterval(30 * 1000, 90 * 1000);
      await this.log(`failed to get directions, sleeping ${sleepSeconds}s`);
      await this.wait(sleepSeconds);
      directions = await this.apiClient.getDirectionsV1(
        rideRequestId,
        startPoint
      );
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
