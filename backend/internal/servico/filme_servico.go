package servico

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/Andydev0/filmes-backend/internal/dominio"
)

// Estrutura para decodificar a resposta da lista de gêneros da API.
type ListaGenerosResponse struct {
	Generos []dominio.Genero `json:"genres"`
}

// A interface agora reflete que BuscarDetalhes retorna nossa nova struct completa.
type FilmeServico interface {
	BuscarFilmes(termo string) ([]dominio.Filme, error)
	BuscarFilmesPorGenero(generoID string) ([]dominio.Filme, error)
	BuscarDetalhes(filmeID int64) (*dominio.DetalhesFilmeCompleto, error)
	ListarGeneros() ([]dominio.Genero, error)
	BuscarFilmeAleatorio(generoID, ano string) (*dominio.Filme, error)
}

type tmdbService struct {
	apiKey      string
	clienteHttp *http.Client
}

func NovoFilmeServico(chaveAPI string) FilmeServico {
	return &tmdbService{
		apiKey:      chaveAPI,
		clienteHttp: &http.Client{},
	}
}

func (s *tmdbService) BuscarFilmes(termo string) ([]dominio.Filme, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("a chave da API do TMDB não foi configurada")
	}

	baseURL := "https://api.themoviedb.org/3/search/movie"
	urlBusca := fmt.Sprintf("%s?api_key=%s&query=%s&language=pt-BR", baseURL, s.apiKey, url.QueryEscape(termo))

	resposta, err := s.clienteHttp.Get(urlBusca)
	if err != nil {
		return nil, fmt.Errorf("falha ao realizar a requisição para o TMDB: %w", err)
	}
	defer resposta.Body.Close()

	if resposta.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("a API do TMDB retornou um status inesperado: %s", resposta.Status)
	}

	corpo, err := ioutil.ReadAll(resposta.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler o corpo da resposta: %w", err)
	}

	var respostaTMDB dominio.RespostaBuscaTMDB
	if err := json.Unmarshal(corpo, &respostaTMDB); err != nil {
		return nil, fmt.Errorf("falha ao decodificar o JSON da resposta: %w", err)
	}

	filmes := make([]dominio.Filme, 0, len(respostaTMDB.Resultados))
	for _, filmeTMDB := range respostaTMDB.Resultados {
		filmes = append(filmes, dominio.Filme{
			ID:             filmeTMDB.ID,
			Titulo:         filmeTMDB.Titulo,
			Sinopse:        filmeTMDB.Sinopse,
			DataLancamento: filmeTMDB.DataLancamento,
			CaminhoPoster:  "https://image.tmdb.org/t/p/w500" + filmeTMDB.CaminhoPoster,
			NotaMedia:      filmeTMDB.NotaMedia,
		})
	}

	return filmes, nil
}

// Método BuscarDetalhes refatorado para buscar detalhes, créditos e vídeos.
func (s *tmdbService) BuscarDetalhes(filmeID int64) (*dominio.DetalhesFilmeCompleto, error) {
	var detalhes dominio.TMDBMovieResult
	var creditos dominio.CreditosTMDB
	var videos dominio.VideosTMDB
	var wg sync.WaitGroup
	var errGlobal error

	wg.Add(3)

	// Goroutine 1: Busca detalhes básicos
	go func() {
		defer wg.Done()
		url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d?api_key=%s&language=pt-BR", filmeID, s.apiKey)
		resp, err := s.clienteHttp.Get(url)
		if err == nil {
			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(&detalhes)
		} else {
			errGlobal = err
		}
	}()

	// Goroutine 2: Busca créditos (elenco/diretor/escritor)
	go func() {
		defer wg.Done()
		url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/credits?api_key=%s&language=pt-BR", filmeID, s.apiKey)
		resp, err := s.clienteHttp.Get(url)
		if err == nil {
			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(&creditos)
		} else {
			errGlobal = err
		}
	}()

	// Goroutine 3: Busca vídeos (trailer)
	go func() {
		defer wg.Done()
		url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/videos?api_key=%s&language=pt-BR", filmeID, s.apiKey)
		resp, err := s.clienteHttp.Get(url)
		if err == nil {
			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(&videos)
		} else {
			errGlobal = err
		}
	}()

	wg.Wait()

	if errGlobal != nil {
		return nil, fmt.Errorf("falha ao buscar dados do TMDB")
	}

	respostaFinal := &dominio.DetalhesFilmeCompleto{
		TMDBMovieResult: &detalhes,
		Elenco:          []dominio.MembroElenco{},
		Escritores:      []string{},
	}

	if len(creditos.Elenco) > 10 {
		respostaFinal.Elenco = creditos.Elenco[:10]
	} else {
		respostaFinal.Elenco = creditos.Elenco
	}

	// Encontra o diretor e os escritores na lista da equipe.
	for _, membro := range creditos.Equipe {
		// CORREÇÃO: Usamos 'membro.Nome' em vez de 'membro.Name'
		if membro.Job == "Director" {
			respostaFinal.Diretor = membro.Nome
		}
		// CORREÇÃO: Usamos 'membro.Nome' em vez de 'membro.Name'
		if membro.Job == "Screenplay" || membro.Job == "Writer" || membro.Job == "Story" {
			respostaFinal.Escritores = append(respostaFinal.Escritores, membro.Nome)
		}
	}

	// Encontra a chave do trailer oficial no YouTube.
	for _, video := range videos.Resultados {
		if video.Site == "YouTube" && video.Tipo == "Trailer" {
			respostaFinal.TrailerKey = video.Key
			break
		}
	}

	return respostaFinal, nil
}

