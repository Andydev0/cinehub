# 🎬 CineHub - Aplicação de Filmes

<div align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB" alt="React">
  <img src="https://img.shields.io/badge/TypeScript-007ACC?style=for-the-badge&logo=typescript&logoColor=white" alt="TypeScript">
  <img src="https://img.shields.io/badge/SQLite-07405E?style=for-the-badge&logo=sqlite&logoColor=white" alt="SQLite">
  <img src="https://img.shields.io/badge/Tailwind_CSS-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white" alt="Tailwind">
</div>

<div align="center">
  <h3>🚀 Uma aplicação completa para descobrir, avaliar e gerenciar seus filmes favoritos</h3>
  <p><em>Projeto desenvolvido para aprimorar habilidades em Go e desenvolvimento full-stack</em></p>
</div>

---

## 📋 Sobre o Projeto

O **CineHub** é uma aplicação web moderna que permite aos usuários descobrir novos filmes, gerenciar suas listas de favoritos, avaliar filmes e participar de um quiz interativo sobre cinema. O projeto foi desenvolvido como uma forma de aprimorar conhecimentos em **Go** no backend e **React/TypeScript** no frontend.

### ✨ Principais Funcionalidades

- 🔍 **Busca Avançada**: Pesquise filmes por título, gênero ou filtros personalizados
- ⭐ **Sistema de Favoritos**: Salve e organize seus filmes preferidos
- 📝 **Avaliações**: Avalie filmes e veja avaliações de outros usuários
- 🎯 **Quiz Interativo**: Teste seus conhecimentos sobre cinema com perguntas personalizadas
- 🎲 **Gerador Aleatório**: Descubra novos filmes com sugestões aleatórias
- 🔐 **Autenticação JWT**: Sistema seguro de login e registro
- 📱 **Design Responsivo**: Interface moderna com tema laranja e glassmorphism

## 🛠️ Tecnologias Utilizadas

### Backend (Go)
- **Gin Framework**: Framework web rápido e minimalista
- **GORM/SQLx**: ORM e query builder para banco de dados
- **SQLite**: Banco de dados leve e eficiente
- **JWT**: Autenticação segura com tokens
- **TMDB API**: Integração com a API do The Movie Database
- **bcrypt**: Hash seguro de senhas

### Frontend (React/TypeScript)
- **React 19**: Biblioteca para interfaces de usuário
- **TypeScript**: Tipagem estática para JavaScript
- **Vite**: Build tool moderna e rápida
- **Tailwind CSS**: Framework CSS utilitário
- **Axios**: Cliente HTTP para requisições
- **React Router**: Roteamento do lado do cliente

## 🚀 Como Executar o Projeto

### Pré-requisitos
- Go 1.24+ instalado
- Node.js 18+ instalado
- Chave da API do TMDB (gratuita)

### 1. Clone o repositório
```bash
git clone https://github.com/Andydev0/cinehub.git
cd cinehub
```

### 2. Configuração do Backend

```bash
cd backend

# Instale as dependências
go mod download

# Configure as variáveis de ambiente
cp .env.example .env
# Edite o arquivo .env com suas configurações
```

### 3. Configuração do Frontend

```bash
cd ../frontend

# Instale as dependências
npm install
```

### 4. Executar a aplicação

**Terminal 1 - Backend:**
```bash
cd backend
go run cmd/api/main.go
```

**Terminal 2 - Frontend:**
```bash
cd frontend
npm run dev
```

A aplicação estará disponível em:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080

## 📁 Estrutura do Projeto

```
cinehub/
├── backend/
│   ├── cmd/api/                 # Ponto de entrada da aplicação
│   ├── internal/
│   │   ├── api/                 # Handlers e rotas
│   │   ├── auth/                # Autenticação JWT
│   │   ├── database/            # Configuração do banco
│   │   ├── dominio/             # Modelos de domínio
│   │   ├── repositorio/         # Camada de dados
│   │   ├── servico/             # Lógica de negócio
│   │   └── tmdb/                # Cliente da API TMDB
│   ├── go.mod
│   └── .env.example
├── frontend/
│   ├── src/
│   │   ├── components/          # Componentes reutilizáveis
│   │   ├── pages/               # Páginas da aplicação
│   │   ├── context/             # Context API do React
│   │   ├── services/            # Serviços de API
│   │   └── types/               # Tipos TypeScript
│   ├── package.json
│   └── tailwind.config.js
└── README.md
```

## 🎯 Funcionalidades Detalhadas

### 🔍 Sistema de Busca
- Busca por título de filme
- Filtros por gênero, ano e avaliação
- Busca avançada com múltiplos critérios
- Resultados paginados e otimizados

### ⭐ Gerenciamento de Favoritos
- Adicionar/remover filmes dos favoritos
- Lista personalizada de filmes salvos
- Avaliações com sistema de estrelas
- Comentários e notas pessoais

### 🎯 Quiz Interativo
- Perguntas sobre diretores, atores, anos e gêneros
- Baseado nos filmes favoritos do usuário
- Sistema de pontuação e estatísticas
- Perguntas variadas e embaralhadas

### 🎲 Gerador Aleatório
- Sugestões personalizadas baseadas no histórico
- Filtros por gênero e década
- Descoberta de filmes menos conhecidos
- Histórico de sugestões

## 🔧 Configuração da API TMDB

1. Acesse [The Movie Database](https://www.themoviedb.org/)
2. Crie uma conta gratuita
3. Vá em Configurações > API
4. Solicite uma chave de API
5. Adicione a chave no arquivo `.env`

## 📊 Endpoints da API

### Autenticação
- `POST /v1/auth/registro` - Registro de usuário
- `POST /v1/auth/login` - Login de usuário

### Filmes
- `GET /v1/filmes/buscar` - Buscar filmes
- `GET /v1/filmes/detalhes/:id` - Detalhes do filme
- `GET /v1/filmes/genero` - Buscar por gênero
- `GET /v1/filmes/aleatorio` - Filme aleatório

### Favoritos
- `GET /v1/favoritos` - Listar favoritos
- `POST /v1/favoritos` - Adicionar favorito
- `DELETE /v1/favoritos/:id` - Remover favorito

### Quiz
- `GET /v1/quiz/pergunta` - Gerar pergunta
- `POST /v1/quiz/resposta` - Enviar resposta

## 🎨 Design e UX

O projeto utiliza um design moderno com:
- **Tema laranja** como cor principal
- **Glassmorphism** para cards e elementos
- **Animações suaves** e transições
- **Layout responsivo** para todos os dispositivos
- **Ícones emoji** para melhor identificação visual
- **Estados de loading** e feedback visual


---

<div align="center">
  <p>⭐ Se este projeto te ajudou, considere dar uma estrela!</p>
  <p>💡 Desenvolvido com 💜 para aprimorar habilidades em Go e React</p>
</div>
