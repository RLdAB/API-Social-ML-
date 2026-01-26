# API-Social-ML - Documentação Completa

## 📋 Visão Geral

A **API Social-ML** é uma API REST desenvolvida em **Go** (Golang) que implementa funcionalidades de rede social voltadas para um marketplace, permitindo que vendedores publiquem produtos e promoções, enquanto usuários possam seguir vendedores e visualizar seus posts.

**Stack Tecnológico:**
- 🔧 **Framework**: Chi (router HTTP)
- 🗄️ **ORM**: GORM (Object-Relational Mapping)
- 📚 **Documentação**: Swagger/OpenAPI
- 🐳 **Containerização**: Docker

---

## 🏗️ Arquitetura da Aplicação

A API segue **Clean Architecture** com separação clara de responsabilidades:

```
┌─────────────────────────────────────────────────────────────┐
│              HTTP HANDLERS (API Layer)                       │
│         UserHandlers / PostHandlers                          │
└──────────────────┬──────────────────────────────────────────┘
                   ↓
┌─────────────────────────────────────────────────────────────┐
│         APPLICATION SERVICES (Business Logic)               │
│  UserService / FollowService / PostService                  │
└──────────────────┬──────────────────────────────────────────┘
                   ↓
┌─────────────────────────────────────────────────────────────┐
│              DOMAIN (Core Business Rules)                    │
│  User / Post / Follow entities + Interfaces                 │
└──────────────────┬──────────────────────────────────────────┘
                   ↓
┌─────────────────────────────────────────────────────────────┐
│      INFRASTRUCTURE/PERSISTENCE (Database Layer)            │
│     GormUserRepository / GormPostRepository                 │
└─────────────────────────────────────────────────────────────┘
```

**Fluxo de uma Requisição:**
1. **HTTP Handler** recebe a request
2. **Service** valida e aplica regras de negócio
3. **Domain** define estruturas e erros
4. **Repository** persiste/recupera dados do banco

---

## 📡 Endpoints Principais

### A. Gestão de Usuários

| Função | Endpoint | Método | Objetivo |
|--------|----------|--------|----------|
| Criar Usuário | `/users` | POST | Registrar novo vendedor ou comprador |
| Listar Usuários | `/users` | GET | Obter lista de todos os usuários |
| Buscar Usuário | `/users/{userId}` | GET | Recuperar dados específicos de um usuário |
| Atualizar Usuário | `/users/{userId}` | PUT | Modificar informações do usuário |

### B. Sistema de Follow/Unfollow

| Função | Endpoint | Método | Objetivo |
|--------|----------|--------|----------|
| Seguir Vendedor | `/users/{userId}/follow/{sellerId}` | POST | Usuário segue um vendedor |
| Deixar de Seguir | `/users/{userId}/follow/{sellerId}` | PUT | Usuário deixa de seguir vendedor |

**Regras de Negócio:**
- ❌ Não pode seguir a si mesmo
- ❌ Só pode seguir usuários marcados como `is_seller = true`
- ❌ Não pode seguir o mesmo vendedor duas vezes

### C. Seguidores e Seguindo

| Função | Endpoint | Método | Objetivo |
|--------|----------|--------|----------|
| Contar Seguidores | `/users/{userId}/followers/count` | GET | Obter quantidade de followers |
| Listar Seguidores | `/users/{userId}/followers/list` | GET | Listar quem segue este usuário (com ordenação) |
| Listar Seguindo | `/users/{userId}/following/list` | GET | Listar vendedores que o usuário segue |

**Parâmetros de Ordenação:**
- `order=name_asc` → Alfabético crescente
- `order=name_desc` → Alfabético decrescente

### D. Publicação de Produtos

| Função | Endpoint | Método | Objetivo |
|--------|----------|--------|----------|
| Publicar Produto | `/products/publish` | POST | Vendedor publica um produto normal |
| Publicar Promoção | `/products/promo-pub` | POST | Vendedor publica com desconto ativo |
| Listar Promoções | `/products/promo-pub/list` | GET | Listar posts em promoção de um vendedor |

### E. Feed de Seguidos

| Função | Endpoint | Método | Objetivo |
|--------|----------|--------|----------|
| Feed Recente | `/products/followed/latest/{userId}` | GET | Posts dos últimos 2 produtos dos vendedores seguidos |

**Parâmetros:**
- `order=date_asc` → Mais antigos primeiro
- `order=date_desc` → Mais recentes primeiro

### F. Métricas de Promoção

| Função | Endpoint | Método | Objetivo |
|--------|----------|--------|----------|
| Contar Promoções | `/sellers/{sellerId}/promotions/count` | GET | Quantidade de posts em promoção de um vendedor |

---

## 📊 Entidades e Modelos

### User (Usuário)

