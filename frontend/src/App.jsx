import { useState } from "react";
import reactLogo from "./assets/react.svg";
import viteLogo from "/vite.svg";
import "./App.css";
import {
  BrowserRouter,
  Routes,
  Route,
  Outlet,
  NavLink,
} from "react-router-dom";
import {
  QueryClient,
  QueryClientProvider,
  useQuery,
} from "@tanstack/react-query";

const queryClient = new QueryClient();
// layout
function LayOut() {
  return (
    <div>
      <nav>
        <NavLink to="/logs" end>
          LogHome
        </NavLink>
        <NavLink to="/logs/all" end>
          LogAll
        </NavLink>
        <NavLink to="/logs/since" end>
          LogSince
        </NavLink>
      </nav>
      <Outlet />
    </div>
  );
}
// logs
function LogHome() {
  return (
    <div>
      <h2>LogHome</h2>
    </div>
  );
}
// logs/all
function LogAll() {
  return (
    <main>
      <header>
        <h2>LogAll</h2>
      </header>
      <section>
        <h3>这里展示日志</h3>
        <QueryClientProvider client={queryClient}>
          <QueryLogAll />
        </QueryClientProvider>
      </section>
      <footer>
        <p>正在搭建ing...</p>
      </footer>
    </main>
  );
}
function QueryLogAll() {
  const { isPending, isError, data, error } = useQuery({
    queryKey: ["fetchlogall"],
    queryFn: fetchlogall,
  });
  if (isPending) {
    return <p>Loading...</p>;
  }
  if (isError) {
    return <p>Error: {error.message}</p>;
  }
  return;
}
// logs/since
function LogSince() {
  return (
    <div>
      <h2>LogSince</h2>
    </div>
  );
}

function App() {
  const [count, setCount] = useState(0);

  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<LayOut />}>
            <Route path="logs" element={<LogHome />} />
            <Route path="/logs/all" element={<LogAll />} />
            <Route path="/logs/since" element={<LogSince />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </>
  );
}

export default App;
