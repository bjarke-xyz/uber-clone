import { fetchEventSource } from "@microsoft/fetch-event-source";
import PersonIcon from "@mui/icons-material/Person";
import { useQuery } from "@tanstack/react-query";
import L from "leaflet";
import "leaflet-rotatedmarker";
import { PropsWithChildren, useEffect, useMemo, useRef, useState } from "react";
import {
  MapContainer,
  Marker,
  Polyline,
  TileLayer,
  Tooltip,
} from "react-leaflet";
import {
  BackendUser,
  Currency,
  LogEvent,
  PositionEvent,
  RideRequest,
  Vehicle,
  backendApi,
  baseUrl,
  decodePolyline,
} from "../api/backend";
import "./OverviewPage.css";
import LocalTaxiIcon from "@mui/icons-material/LocalTaxi";
import { takeRight } from "lodash";
import { format, parseISO } from "date-fns";
import markerIconPng from "leaflet/dist/images/marker-icon.png";
import { Icon } from "leaflet";

const defaultIcon = new Icon({
  iconUrl: markerIconPng,
  iconSize: [25, 41],
  iconAnchor: [12, 41],
});

const position = new L.LatLng(55.82905, 10.20355);

type SideSectionProps = {
  title: string;
  className?: string;
};
function SideSection(props: PropsWithChildren<SideSectionProps>) {
  return (
    <div className={`bg-slate-100 mb-4 p-2 ${props.className ?? ""}`}>
      <h3 className="text-xl font-bold p-4">{props.title}</h3>
      {props.children}
    </div>
  );
}

const maxLogLines = 100;

