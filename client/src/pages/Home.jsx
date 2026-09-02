import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { getMovies } from '../api';

export default function Home() {
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    getMovies()
      .then((res) => setMovies(res.data || []))
      .catch(() => setError('Failed to load movies'))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="loading">Loading movies...</p>;
  if (error) return <p className="error">{error}</p>;

  return (
    <div className="container">
      <h1>Movies</h1>
      {movies.length === 0 ? (
        <p>No movies found.</p>
      ) : (
        <div className="movie-grid">
          {movies.map((movie) => (
            <Link to={`/movies/${movie.imdb_id}`} key={movie.imdb_id} className="movie-card">
              <img src={movie.poster_path} alt={movie.title} className="movie-poster" />
              <div className="movie-info">
                <h3>{movie.title}</h3>
                {movie.ranking && (
                  <span className="ranking">{movie.ranking.ranking_name}</span>
                )}
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
