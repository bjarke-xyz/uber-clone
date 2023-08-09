import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Outlet, Route, Routes } from "react-router-dom";
import "./App.css";
import { DriverPage } from "./pages/DriverPage";
import { NavigationPage } from "./pages/NavigationPage";
import { OverviewPage } from "./pages/OverviewPage";
import { RiderPage } from "./pages/RiderPage";

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Routes>
        <Route element={<Layout />}>
          <Route index element={<NavigationPage />} />
          <Route path="home" element={<NavigationPage />} />
          <Route path="driver" element={<DriverPage />} />
          <Route path="rider" element={<RiderPage />} />
          <Route path="overview" element={<OverviewPage />} />
          <Route path="*" element={<p>There's nothing here: 404!</p>} />
        </Route>
      </Routes>
    </QueryClientProvider>
  );
}

const Layout = () => {
  return (
    <>
      <header className="p-2 flex">
        <a href="/">
          <h1 className="text-3xl text-blue-600 font-bold align-top">
            Uber clone
          </h1>
        </a>
      </header>

      <main className="">
        <Outlet />
      </main>
    </>
  );
};

export default App;
