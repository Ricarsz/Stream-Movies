import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getMovie, deleteMovie } from '../api';
import { useAuth } from '../context/AuthContext';

export default function MovieDetail() {
  const { imdb_id } = useParams();
  const { user } = useAuth();
  const navigate = useNavigate();
  const [movie, setMovie] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getMovie(imdb_id)
      .then((res) => setMovie(res.data))
      .catch(() => navigate('/'))
      .finally(() => setLoading(false));
  }, [imdb_id, navigate]);

  const handleDelete = async () => {
    if (!window.confirm('Delete this movie?')) return;
    try {
      await deleteMovie(imdb_id);
      navigate('/');
    } catch {
      alert('Delete failed');
    }
  };

  if (loading) return <p className="loading">Loading...</p>;
  if (!movie) return null;

  return (
    <div className="container movie-detail">
      <img src={movie.poster_path} alt={movie.title} className="detail-poster" />
      <div className="detail-info">
        <h1>{movie.title}</h1>
        <p><strong>IMDB ID:</strong> {movie.imdb_id}</p>
        {movie.ranking && (
          <p><strong>Ranking:</strong> {movie.ranking.ranking_name} ({movie.ranking.ranking_value}/10)</p>
        )}
        {movie.genre && movie.genre.length > 0 && (
          <p><strong>Genre:</strong> {movie.genre.map((g) => g.genre_name).join(', ')}</p>
        )}
        {movie.admin_review && <p><strong>Review:</strong> {movie.admin_review}</p>}
        {movie.youtube_id && (
          <a
            href={`https://www.youtube.com/watch?v=${movie.youtube_id}`}
            target="_blank"
            rel="noreferrer"
            className="yt-link"
          >
            Watch Trailer on YouTube
          </a>
        )}
        {user && user.role === 'ADMIN' && (
          <div className="admin-actions">
            <button onClick={() => navigate(`/edit-movie/${imdb_id}`)} className="btn-edit">Edit</button>
            <button onClick={handleDelete} className="btn-delete">Delete</button>
          </div>
        )}
      </div>
    </div>
  );
}