export function OverviewPage() {
  const mapRef = useRef<L.Map | null>(null);

  const logsDiv = useRef<HTMLDivElement | null>(null);
  // useMemo(() => {
  //   setInterval(() => {
  //     if (logsDiv?.current) {
  //       logsDiv.current.scroll({
  //         top: logsDiv.current.scrollHeight,
  //         behavior: "smooth",
  //       });
  //     }
  //   }, 1000);
  // }, []);

  const [eventMap, setEventMap] = useState<Record<number, PositionEvent>>({});

  const [logs, setLogs] = useState<LogEvent[]>([]);
  useEffect(() => {
    async function getData() {
      const logs = await backendApi.getRecentLogs();
      const logEvents: LogEvent[] = logs.map((x) => ({ data: x }));
      setLogs(logEvents);
    }
    getData();
  }, []);

  const [currencies, setCurrencies] = useState<Record<string, Currency>>({});
  useEffect(() => {
    async function getData() {
      const currencies = await backendApi.getCurrencies();
      const currencyMap: Record<string, Currency> = {};
      for (const c of currencies) {
        currencyMap[c.symbol] = c;
      }
      setCurrencies(currencyMap);
    }
    getData();
  }, []);

  const vehiclesQuery = useQuery({
    queryKey: ["vehicles"],
    queryFn: backendApi.getVehicles,
  });
  const [vehiclesMap, setVehiclesMap] = useState<Record<number, Vehicle>>({});
  useMemo(() => {
    const _vehiclesMap: Record<number, Vehicle> = {};
    for (const vehicle of vehiclesQuery?.data ?? []) {
      _vehiclesMap[vehicle.ID] = vehicle;
      if (vehicle.lastRecordedPosition) {
        const fakePositionEvent: PositionEvent = {
          data: {
            ...vehicle.lastRecordedPosition,
          },
        };
        setEventMap((m) => {
          return {
            ...m,
            [vehicle.ID]: fakePositionEvent,
          };
        });
      }
    }
    setVehiclesMap(_vehiclesMap);
  }, [vehiclesQuery?.data]);

  const usersQuery = useQuery({
    queryKey: ["users"],
    queryFn: backendApi.getUsers,
  });
  const [usersMap, setUsersMap] = useState<Record<number, BackendUser>>({});
  useMemo(() => {
    const _usersMap: Record<number, BackendUser> = {};
    for (const user of usersQuery?.data ?? []) {
      _usersMap[user.id] = user;
    }
    setUsersMap(_usersMap);
  }, [usersQuery?.data]);

  const ridesQuery = useQuery({
    queryKey: ["rides"],
    queryFn: backendApi.getRideRequests,
    refetchInterval: 5000,
  });

  useEffect(() => {
    const doSse = async () => {
      await fetchEventSource(`${baseUrl}/v1/sim/events`, {
        onmessage(ev) {
          const evData = ev.data?.trim();
          if (evData) {
            switch (ev.event) {
              case "position-update": {
                const event = JSON.parse(ev.data) as PositionEvent;
                setEventMap((m) => {
                  return {
                    ...m,
                    [event.data.vehicleId]: event,
                  };
                });
                break;
              }
              case "user-log": {
                const event = JSON.parse(ev.data) as LogEvent;
                setLogs((logs) => {
                  const newLogs = takeRight(logs, maxLogLines);
                  return [event, ...newLogs];
                });
              }
            }
          }
        },
        onerror(error) {
          console.log("error", error);
        },
        openWhenHidden: true,
      });
    };
    doSse();
  }, []);

  const [activeRide, setActiveRide] = useState<RideRequest | null>(null);

  const onActiveRideClickedHandler = (
    e: React.MouseEvent<HTMLDivElement>,
    ride: RideRequest | null
  ) => {
    e.preventDefault();
    if (ride == null) {
      setActiveRide(null);
    } else {
      if (ride.id === activeRide?.id) {
        setActiveRide(null);
      } else {
        setActiveRide(ride);
      }
    }
  };

  return (
    <div className="w-full h-full flex">
      <aside className="w-96 bg-blue-50">
        <SideSection title="Active rides" className="min-h-[28rem]">
          {ridesQuery?.data?.map((ride) => (
            <div
              onClick={(e) => onActiveRideClickedHandler(e, ride)}
              key={ride.id}
              className={`mb-4 p-2 bg-sky-100 hover:bg-sky-200 rounded-md cursor-pointer ${
                activeRide?.id === ride.id ? "bg-sky-200" : ""
              }`}
            >
              <div className="flex items-center justify-between">
                <div className="flex">
                  <PersonIcon />
                  <div className="ml-2 text-lg">
                    {usersMap[ride.riderId]?.name ?? "Unknown user"}
                  </div>
                </div>
                <div className="flex self-end">
                  {ride.price / 100} {currencies[ride.currency]?.icon}
                </div>
              </div>
              {ride.driverId && (
                <div className="flex">
                  <LocalTaxiIcon />
                  <div className="ml-2 text-lg">
                    {usersMap[ride.driverId]?.name ?? "Unknown user"}
                  </div>
                </div>
              )}
            </div>
          ))}
        </SideSection>
        <SideSection title="Logs" className="bg-slate-800 text-gray-50 h-96">
          <div ref={logsDiv} className="h-72 overflow-y-auto">
            {logs.map((log, i) => (
              <div
                key={`${log.data.timestamp}-${i}`}
                className={`text-xs font-mono p-2 rounded-md  ${
                  activeRide
                    ? activeRide.driverId === log.data.userId ||
                      activeRide.riderId === log.data.userId
                      ? "text-lime-200"
                      : ""
                    : ""
                }`}
              >
                {format(parseISO(log.data.timestamp), "HH:mm:ss")}{" "}
                {usersMap[log.data.userId]?.name ?? "Unknown user"} [
                {log.data.tag}] | {log.data.message}
              </div>
            ))}
          </div>
        </SideSection>
      </aside>
      <main className="h-[calc(100vh-3.5rem)] flex-1 bg-lime-100">
        <MapContainer
          ref={mapRef}
          className="h-full w-full"
          center={position}
          zoom={8}
        >
          <TileLayer
            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          {activeRide && (
            <>
              <Marker
                icon={defaultIcon}
                position={new L.LatLng(activeRide.fromLat, activeRide.fromLng)}
              >
                <Tooltip direction="bottom" permanent>
                  A | {activeRide.fromName}
                </Tooltip>
              </Marker>
              <Marker
                icon={defaultIcon}
                position={new L.LatLng(activeRide.toLat, activeRide.toLng)}
              >
                <Tooltip direction="bottom" permanent>
                  B | {activeRide.toName}
                </Tooltip>
              </Marker>
              <Polyline
                positions={decodePolyline(
                  activeRide.directions?.routes[0]?.geometry ?? ""
                ).map(([lat, lng]) => new L.LatLng(lat, lng))}
              ></Polyline>
            </>
          )}
          {Object.keys(eventMap)
            .map((x) => parseInt(x))
            .map((key) => {
              const data = eventMap[key]?.data;
              const vehicle = vehiclesMap[key];
              const user = usersMap[vehicle?.OwnerID];
              const isActiveRide = activeRide?.driverId === user?.id;
              return !data || !vehicle || !user ? null : (
                <VehicleMarker
                  key={key}
                  eventData={data}
                  vehicle={vehicle}
                  user={user}
                  isActiveRide={isActiveRide}
                />
              );
            })}
        </MapContainer>
      </main>
    </div>
  );
}

interface VehicleMarkerProps {
  eventData: PositionEvent["data"];
  vehicle: Vehicle;
  user: BackendUser;
  isActiveRide: boolean;
}
function VehicleMarker({
  eventData,
  vehicle,
  user,
  isActiveRide,
}: VehicleMarkerProps) {
  const markerRef = useRef<L.Marker | null>(null);
  useEffect(() => {
    if (markerRef?.current) {
      markerRef.current.setRotationAngle(eventData.bearing);
    }
  }, [eventData, markerRef]);
  let iconName = vehicle.Icon;
  if (!iconName) {
    iconName = "vehicle001.png";
  }
  const iconUrl = `https://static.bjarke.xyz/uber-clone/vehicles/${iconName}`;
  const icon = new L.Icon({
    iconUrl,
    iconSize: [41, 41],
    iconAnchor: [20, 20],
    popupAnchor: [1, -34],
    className: `marker-icon`,
  });
  return (
    <Marker
      ref={markerRef}
      position={new L.LatLng(eventData.lat, eventData.lng)}
      icon={icon}
      rotationAngle={eventData.bearing}
      rotationOrigin="center"
    >
      {isActiveRide && (
        <Tooltip direction="bottom" offset={[0, 20]} permanent={true}>
          {user.name} | {Math.round(eventData.speed)} km/h
        </Tooltip>
      )}
    </Marker>
  );
}
