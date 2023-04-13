import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Home from './pages/Home';
import CreatePackage from './pages/CreatePackage';
import UpdatePackage from './pages/UpdatePackage';
import PackageDetails from './pages/PackageDetails';
import SearchPackages from './pages/SearchPackages';
import PackagesList from './pages/PackagesList';
import RatePackage from './pages/RatePackage';

function App() {
  return (
    <Router>
      <div className="App">
        {/* Add a Header and/or Navigation component here */}
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/create" element={<CreatePackage />} />
          <Route path="/update/:id" element={<UpdatePackage />} />
          <Route path="/package/:id" element={<PackageDetails />} />
          <Route path="/search" element={<SearchPackages />} />
          <Route path="/package/:id/rate" element={<RatePackage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
