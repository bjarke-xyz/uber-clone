import { NavLink } from "react-router-dom";

export function NavigationPage() {
  return (
    <div>
      <ul>
        <li>
          <NavLink to="/overview">Overview</NavLink>
        </li>
        <li>
          <NavLink to="/driver">Driver</NavLink>
        </li>
        <li>
          <NavLink to="/rider">Rider</NavLink>
        </li>
      </ul>
    </div>
  );
}
