import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useAtom } from "jotai";
import { useEffect } from "react";
import { Outlet, Route, Routes } from "react-router-dom";
import "./App.css";
import { firebaseAuth } from "./api/firebase";
import { DriverPage } from "./pages/DriverPage";
import { LoginPage } from "./pages/LoginPage";
import { NavigationPage } from "./pages/NavigationPage";
import { OverviewPage } from "./pages/OverviewPage";
import { RiderPage } from "./pages/RiderPage";
import { userAtom } from "./store/user";

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
          <Route path="login" element={<LoginPage />} />
          <Route path="*" element={<p>There's nothing here: 404!</p>} />
        </Route>
      </Routes>
    </QueryClientProvider>
  );
}

const Layout = () => {
  const [user, setUser] = useAtom(userAtom);
  useEffect(() => {
    const unsubscribe = firebaseAuth.onAuthStateChanged((user) => {
      setUser(user);
    });
    return unsubscribe;
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
  return (
    <>
      <header className="p-2 flex justify-between">
        <a href="/">
          <h1 className="text-3xl text-blue-600 font-bold align-top">
            Uber clone
          </h1>
        </a>
        <span onClick={() => firebaseAuth.signOut()}>{user && user.email}</span>
      </header>

      <main className="">
        <Outlet />
      </main>
    </>
  );
};

export default App;
