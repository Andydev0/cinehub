package servico

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Andydev0/filmes-backend/internal/dominio"
	"github.com/Andydev0/filmes-backend/internal/repositorio"
)

type RecomendacaoServico interface {
	RecomendarFilmes(usuarioID int64) ([]dominio.Filme, error)
}

type recomendacaoServicoImpl struct {
	favoritoRepo repositorio.FavoritoRepositorio
	filmeServico FilmeServico
	apiKey       string
}

func NovoRecomendacaoServico(favoritoRepo repositorio.FavoritoRepositorio, filmeServico FilmeServico, apiKey string) RecomendacaoServico {
	return &recomendacaoServicoImpl{
		favoritoRepo: favoritoRepo,
		filmeServico: filmeServico,
		apiKey:       apiKey,
	}
}

func (s *recomendacaoServicoImpl) RecomendarFilmes(usuarioID int64) ([]dominio.Filme, error) {
	favoritos, err := s.favoritoRepo.ListarPorUsuarioID(usuarioID)
	if err != nil || len(favoritos) == 0 {
		return make([]dominio.Filme, 0), err
	}

	contagemGeneros := make(map[int]int)
	mapaFavoritos := make(map[int64]bool)
	for _, fav := range favoritos {
		mapaFavoritos[fav.FilmeID] = true
		detalhes, err := s.filmeServico.BuscarDetalhes(fav.FilmeID)
		if err == nil {
			for _, genero := range detalhes.Generos {
				contagemGeneros[genero.ID]++
			}
		}
	}

	if len(contagemGeneros) == 0 {
		return make([]dominio.Filme, 0), nil
	}

	generoMaisComumID := -1
	maxContagem := 0
	for id, contagem := range contagemGeneros {
		if contagem > maxContagem {
			maxContagem = contagem
			generoMaisComumID = id
		}
	}

	url := fmt.Sprintf("https://api.themoviedb.org/3/discover/movie?api_key=%s&with_genres=%d&language=pt-BR&sort_by=popularity.desc", s.apiKey, generoMaisComumID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var discoverResponse dominio.RespostaBuscaTMDB
	if err := json.NewDecoder(resp.Body).Decode(&discoverResponse); err != nil {
		return nil, err
	}

	var recomendacoes []dominio.Filme
	for _, filmeTMDB := range discoverResponse.Resultados {
		if _, ehFavorito := mapaFavoritos[int64(filmeTMDB.ID)]; !ehFavorito {
			recomendacoes = append(recomendacoes, dominio.Filme{
				ID:             filmeTMDB.ID,
				Titulo:         filmeTMDB.Titulo,
				Sinopse:        filmeTMDB.Sinopse,
				DataLancamento: filmeTMDB.DataLancamento,
				CaminhoPoster:  "https://image.tmdb.org/t/p/w500" + filmeTMDB.CaminhoPoster,
				NotaMedia:      filmeTMDB.NotaMedia,
			})
		}
	}

	return recomendacoes, nil
}
