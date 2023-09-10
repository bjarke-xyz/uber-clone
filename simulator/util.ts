import path from "path";
import { CityData } from "./types";
import { readFile } from "fs/promises";

export function randomIntFromInterval(min: number, max: number) {
  // min and max included
  return Math.floor(Math.random() * (max - min + 1) + min);
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
  let points = [];
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
      let location = [lat / 1e5, lng / 1e5];
      if (includeElevation) location.push(ele / 100);
      points.push(location);
    } catch (e) {
      console.log(e);
    }
  }
  return points;
};

export async function getCityData(city: string): Promise<CityData> {
  const resource = `${city}_clean.geojson`;
  const url = path.resolve(__dirname, `./data/random-city-data/${resource}`);
  const file = await readFile(url, { encoding: "utf-8" });
  const cityData = JSON.parse(file) as CityData;
  return cityData;
}

class Urls {
  private _backendApiBaseUrl: string = "missing API_BASE_URL";
  private _authUrl: string = "missing AUTH_URL";
  public get backendApiBaseUrl() {
    return this._backendApiBaseUrl;
  }
  public get authUrl() {
    return this._authUrl;
  }
  public load() {
    this._backendApiBaseUrl =
      process.env.API_BASE_URL ?? "Missing API_BASE_URL";
    this._authUrl = process.env.AUTH_URL ?? "Missing AUTH_URL";
  }
}
export const urls = new Urls();
