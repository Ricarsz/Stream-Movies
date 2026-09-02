import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { addMovie } from '../api';

export default function AddMovie() {
  const navigate = useNavigate();
  const [form, setForm] = useState({
    imdb_id: '',
    title: '',
    poster_path: '',
    youtube_id: '',
    admin_review: '',
    ranking: { ranking_value: 5, ranking_name: 'Good' },
    genre: [],
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
      genre: [...form.genre, { genre_id: id, genre_name: name }],
    });
    setGenreInput('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await addMovie(form);
      navigate('/');
    } catch (err) {
      setError(err.response?.data?.error || 'Add movie failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container auth-page">
      <h1>Add Movie</h1>
      {error && <p className="error">{error}</p>}
      <form onSubmit={handleSubmit} className="auth-form">
        <input
          type="text"
          placeholder="IMDB ID (e.g. tt1234567)"
          value={form.imdb_id}
          onChange={(e) => setForm({ ...form, imdb_id: e.target.value })}
          required
        />
        <input
          type="text"
          placeholder="Title"
          value={form.title}
          onChange={(e) => setForm({ ...form, title: e.target.value })}
          required
        />
        <input
          type="url"
          placeholder="Poster URL"
          value={form.poster_path}
          onChange={(e) => setForm({ ...form, poster_path: e.target.value })}
          required
        />
        <input
          type="text"
          placeholder="YouTube Video ID"
          value={form.youtube_id}
          onChange={(e) => setForm({ ...form, youtube_id: e.target.value })}
          required
        />
        <textarea
          placeholder="Admin Review"
          value={form.admin_review}
          onChange={(e) => setForm({ ...form, admin_review: e.target.value })}
        />
        <select
          value={form.ranking.ranking_name}
          onChange={(e) => {
            const name = e.target.value;
            const vals = { Excellent: 10, Good: 7, Okay: 5, Bad: 3, Terrible: 1 };
            setForm({ ...form, ranking: { ranking_value: vals[name], ranking_name: name } });
          }}
        >
          <option value="Excellent">Excellent (10)</option>
          <option value="Good">Good (7)</option>
          <option value="Okay">Okay (5)</option>
          <option value="Bad">Bad (3)</option>
          <option value="Terrible">Terrible (1)</option>
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
        {form.genre.length > 0 && (
          <div className="genre-tags">
            {form.genre.map((g, i) => (
              <span key={i} className="genre-tag">
                {g.genre_name}
                <button type="button" onClick={() => setForm({ ...form, genre: form.genre.filter((_, j) => j !== i) })}>×</button>
              </span>
            ))}
          </div>
        )}
        <button type="submit" disabled={loading}>
          {loading ? 'Adding...' : 'Add Movie'}
        </button>
      </form>
    </div>
  );
}
