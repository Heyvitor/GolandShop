import { Switch, Route } from 'wouter';
import { AuthProvider } from './context/AuthContext';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import StoreView from './pages/StoreView';

function App() {
  return (
    <AuthProvider>
      <Switch>
        <Route path="/login" component={Login} />
        <Route path="/register" component={Register} />
        <Route path="/" component={Dashboard} />
        <Route path="/:slug" component={StoreView} />
      </Switch>
    </AuthProvider>
  );
}

export default App;
