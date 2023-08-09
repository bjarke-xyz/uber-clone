import { LatLng, OSMSearchResult } from "./types";

export class OSMAPI {
  private cache = new Map<string, string>();
  async search(query: string): Promise<ProcessedSearchResult | null> {
    try {
      const cacheKey = `SEARCH:${query}`;
      let json: OSMSearchResult[] = [];
      if (this.cache.has(cacheKey)) {
        const jsonStr = this.cache.get(cacheKey) ?? "";
        json = JSON.parse(jsonStr) as OSMSearchResult[];
      } else {
        const resp = await fetch(
          `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(
            query
          )}&format=json&limit=1`
        );
        const jsonStr = await resp.text();
        this.cache.set(cacheKey, jsonStr);
        json = JSON.parse(jsonStr) as OSMSearchResult[];
      }
      if (json.length === 0) {
        return null;
      }
      const result = json[0];
      return {
        location: new LatLng(parseFloat(result.lat), parseFloat(result.lon)),
        displayName: result.display_name,
      };
    } catch (error) {
      console.error("OSM search failed", error);
      return null;
    }
  }

  // async getDirections(points: LatLng[]): Promise<Directions | null> {
  //   try {
  //     const formattedPointsList = points.map((x) => `${x.lng},${x.lat}`);
  //     const formattedPointStr = formattedPointsList.join(";");
  //     const cacheKey = `DIRECTION:${formattedPointStr}`;
  //     let directions: Directions | null;
  //     if (this.cache.has(cacheKey)) {
  //       const jsonStr = this.cache.get(cacheKey) ?? "";
  //       directions = JSON.parse(jsonStr) as Directions;
  //     } else {
  //       const resp = await fetch(
  //         `https://routing.openstreetmap.de/routed-car/route/v1/driving/${formattedPointStr}?overview=false&geometries=geojson&steps=true`
  //       );
  //       if (resp.status !== 200) {
  //         console.log("getDirections got status", resp.status);
  //       }
  //       const jsonStr = await resp.text();
  //       this.cache.set(cacheKey, jsonStr);
  //       directions = JSON.parse(jsonStr) as Directions;
  //     }
  //     return directions;
  //   } catch (error) {
  //     console.error("getDirections failed", error);
  //     return null;
  //   }
  // }
}

export const osmApi = new OSMAPI();

export interface ProcessedSearchResult {
  location: LatLng;
  displayName: string;
}
