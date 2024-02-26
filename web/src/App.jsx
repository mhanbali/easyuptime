import {
  BrowserRouter as Router,
  Routes,
  Route,
  Link,
} from "react-router-dom";
import Dashboard from "./components/Dashboard";
import Domain from "./components/Domain";
const App = () => {

  return (
    <Router>
      <div className="App">
        <ul className="App-header">
          <li>
            <Link to="/">Dashboard</Link>
          </li>
        </ul>
        <Routes>
          <Route
            exact
            path="/"
            element={<Dashboard />}>
          </Route>
          <Route
            exact
            path="/domain/:domain"
            element={<Domain />}>
          </Route>
        </Routes>
      </div>
    </Router>
  );

}

export default App;