// ListarGeneros busca a lista de todos os gêneros de filmes disponíveis.
func (s *tmdbService) ListarGeneros() ([]dominio.Genero, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/genre/movie/list?api_key=%s&language=pt-BR", s.apiKey)
	resp, err := s.clienteHttp.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var listaGeneros ListaGenerosResponse
	if err := json.NewDecoder(resp.Body).Decode(&listaGeneros); err != nil {
		return nil, err
	}
	return listaGeneros.Generos, nil
}

// BuscarFilmeAleatorio usa a API Discover para encontrar um filme com base nos filtros.
func (s *tmdbService) BuscarFilmeAleatorio(generoID, ano string) (*dominio.Filme, error) {
	baseURL := "https://api.themoviedb.org/3/discover/movie"
	queryParams := url.Values{}
	queryParams.Add("api_key", s.apiKey)
	queryParams.Add("language", "pt-BR")
	queryParams.Add("sort_by", "popularity.desc")

	if generoID != "" {
		queryParams.Add("with_genres", generoID)
	}
	if ano != "" {
		queryParams.Add("primary_release_year", ano)
	}

	paginaAleatoria := rand.Intn(50) + 1
	queryParams.Add("page", strconv.Itoa(paginaAleatoria))

	resp, err := s.clienteHttp.Get(baseURL + "?" + queryParams.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var discoverResponse dominio.RespostaBuscaTMDB
	if err := json.NewDecoder(resp.Body).Decode(&discoverResponse); err != nil {
		return nil, err
	}

	if len(discoverResponse.Resultados) == 0 {
		return nil, fmt.Errorf("nenhum filme encontrado com os filtros fornecidos")
	}

	filmeAleatorioTMDB := discoverResponse.Resultados[rand.Intn(len(discoverResponse.Resultados))]

	filme := &dominio.Filme{
		ID:             filmeAleatorioTMDB.ID,
		Titulo:         filmeAleatorioTMDB.Titulo,
		Sinopse:        filmeAleatorioTMDB.Sinopse,
		DataLancamento: filmeAleatorioTMDB.DataLancamento,
		CaminhoPoster:  "https://image.tmdb.org/t/p/w500" + filmeAleatorioTMDB.CaminhoPoster,
		NotaMedia:      filmeAleatorioTMDB.NotaMedia,
	}

	return filme, nil
}

// BuscarFilmesPorGenero busca filmes de um gênero específico usando a API Discover
func (s *tmdbService) BuscarFilmesPorGenero(generoID string) ([]dominio.Filme, error) {
	baseURL := "https://api.themoviedb.org/3/discover/movie"
	queryParams := url.Values{}
	queryParams.Add("api_key", s.apiKey)
	queryParams.Add("language", "pt-BR")
	queryParams.Add("sort_by", "popularity.desc")
	queryParams.Add("with_genres", generoID)

	resp, err := s.clienteHttp.Get(baseURL + "?" + queryParams.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var discoverResponse dominio.RespostaBuscaTMDB
	if err := json.NewDecoder(resp.Body).Decode(&discoverResponse); err != nil {
		return nil, err
	}

	filmes := make([]dominio.Filme, 0, len(discoverResponse.Resultados))
	for _, filmeTMDB := range discoverResponse.Resultados {
		filme := &dominio.Filme{
			ID:             filmeTMDB.ID,
			Titulo:         filmeTMDB.Titulo,
			Sinopse:        filmeTMDB.Sinopse,
			DataLancamento: filmeTMDB.DataLancamento,
			CaminhoPoster:  "https://image.tmdb.org/t/p/w500" + filmeTMDB.CaminhoPoster,
			NotaMedia:      filmeTMDB.NotaMedia,
		}
		filmes = append(filmes, *filme)
	}

	return filmes, nil
}
