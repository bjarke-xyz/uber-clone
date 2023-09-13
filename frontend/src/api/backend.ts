import { User } from "firebase/auth";

export const baseUrl = import.meta.env.VITE_API_BASE_URL;
const simBaseUrl = import.meta.env.VITE_SIM_BASE_URL;

export class BackendApi {
  async getVehicles(): Promise<Vehicle[]> {
    return fetch(`${baseUrl}/v1/sim-vehicles`).then((res) => res.json());
  }

  async getUsers(): Promise<BackendUser[]> {
    return fetch(`${baseUrl}/v1/sim-users`).then((res) => res.json());
  }

  async getRideRequests(): Promise<RideRequest[]> {
    return fetch(`${baseUrl}/v1/sim-rides`).then((res) => res.json());
  }

  async getRecentLogs(): Promise<LogEvent["data"][]> {
    return fetch(`${baseUrl}/v1/sim/logs`).then((res) => res.json());
  }

  async getCurrencies(): Promise<Currency[]> {
    return fetch(`${baseUrl}/v1/payments/currencies`).then((res) => res.json());
  }

  async getSimStatus(): Promise<SimStatus[]> {
    return fetch(`${simBaseUrl}/api/admin/status`).then((res) => res.json());
  }
  async startSim(): Promise<Response> {
    return fetch(`${simBaseUrl}/api/admin/start`, { method: "POST" });
  }
  async stopSim(user: User): Promise<Response> {
    return fetch(`${simBaseUrl}/api/admin/stop`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${await user.getIdToken()}`,
      },
    });
  }
}
export const backendApi = new BackendApi();

export interface SimStatus {
  email: string;
  user: BackendUser | null;
  isRider: boolean;
  state: "STARTING" | "STARTED" | "STOPPING" | "STOPPED";
}

export interface PositionEvent {
  data: {
    vehicleId: number;
    lat: number;
    lng: number;
    recordedAt: string;
    bearing: number;
    speed: number;
  };
}

export interface LogEvent {
  data: {
    userId: number;
    tag: string;
    message: string;
    timestamp: string;
  };
}

export interface Vehicle {
  ID: number;
  RegistrationCountry: number;
  RegistrationNumber: number;
  OwnerID: number;
  Icon: string;
  lastRecordedPosition: PositionEvent["data"] | null;
}

export interface BackendUser {
  id: number;
  name: string;
  userId: string;
  simulated: boolean;
}

export interface RideRequest {
  id: number;
  riderId: number;
  driverId: number | null;
  fromLat: number;
  fromLng: number;
  fromName: string;
  toLat: number;
  toLng: number;
  toName: string;
  state: number;
  directionsVersion: number | null;
  directions: DirectionsV1 | null;
  price: number;
  currency: string;
  createdAt: string;
  updatedAt: string;
}

export interface DirectionsV1 {
  bbox: number[];
  routes: Route[];
  metadata: Metadata;
}

export interface Route {
  summary: Summary;
  segments: Segment[];
  bbox: number[];
  geometry: string;
  way_points: number[];
}

export interface Summary {
  distance: number;
  duration: number;
}

export interface Segment {
  distance: number;
  duration: number;
  steps: Step[];
}

export interface Step {
  distance: number;
  duration: number;
  type: number;
  instruction: string;
  name: string;
  way_points: number[];
  maneuver: Maneuver;
}

export interface Maneuver {
  location: number[];
  bearing_before: number;
  bearing_after: number;
}

export interface Metadata {
  attribution: string;
  service: string;
  timestamp: number;
  query: Query;
  engine: Engine;
}

export interface Query {
  coordinates: number[][];
  profile: string;
  format: string;
}

export interface Engine {
  version: string;
  build_date: string;
  graph_date: string;
}

/**
 * Decode an x,y or x,y,z encoded polyline
 * @param {*} encodedPolyline
 * @param {Boolean} includeElevation - true for x,y,z polyline
 * @returns {Array} of coordinates
 */
export const decodePolyline = (
  encodedPolyline: string,
  includeElevation = false
) => {
  // array that holds the points
  const points = [];
  let index = 0;
  const len = encodedPolyline.length;
  let lat = 0;
  let lng = 0;
  let ele = 0;
  while (index < len) {
    let b;
    let shift = 0;
    let result = 0;
    do {
      b = encodedPolyline.charAt(index++).charCodeAt(0) - 63; // finds ascii
      // and subtract it by 63
      result |= (b & 0x1f) << shift;
      shift += 5;
    } while (b >= 0x20);

    lat += (result & 1) !== 0 ? ~(result >> 1) : result >> 1;
    shift = 0;
    result = 0;
    do {
      b = encodedPolyline.charAt(index++).charCodeAt(0) - 63;
      result |= (b & 0x1f) << shift;
      shift += 5;
    } while (b >= 0x20);
    lng += (result & 1) !== 0 ? ~(result >> 1) : result >> 1;

    if (includeElevation) {
      shift = 0;
      result = 0;
      do {
        b = encodedPolyline.charAt(index++).charCodeAt(0) - 63;
        result |= (b & 0x1f) << shift;
        shift += 5;
      } while (b >= 0x20);
      ele += (result & 1) !== 0 ? ~(result >> 1) : result >> 1;
    }
    try {
      const location = [lat / 1e5, lng / 1e5];
      if (includeElevation) location.push(ele / 100);
      points.push(location);
    } catch (e) {
      console.log(e);
    }
  }
  return points;
};

export interface Currency {
  symbol: string;
  icon: string;
}
