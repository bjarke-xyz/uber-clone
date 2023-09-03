import { FirebaseOptions, initializeApp } from "firebase/app";
import { User, getAuth, signInWithEmailAndPassword } from "firebase/auth";
import {
  BackendUser,
  DirectionsV1,
  LatLng,
  PostLogInput,
  RideRequest,
  Vehicle,
} from "./types";

export class BackendApiClient {
  private user: User | null = null;
  constructor(private baseUrl: string) {}

  async postLog(input: PostLogInput): Promise<void> {
    try {
      const idToken = await this.mustGetToken();
      const resp = await fetch(`${this.baseUrl}/v1/me/log`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${idToken}`,
        },
        body: JSON.stringify(input),
      });
      if (resp.status > 299) {
        throw new Error(`got status ${resp.status}`);
      }
    } catch (error) {
      console.error("postLog failed", error);
    }
  }

  async getVehicle(): Promise<Vehicle | null> {
    try {
      const idToken = await this.mustGetToken();
      const resp = await fetch(`${this.baseUrl}/v1/vehicles`, {
        headers: {
          Authorization: `Bearer ${idToken}`,
        },
      });
      if (resp.status !== 200) {
        throw new Error(`got status ${resp.status}`);
      }
      const vehicles = (await resp.json()) as Vehicle[];
      return vehicles[0];
    } catch (error) {
      console.error("getVehicle failed", error);
      return null;
    }
  }

  async updateLocation(
    vehicleId: number,
    lat: number,
    lng: number,
    bearing: number,
    speed: number
  ) {
    try {
      const idToken = await this.mustGetToken();
      await fetch(`${this.baseUrl}/v1/vehicles/${vehicleId}/position`, {
        method: "PUT",
        body: JSON.stringify({
          lat,
          lng,
          bearing,
          speed,
        }),
        headers: {
          Authorization: `Bearer ${idToken}`,
        },
      });
    } catch (error) {
      console.error("updateLocation failed", error);
    }
  }

  async createRideRequest(
    from: LatLng,
    fromName: string,
    to: LatLng,
    toName: string
  ): Promise<RideRequest | null> {
    try {
      const idtoken = await this.mustGetToken();
      const resp = await fetch(`${this.baseUrl}/v1/rides/`, {
        method: "POST",
        body: JSON.stringify({
          fromLat: from.lat,
          fromLng: from.lng,
          fromName: fromName,
          toLat: to.lat,
          toLng: to.lng,
          toName: toName,
        }),
        headers: {
          Authorization: `Bearer ${idtoken}`,
        },
      });
      if (resp.status !== 200) {
        throw new Error(`requestRide returned ${resp.status}`);
      }
      const json = (await resp.json()) as RideRequest;
      return json;
    } catch (error) {
      console.error("requestRide failed", error);
      return null;
    }
  }

  async getRideRequest(id: number): Promise<RideRequest | null> {
    const myRides = await this.getMyRides();
    const rideRequest = myRides.find((x) => x.id === id) ?? null;
    return rideRequest;
  }

  async getMyRides(): Promise<RideRequest[]> {
    try {
      const idtoken = await this.mustGetToken();
      const resp = await fetch(`${this.baseUrl}/v1/rides/mine`, {
        headers: {
          Authorization: `Bearer ${idtoken}`,
        },
      });
      if (resp.status !== 200) {
        throw new Error(`getMyRides returned ${resp.status}`);
      }
      const json = (await resp.json()) as RideRequest[];
      return json;
    } catch (error) {
      console.error(`getMyRides failed`, error);
      return [];
    }
  }

  async getMyUser(): Promise<BackendUser | null> {
    try {
      const idtoken = await this.mustGetToken();
      const resp = await fetch(`${this.baseUrl}/v1/me/user`, {
        headers: {
          Authorization: `Bearer ${idtoken}`,
        },
      });
      if (resp.status !== 200) {
        throw new Error(`getMyUser returned ${resp.status}`);
      }
      const json = (await resp.json()) as BackendUser;
      return json;
    } catch (error) {
      console.error("getMyUsers failed", error);
      return null;
    }
  }

  async getAvailableRideRequests(): Promise<RideRequest[]> {
    try {
      const idtoken = await this.mustGetToken();
      const resp = await fetch(`${this.baseUrl}/v1/rides/available`, {
        headers: {
          Authorization: `Bearer ${idtoken}`,
        },
      });
      if (resp.status !== 200) {
        throw new Error(`getAvailableRideRequests returned ${resp.status}`);
      }
      const json = (await resp.json()) as RideRequest[];
      return json;
    } catch (error) {
      console.error(`getAvaialbleRideRequests failed`, error);
      return [];
    }
  }

  async claimRideRequest(id: number): Promise<boolean> {
    try {
      const idtoken = await this.mustGetToken();
      const resp = await fetch(`${this.baseUrl}/v1/rides/${id}/claim`, {
        method: "PUT",
        headers: {
          Authorization: `Bearer ${idtoken}`,
        },
      });
      if (resp.status !== 200) {
        throw new Error(`claimRideRequest ${id} returned ${resp.status}`);
      }
      return true;
    } catch (error) {
      console.log(`claimRideRequest ${id} failed`, error);
      return false;
    }
  }

  async finishRideRequest(id: number): Promise<boolean> {
    try {
      const idtoken = await this.mustGetToken();
      const resp = await fetch(`${this.baseUrl}/v1/rides/${id}/finish`, {
        method: "PUT",
        headers: {
          Authorization: `Bearer ${idtoken}`,
        },
      });
      if (resp.status !== 200) {
        throw new Error(`finishRideRequest ${id} returned ${resp.status}`);
      }
      return true;
    } catch (error) {
      console.log(`finishRideRequest ${id} failed`, error);
      return false;
    }
  }

  async getDirectionsV1(
    rideRequestId: number,
    startPoint: LatLng | null = null
  ): Promise<DirectionsV1 | null> {
    try {
      const idToken = await this.mustGetToken();
      const resp = await fetch(
        `${this.baseUrl}/v1/rides/${rideRequestId}/directions?startLat=${
          startPoint?.lat ?? ""
        }&startLng=${startPoint?.lng ?? ""}`,
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${idToken}`,
          },
        }
      );
      if (resp.status !== 200) {
        throw new Error(`getDirectionsV1 returned status ${resp.status}`);
      }
      const json = (await resp.json()) as DirectionsV1;
      return json;
    } catch (error) {
      console.log(`getDirectionsV1 failed ${rideRequestId}`, error);
      return null;
    }
  }

  private async mustGetToken(): Promise<string> {
    const idToken = await this.user?.getIdToken();
    if (!idToken) {
      throw new Error("Could not get token");
    }
    return idToken;
  }

  public async signIn(email: string, password: string) {
    if (this.user) {
      return;
    }
    const firebaseConfig: FirebaseOptions = {
      apiKey: process.env.FIREBASE_API_KEY,
    };

    const app = initializeApp(firebaseConfig);
    const auth = getAuth(app);
    const resp = await signInWithEmailAndPassword(auth, email, password);
    this.user = resp.user;
    const idToken = await this.user.getIdToken();
  }
}
