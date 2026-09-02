import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function Register() {
  const { register, login } = useAuth();
  const navigate = useNavigate();
  const [form, setForm] = useState({
    first_name: '',
    last_name: '',
    email: '',
    password: '',
    role: 'USER',
    favourite_genres: [],
  });
  const [genreInput, setGenreInput] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const addGenre = () => {
    if (!genreInput.trim()) return;
    const parts = genreInput.split(':');
    const id = parseInt(parts[0]);
    const name = parts[1] || parts[0];
    if (isNaN(id)) return;
    setForm({
      ...form,
      favourite_genres: [...form.favourite_genres, { genre_id: id, genre_name: name }],
    });
    setGenreInput('');
  };

  const removeGenre = (idx) => {
    setForm({
      ...form,
      favourite_genres: form.favourite_genres.filter((_, i) => i !== idx),
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await register(form);
      await login(form.email, form.password);
      navigate('/');
    } catch (err) {
      setError(err.response?.data?.error || 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container auth-page">
      <h1>Register</h1>
      {error && <p className="error">{error}</p>}
      <form onSubmit={handleSubmit} className="auth-form">
        <input
          type="text"
          placeholder="First Name"
          value={form.first_name}
          onChange={(e) => setForm({ ...form, first_name: e.target.value })}
          required
        />
        <input
          type="text"
          placeholder="Last Name"
          value={form.last_name}
          onChange={(e) => setForm({ ...form, last_name: e.target.value })}
          required
        />
        <input
          type="email"
          placeholder="Email"
          value={form.email}
          onChange={(e) => setForm({ ...form, email: e.target.value })}
          required
        />
        <input
          type="password"
          placeholder="Password (min 6 chars)"
          value={form.password}
          onChange={(e) => setForm({ ...form, password: e.target.value })}
          required
          minLength={6}
        />
        <select value={form.role} onChange={(e) => setForm({ ...form, role: e.target.value })}>
          <option value="USER">User</option>
          <option value="ADMIN">Admin</option>
        </select>
        <div className="genre-input-row">
          <input
            type="text"
            placeholder="Genre (e.g. 1:Action)"
            value={genreInput}
            onChange={(e) => setGenreInput(e.target.value)}
          />
          <button type="button" onClick={addGenre}>Add</button>
        </div>
        {form.favourite_genres.length > 0 && (
          <div className="genre-tags">
            {form.favourite_genres.map((g, i) => (
              <span key={i} className="genre-tag">
                {g.genre_name}
                <button type="button" onClick={() => removeGenre(i)}>×</button>
              </span>
            ))}
          </div>
        )}
        <button type="submit" disabled={loading}>
          {loading ? 'Registering...' : 'Register'}
        </button>
      </form>
      <p className="auth-switch">
        Already have an account? <Link to="/login">Login</Link>
      </p>
    </div>
  );
}