```go
type User struct {
    ID         uint              // Identificador único
    CreatedAt  time.Time         // Data de criação
    UpdatedAt  time.Time         // Data de atualização
    DeletedAt  gorm.DeletedAt    // Soft Delete
    Name       string            // Nome (max 15 chars)
    IsSeller   bool              // Flag: é vendedor?
}
```

### Follow (Relacionamento de Seguimento)

```go
type Follow struct {
    FollowerID uint   // Quem segue (FK para User)
    SellerID   uint   // Quem é seguido (FK para User)
}
```

**Características:**
- Chave primária composta: `(FollowerID, SellerID)`
- Representa a relação "um usuário segue um vendedor"

### Post (Publicação)

```go
type Post struct {
    ID            uint          // Identificador único
    CreatedAt     time.Time     // Data de criação
    UpdatedAt     time.Time     // Data de atualização
    DeletedAt     gorm.DeletedAt
    
    UserID        uint          // FK para User (vendedor)
    ProductName   string        // Nome do produto
    ProductType   string        // Tipo (ex: Eletrônicos)
    ProductBrand  string        // Marca
    Category      string        // Categoria
    Content       string        // Descrição completa
    Price         float64       // Preço (DECIMAL 10,2)
    
    HasPromo      bool          // Tem promoção?
    Discount      float64       // % desconto (0-100)
    PromoEndsAt   time.Time     // Quando termina promoção
}
```

---

## 🔗 Relacionamentos e Cardinalidade

```
┌──────────────┐         ┌──────────────┐
│    User      │         │   Follow     │
├──────────────┤         ├──────────────┤
│ ID (PK)      │◄────┐   │ FollowerID   │
│ Name         │     └───┤ (FK, PK)     │
│ IsSeller     │         │              │
└──────────────┘      ┌──┤ SellerID     │
       ▲              │  │ (FK, PK)     │
       │              │  └──────────────┘
       │              │
       │              └──► User (Sellers)
       │
       │ 1:N (um usuário publica muitos posts)
       │
    ┌──┴────────────────┐
    │      Post         │
    ├───────────────────┤
    │ ID (PK)           │
    │ UserID (FK)       │
    │ ProductName       │
    │ Price             │
    │ HasPromo          │
    │ Discount          │
    └───────────────────┘
```

| Relação | Cardinalidade | Descrição |
|---------|---------------|-----------|
| User → Follow | 1:N | Um usuário pode seguir múltiplos vendedores |
| Follow → User | N:M | Muitos usuários podem seguir o mesmo vendedor |
| User → Post | 1:N | Um vendedor publica muitos posts |
| Post → User | N:1 | Cada post pertence a um único vendedor |

---

## 💡 Melhorias Propostas

### Funcionalidades Ausentes
- ⭕ **Comentários/Reviews**: Usuários não podem avaliar ou comentar posts
- ⭕ **Sistema de Curtidas**: Não há funcionalidade de like/favorito
- ⭕ **Busca Avançada**: Sem filtros ou busca por nome/categoria
- ⭕ **Autenticação/Autorização**: Sem JWT ou OAuth
- ⭕ **Paginação**: Endpoints não suportam `limit` e `offset`
- ⭕ **Notificações**: Sem sistema de notificações

### Validações e Segurança
- 🔒 Validação de Email único
- 🔒 Hash de Senha
- 🔒 Rate Limiting
- 🔒 Configuração de CORS
- 🔒 Validação avançada de entrada

### Performance e Escalabilidade
- ⚡ Cache para seguidores/posts frequentes
- ⚡ Índices compostos em `follows` e `posts`
- ⚡ Paginação com Cursor
- ⚡ Denormalização de contagens

### Observabilidade
- 📊 Logging Estruturado
- 📊 Métricas (Prometheus)
- 📊 Tracing Distribuído
- 📊 Health Checks

### Arquitetura
- 🏗️ Testes Unitários
- 🏗️ Integração com Message Queue
- 🏗️ Versionamento de API (`/v1`, `/v2`)
- 🏗️ Melhorar Soft Delete

---

## 📂 Estrutura de Pastas

```
├── cmd/server/              # Ponto de entrada
│   ├── main.go
│   └── routes.go
├── internal/                # Código privado
│   ├── user/
│   │   ├── application/
│   │   ├── domain/
│   │   └── infrastructure/
│   └── post/
│       ├── application/
│       ├── domain/
│       └── infrastructure/
├── docs/                    # Swagger
├── migrations/              # SQL
└── go.mod
```

---

## 🚀 Como Executar

```bash
# Instalar dependências
go mod download

# Configurar variáveis de ambiente
export DATABASE_URL="postgres://user:pass@localhost/socialml"
export PORT=8080

# Executar
go run cmd/server/main.go

# Acessar Swagger
# http://localhost:8080/swagger/
```

---

**Última atualização**: 26/01/2026
**Versão**: 1.0
