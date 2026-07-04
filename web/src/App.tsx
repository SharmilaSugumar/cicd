import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Layout } from './components/layout/Layout';
import { Login } from './pages/Login';
import { Dashboard } from './pages/Dashboard';
import { Organizations } from './pages/Organizations';
import { Projects } from './pages/Projects';
import { Pipelines } from './pages/Pipelines';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        
        <Route element={<Layout />}>
          <Route path="/" element={<Dashboard />} />
          <Route path="/organizations" element={<Organizations />} />
          <Route path="/projects" element={<Projects />} />
          <Route path="/pipelines" element={<Pipelines />} />

        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
