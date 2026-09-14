# Job Bot

Bot open-source para Discord focado em **descobrir, organizar e entregar oportunidades profissionais de forma contínua**.

O Job Bot foi pensado como um produto de uso real: ele conecta fontes de oportunidades, persiste dados no MongoDB, mantém métricas da operação e entrega essas informações pelo Discord.

> **Status:** em desenvolvimento ativo. A base de Discord, MongoDB, Active Record, stats e presence dinâmica já está sendo estruturada, enquanto os módulos de fontes e oportunidades evoluem.

## O que o Job Bot quer resolver

Encontrar boas oportunidades exige acompanhar várias fontes, repetir buscas, filtrar resultados e evitar duplicações.

O Job Bot centraliza esse fluxo:

```text
Fontes de oportunidades
        ↓
Coleta e normalização
        ↓
MongoDB
        ↓
Filtros e regras
        ↓
Discord
```

A proposta é transformar o Discord em uma interface prática para acompanhar oportunidades sem depender de buscas manuais constantes.

## Funcionalidades

### Já disponível / em implementação

- integração com Discord através de `discordgo`;
- slash commands;
- handlers de eventos;
- conexão persistente com MongoDB;
- models com Active Record;
- criação automática de índices;
- estatísticas da aplicação;
- presence dinâmica;
- scheduler periódico para atualização da presence;
- startup e shutdown coordenados;
- configuração por variáveis de ambiente.

### Em desenvolvimento

- cadastro e gerenciamento de fontes;
- coleta automática de oportunidades;
- model de oportunidades;
- deduplicação;
- filtros por área, tecnologia e localização;
- preferências por usuário;
- notificações automáticas;
- comandos para consulta;
- novas métricas operacionais.

## Exemplo de uso

A ideia é que o bot consiga acompanhar o estado da própria operação e, futuramente, entregar comandos e notificações como:

```text
/novas
/fontes
/stats
```

Além disso, a presence do bot pode refletir métricas reais:

```text
42 usuários cadastrados
18 fontes monitoradas
1.284 oportunidades encontradas
```

## Stack

- Go
- DiscordGo
- MongoDB
- MongoDB Go Driver v2
- Docker / Docker Compose

## Executando localmente

### 1. Configure o ambiente

Copie:

```bash
cp .env.example .env
```

Preencha:

```env
DISCORD_TOKEN=
DISCORD_APP_ID=

MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=job_bot
```

### 2. Suba o MongoDB

```bash
docker compose up -d
```

### 3. Execute o bot

```bash
go run ./cmd/job-bot
```

## Como o startup funciona

O bot só fica disponível como pronto depois que suas dependências essenciais estiverem carregadas:

```text
Configuração
      ↓
Discord conectado
      ↓
idle / "inicializando"
      ↓
MongoDB conectado
      ↓
Índices verificados
      ↓
Slash commands sincronizados
      ↓
Stats carregados
      ↓
online / presence dinâmica
      ↓
Schedulers iniciados
```

Se uma dependência essencial falhar, a aplicação encerra em vez de continuar parcialmente funcional.

## Estrutura do repositório

```text
job-bot/
├── cmd/
│   └── job-bot/
├── internal/
│   ├── config/
│   ├── database/
│   ├── discord/
│   ├── stats/
│   └── users/
├── docs/
├── compose.yaml
├── .env.example
├── LICENSE
├── README.md
└── DEVELOPMENT.md
```

## Desenvolvimento

A documentação técnica, padrões de código e instruções para criar commands, handlers, models e integrações ficam em:

**[`DEVELOPMENT.md`](DEVELOPMENT.md)**

A documentação complementar do projeto fica em:

**[`docs/`](docs/)**

## Contribuindo

Contribuições são bem-vindas.

Antes de abrir um Pull Request:

1. mantenha a mudança focada;
2. siga as convenções descritas em `DEVELOPMENT.md`;
3. formate o código com `go fmt`;
4. execute `go test ./...`;
5. execute `go vet ./...`;
6. atualize a documentação quando necessário.

Para mudanças maiores, prefira abrir uma issue antes de implementar.

## Segurança

Nunca publique:

- token do Discord;
- credenciais do MongoDB;
- `.env`;
- chaves de APIs;
- secrets de produção.

Se um segredo for exposto, faça a rotação imediatamente.

## Licença

Distribuído sob a licença **MIT**.

Consulte [`LICENSE`](LICENSE).
