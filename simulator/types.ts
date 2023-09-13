import { Position, Units, distance } from "@turf/turf";
import { BackendApiClient } from "./api-client";
import { setTimeout } from "timers/promises";
import { User } from "firebase/auth";

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
  legs: any[];
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

export interface OSMSearchResult {
  place_id: number;
  licence: string;
  osm_type: string;
  osm_id: number;
  boundingbox: string[];
  lat: string;
  lon: string;
  display_name: string;
  class: string;
  type: string;
  importance: number;
}

export class LatLng {
  constructor(public lat: number, public lng: number) {}

  equals(other: LatLng): boolean {
    return this.lat === other.lat && this.lng == other.lng;
  }

  distanceTo(other?: LatLng, units: Units = "meters"): number {
    if (!other) {
      return 0;
    }
    return distance(this.toPosition(), other.toPosition(), { units });
  }

  toPosition(): Position {
    return [this.lat, this.lng];
  }

  toString(): string {
    return `[${this.lat}, ${this.lng}]`;
  }
}

export function getAngleInDegrees(p1: LatLng, p2: LatLng): number {
  const pp1 = {
    x: p1.lat,
    y: p1.lng,
  };

  const pp2 = {
    x: p2.lat,
    y: p2.lng,
  };
  // angle in radians
  const angleRadians = Math.atan2(pp2.y - pp1.y, pp2.x - pp1.x);
  // angle in degrees
  const angleDeg = (Math.atan2(pp2.y - pp1.y, pp2.x - pp1.x) * 180) / Math.PI;

  // document.getElementById("rotation").innerHTML = "Rotation : " + angleDeg;
  return angleDeg;
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
export interface Vehicle {
  ID: number;
  RegistrationCountry: string;
  RegistrationNumber: string;
  OwnerID: number;
  lastRecordedPosition: PositionEvent["data"] | null;
}

export interface SimUser {
  email: string;
  password: string;
  isRider?: boolean;
  city: string;
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
  price: number;
  currency: string;
  createdAt: string;
  updatedAt: string;
}

export enum RideRequestState {
  Available,
  Accepted,
  InProgress,
  Finished,
}

export interface BackendUser {
  id: number;
  name: string;
  simulated: boolean;
  userId: string;
}

export interface CityData {
  type: string;
  generator: string;
  copyright: string;
  timestamp: string;
  features: CityDataFeature[];
}

export interface CityDataFeature {
  type: string;
  properties: CityDataProperties;
  geometry: CityDataGeometry;
  id: string;
}

export interface CityDataProperties {
  "@id": string;
  "addr:city": string;
  "addr:country": string;
  "addr:housenumber": string;
  "addr:municipality": string;
  "addr:postcode": string;
  "addr:street": string;
  "osak:identifier": string;
  source: string;
  "addr:place"?: string;
}

export interface CityDataGeometry {
  type: string;
  coordinates: number[];
}

export interface NamedPoint {
  name: string;
  location: LatLng;
}
export function cityDataFeatureToNamedPoints(
  feature: CityDataFeature
): NamedPoint {
  const name = `${feature.properties["addr:street"]} ${feature.properties["addr:housenumber"]}, ${feature.properties["addr:postcode"]} ${feature.properties["addr:city"]}`;
  const location = new LatLng(
    feature.geometry.coordinates[1],
    feature.geometry.coordinates[0]
  );
  return {
    name,
    location,
  };
}

export interface PostLogInput {
  tag: string;
  message: string;
}

export type RunnerState = "STARTING" | "STARTED" | "STOPPING" | "STOPPED";
export abstract class SimRunner {
  protected timeMultiplier: number = 16;
  protected running = false;
  protected state: RunnerState = "STOPPED";
  constructor(
    protected abortController: AbortController,
    protected apiClient: BackendApiClient,
    protected userEmail: string,
    protected userPassword: string,
    protected tag: string
  ) {
    this.signIn();
  }
  public abstract run(): Promise<void>;
  public setTimeMultiplier(timeMultiplier: number) {
    this.timeMultiplier = timeMultiplier;
  }
  public setAbortController(abortController: AbortController) {
    this.abortController = abortController;
  }
  public async signIn(): Promise<void> {
    this.apiClient.signIn(this.userEmail, this.userPassword);
  }
  public stop() {
    if (this.state === "STOPPED" || this.state === "STOPPING") {
      return;
    }
    this.log("stopping...");
    this.running = false;
    this.state = "STOPPING";
    this.abortController.abort("STOP");
  }
  public async stopped(): Promise<void> {
    this.state = "STOPPED";
    this.log("stopped");
  }
  public async starting(): Promise<void> {
    this.log("starting...");
    this.state = "STARTING";
    this.running = true;
  }
  public async started(): Promise<void> {
    this.log("started");
    this.state = "STARTED";
  }

  public getState() {
    return this.state;
  }
  public getUser(): BackendUser | null {
    return this.apiClient.getBackendUser();
  }
  async wait(timeMs: number) {
    await setTimeout(timeMs, null, { signal: this.abortController.signal });
  }
  protected async log(message?: any, ...optionalParams: any[]) {
    console.log(
      `${this.userEmail} [${this.tag}] | ${message}`,
      ...optionalParams
    );
    await this.apiClient.postLog({ tag: this.tag, message });
  }
}

export function isAbortError(error: unknown): boolean {
  return (error as any)?.code === "ABORT_ERR";
}
