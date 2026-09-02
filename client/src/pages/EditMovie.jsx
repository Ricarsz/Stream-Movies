import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getMovie, updateMovie } from '../api';

export default function EditMovie() {
  const { imdb_id } = useParams();
  const navigate = useNavigate();
  const [form, setForm] = useState(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getMovie(imdb_id)
      .then((res) => {
        const m = res.data;
        setForm({
          imdb_id: m.imdb_id,
          title: m.title,
          poster_path: m.poster_path,
          youtube_id: m.youtube_id || '',
          admin_review: m.admin_review || '',
          ranking: m.ranking || { ranking_value: 5, ranking_name: 'Good' },
          genre: m.genre || [],
        });
      })
      .catch(() => navigate('/'))
      .finally(() => setLoading(false));
  }, [imdb_id, navigate]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      await updateMovie(imdb_id, form);
      navigate(`/movies/${imdb_id}`);
    } catch (err) {
      setError(err.response?.data?.error || 'Update failed');
    }
  };

  if (loading) return <p className="loading">Loading...</p>;
  if (!form) return null;

  return (
    <div className="container auth-page">
      <h1>Edit Movie</h1>
      {error && <p className="error">{error}</p>}
      <form onSubmit={handleSubmit} className="auth-form">
        <input type="text" value={form.imdb_id} disabled />
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
        <button type="submit">Update Movie</button>
      </form>
    </div>
  );
}
