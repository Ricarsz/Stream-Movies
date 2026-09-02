import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function Navbar() {
  const { user, logout } = useAuth();

  return (
    <nav className="navbar">
      <Link to="/" className="nav-brand">goMovies</Link>
      <div className="nav-links">
        {user ? (
          <>
            <span className="nav-user">Hi, {user.first_name}</span>
            {user.role === 'ADMIN' && (
              <Link to="/add-movie" className="nav-link">Add Movie</Link>
            )}
            <button onClick={logout} className="nav-btn">Logout</button>
          </>
        ) : (
          <>
            <Link to="/login" className="nav-link">Login</Link>
            <Link to="/register" className="nav-link">Register</Link>
          </>
        )}
      </div>
    </nav>
  );
}
