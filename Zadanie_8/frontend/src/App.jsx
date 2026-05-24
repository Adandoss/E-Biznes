import { BrowserRouter as Router, Routes, Route, NavLink, Navigate } from 'react-router-dom';
import { CartProvider } from './context/CartContext';
import { AuthProvider, useAuth } from './context/AuthContext';
import Products from './components/Products';
import Cart from './components/Cart';
import Payments from './components/Payments';
import Login from './components/Login';
import Register from './components/Register';
import AuthCallback from './components/AuthCallback';

function AppContent() {
  const { user, logout, loading } = useAuth();

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-100">
        <p className="text-xl font-semibold text-gray-600">Ładowanie...</p>
      </div>
    );
  }

  const mainContent = (
    <div className="min-h-screen bg-gray-50 flex flex-col font-sans">
      <header className="bg-blue-600 text-white shadow-md">
        <div className="max-w-6xl mx-auto px-4 py-3 flex flex-wrap justify-between items-center">
          <h1 className="text-2xl font-bold tracking-tight">Mój Sklep</h1>
          <nav className="flex items-center gap-3">
            {user ? (
              <>
                <NavLink 
                  to="/" 
                  end 
                  className={({ isActive }) => 
                    `px-3 py-2 rounded-md text-sm font-medium transition-colors ${isActive ? 'bg-blue-800 text-white' : 'text-blue-100 hover:bg-blue-700 hover:text-white'}`
                  }
                >
                  Produkty
                </NavLink>
                <NavLink 
                  to="/cart" 
                  className={({ isActive }) => 
                    `px-3 py-2 rounded-md text-sm font-medium transition-colors ${isActive ? 'bg-blue-800 text-white' : 'text-blue-100 hover:bg-blue-700 hover:text-white'}`
                  }
                >
                  Koszyk
                </NavLink>
                <NavLink 
                  to="/payments" 
                  className={({ isActive }) => 
                    `px-3 py-2 rounded-md text-sm font-medium transition-colors ${isActive ? 'bg-blue-800 text-white' : 'text-blue-100 hover:bg-blue-700 hover:text-white'}`
                  }
                >
                  Płatności
                </NavLink>
                <div className="flex items-center gap-2 border-l border-blue-500 pl-3 ml-1">
                  <span className="text-sm font-semibold">{user.name || user.email}</span>
                  <span className="bg-blue-500 text-blue-100 text-xs px-2 py-0.5 rounded-full uppercase">{user.provider}</span>
                  <button onClick={logout} className="ml-2 bg-red-500 hover:bg-red-600 text-white text-xs px-3 py-1.5 rounded font-medium transition-colors">
                    Wyloguj
                  </button>
                </div>
              </>
            ) : (
              <NavLink 
                to="/login" 
                className={({ isActive }) => 
                  `px-3 py-2 rounded-md text-sm font-medium transition-colors ${isActive ? 'bg-blue-800 text-white' : 'text-blue-100 hover:bg-blue-700 hover:text-white'}`
                }
              >
                Zaloguj się
              </NavLink>
            )}
          </nav>
        </div>
      </header>

      <main className="max-w-6xl w-full mx-auto px-4 py-8 flex-grow">
        <Routes>
          <Route path="/login" element={user ? <Navigate to="/" /> : <Login />} />
          <Route path="/register" element={user ? <Navigate to="/" /> : <Register />} />
          <Route path="/auth/callback" element={<AuthCallback />} />
          <Route path="/" element={user ? <Products /> : <Navigate to="/login" />} />
          <Route path="/cart" element={user ? <Cart /> : <Navigate to="/login" />} />
          <Route path="/payments" element={user ? <Payments /> : <Navigate to="/login" />} />
        </Routes>
      </main>
    </div>
  );

  return user ? <CartProvider>{mainContent}</CartProvider> : mainContent;
}

function App() {
  return (
    <Router>
      <AuthProvider>
        <AppContent />
      </AuthProvider>
    </Router>
  );
}

export default App;
