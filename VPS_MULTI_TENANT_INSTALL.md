# Instalacao em VPS Limpa para SaaS de Lojas

Este guia assume que voce vai:

1. terminar primeiro a aplicacao das lojas;
2. depois subir tudo em uma VPS limpa;
3. operar um SaaS multi-tenant com:
   - frontend React;
   - backend Go;
   - PostgreSQL;
   - Redis;
   - Traefik na frente;
   - dominios proprios de clientes;
   - dois DNS seus, por exemplo `dns1.vrprime.com.br` e `dns2.vrprime.com.br`.

O objetivo final e este:

- a loja continua existindo internamente por `slug`;
- o cliente pode usar `rainha.com` ou `www.rainha.com`;
- o trafego entra pelo Traefik;
- o backend identifica a loja pelo `Host`;
- o site abre no dominio do cliente, sem mostrar `/slug`.

## 1. Decisao de arquitetura

Para este projeto, a base recomendada e:

- `Ubuntu 24.04 LTS`
- `Docker Engine`
- `Docker Compose`
- `Traefik v3`
- `PostgreSQL 16`
- `Redis 7`
- sua app atual dividida em:
  - `frontend` com Nginx interno;
  - `api` Go;
  - `postgres`;
  - `redis`.

Este guia assume instalacao limpa, sem Dokploy.

Motivo:

- dominio customizado de cliente exige mais controle do proxy;
- SSL por dominio exige automacao previsivel;
- o Traefik precisa ser configurado com mais liberdade do que normalmente se usa num painel de deploy.

## 2. Como o sistema vai funcionar

### 2.1 Loja por slug

Cada loja tem um `slug`, por exemplo:

- `rainha`
- `lojadomario`
- `ofertas-centro`

Isso continua sendo o identificador interno.

### 2.2 Dominio do cliente

O cliente cadastra no painel algo como:

- `rainha.com`
- `www.rainha.com`

Voce salva esse dominio na base e vincula a uma loja.

Exemplo de mapeamento:

- `rainha.com -> loja slug rainha`
- `www.rainha.com -> loja slug rainha`
- `lojadomario.com.br -> loja slug lojadomario`

### 2.3 DNS `dns1` e `dns2`

Voce quer dois alvos seus:

- `dns1.vrprime.com.br`
- `dns2.vrprime.com.br`

Esses nomes vao apontar para sua infraestrutura. O cliente usa esses alvos no DNS dele.

Exemplos:

- `www.rainha.com CNAME dns1.vrprime.com.br`
- `www.lojamaria.com CNAME dns2.vrprime.com.br`

Para dominio raiz `@`:

- muitos provedores nao aceitam `CNAME` no apex;
- nesse caso o cliente vai usar `A`, `ALIAS` ou `ANAME`.

Exemplo:

- `rainha.com A IP_DO_PROXY`
- ou `rainha.com ALIAS dns1.vrprime.com.br`

## 3. Ponto importante sobre SSL

Aqui esta a parte que mais costuma confundir.

### 3.1 Wildcard do seu dominio

Voce consegue cobrir estes hosts com certificado wildcard:

- `*.vrprime.com.br`
- `vrprime.com.br`

Isso serve para:

- `dns1.vrprime.com.br`
- `dns2.vrprime.com.br`
- `loja123.vrprime.com.br`

Para wildcard, o Traefik usa `DNS-01 challenge`.

### 3.2 Dominio do cliente

O wildcard `*.vrprime.com.br` **nao** cobre:

- `rainha.com`
- `www.rainha.com`
- `lojamaria.com.br`

Cada dominio de cliente precisa do proprio certificado.

No Traefik, certificados automaticos sao emitidos para os dominios presentes na regra `Host(...)` do router, ou definidos explicitamente em `tls.domains`.

Traduzindo para sua operacao:

- um router generico `HostRegexp` resolve roteamento;
- mas dominio de cliente com SSL automatico precisa entrar na configuracao do Traefik como dominio conhecido.

Por isso, o onboarding de dominio customizado precisa:

1. validar DNS;
2. salvar no banco;
3. gerar ou atualizar configuracao dinamica do Traefik para aquele dominio;
4. deixar o Traefik emitir o certificado.

## 4. Ordem recomendada de implementacao

Nao tente fazer tudo de uma vez.

### Fase 1

Suba a plataforma com:

- lojas por `slug`;
- um dominio principal seu;
- deploy estavel;
- backup;
- monitoramento basico.

### Fase 2

Adicione:

- tabela de dominios customizados;
- onboarding de dominio;
- verificacao DNS;
- resolucao por `Host`;
- configuracao dinamica do Traefik.

### Fase 3

Se o volume crescer:

- `dns1` e `dns2` em servidores diferentes;
- replica de banco ou backup mais robusto;
- cache e rate limit melhores;
- CDN;
- observabilidade.

## 5. Requisitos da VPS

Para comecar com seguranca:

- `2 vCPU`
- `4 GB RAM`
- `80 GB SSD`
- Ubuntu 24.04

Se voce pretende subir varias lojas pequenas, isso segura o MVP.

Para um ambiente mais folgado:

- `4 vCPU`
- `8 GB RAM`
- `160 GB SSD`

## 6. DNS que voce precisa preparar

No painel DNS da `vrprime.com.br`, crie:

- `A dns1.vrprime.com.br -> IP_PUBLICO_DA_VPS`
- `A dns2.vrprime.com.br -> IP_PUBLICO_DA_VPS_2`

Se no inicio voce tiver apenas uma VPS, pode fazer:

- `A dns1.vrprime.com.br -> IP_PUBLICO_DA_VPS`
- `A dns2.vrprime.com.br -> IP_PUBLICO_DA_VPS`

Opcional, mas util:

- `A app.vrprime.com.br -> IP_PUBLICO_DA_VPS`
- `A api.vrprime.com.br -> IP_PUBLICO_DA_VPS`

Se quiser subdominios seus por loja:

- `A *.vrprime.com.br -> IP_PUBLICO_DA_VPS`

## 7. Preparacao inicial da VPS

Entre na VPS por SSH.

Atualize o sistema:

```bash
sudo apt update && sudo apt upgrade -y
```

Instale pacotes basicos:

```bash
sudo apt install -y ca-certificates curl gnupg lsb-release git unzip
```

Configure timezone:

```bash
sudo timedatectl set-timezone America/Sao_Paulo
```

Crie um usuario de deploy, se quiser separar do root:

```bash
sudo adduser deploy
sudo usermod -aG sudo deploy
```

## 8. Instalacao do Docker

Remova pacotes antigos:

```bash
sudo apt remove -y docker docker-engine docker.io containerd runc
```

Adicione a chave oficial:

```bash
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
```

Adicione o repositorio:

```bash
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```

Instale:

```bash
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

Permita uso sem `sudo`:

```bash
sudo usermod -aG docker $USER
newgrp docker
```

Teste:

```bash
docker version
docker compose version
```

## 9. Estrutura de diretorios no servidor

Sugestao:

```text
/srv/goapi/
  app/
    backend/
    frontend/
    docker-compose.prod.yml
    .env
  traefik/
    docker-compose.yml
    traefik.yml
    dynamic/
    letsencrypt/
```

Crie:

```bash
sudo mkdir -p /srv/goapi/app
sudo mkdir -p /srv/goapi/traefik/dynamic
sudo mkdir -p /srv/goapi/traefik/letsencrypt
sudo chown -R $USER:$USER /srv/goapi
```

## 10. Clonar o projeto

```bash
cd /srv/goapi/app
git clone SEU_REPOSITORIO .
```

Se for privado, configure chave SSH ou token antes.

## 11. Variaveis de ambiente de producao

Crie o arquivo:

```bash
cd /srv/goapi/app
cp backend/.env.example .env 2>/dev/null || touch .env
```

Exemplo recomendado:

```env
APP_ENV=production
HTTP_ADDR=:8080

POSTGRES_DB=goapi_prod
POSTGRES_USER=goapi_prod
POSTGRES_PASSWORD=troque_essa_senha
DATABASE_URL=postgres://goapi_prod:troque_essa_senha@postgres:5432/goapi_prod?sslmode=disable

REDIS_URL=redis://redis:6379/0

JWT_SECRET=troque_por_uma_chave_muito_longa
JWT_TTL=24h
BCRYPT_COST=12
DB_MAX_CONNS=40
DB_MIN_CONNS=5

SMTP_HOST=smtp.seudominio.com
SMTP_PORT=587
SMTP_USERNAME=seu_usuario
SMTP_PASSWORD=sua_senha
SMTP_FROM=noreply@vrprime.com.br
```

Gere um segredo forte:

```bash
openssl rand -base64 64
```

## 12. Rede Docker de producao

Crie uma rede compartilhada para o proxy:

```bash
docker network create web
```

Essa rede sera usada pelo Traefik e pelo frontend.

## 13. Traefik em producao

### 13.1 Arquivo `traefik.yml`

Crie `/srv/goapi/traefik/traefik.yml`:

```yaml
api:
  dashboard: true

entryPoints:
  web:
    address: ":80"
  websecure:
    address: ":443"

providers:
  docker:
    exposedByDefault: false
    network: web
  file:
    directory: /dynamic
    watch: true

certificatesResolvers:
  letsencrypt:
    acme:
      email: [email protected]
      storage: /letsencrypt/acme.json
      httpChallenge:
        entryPoint: web

log:
  level: INFO
```

Observacao:

- esse modelo usa `HTTP-01` para dominios publicos que apontam para a VPS;
- para wildcard `*.vrprime.com.br`, troque depois para `DNS-01`.

### 13.2 Arquivo `docker-compose.yml` do Traefik

Crie `/srv/goapi/traefik/docker-compose.yml`:

```yaml
services:
  traefik:
    image: traefik:v3.1
    container_name: traefik
    restart: unless-stopped
    command:
      - --configFile=/etc/traefik/traefik.yml
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./traefik.yml:/etc/traefik/traefik.yml:ro
      - ./dynamic:/dynamic
      - ./letsencrypt:/letsencrypt
    networks:
      - web

networks:
  web:
    external: true
```

Suba o Traefik:

```bash
cd /srv/goapi/traefik
touch letsencrypt/acme.json
chmod 600 letsencrypt/acme.json
docker compose up -d
```

Verifique:

```bash
docker ps
docker logs traefik --tail 100
```

Recomendacao:

- no inicio, rode apenas uma instancia do Traefik;
- nao tente compartilhar `acme.json` entre varios proxies sem um desenho de HA bem definido.

## 14. Compose de producao da aplicacao

Crie `/srv/goapi/app/docker-compose.prod.yml`.

Exemplo:

```yaml
services:
  postgres:
    image: postgres:16-alpine
    restart: unless-stopped
    env_file:
      - .env
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backend/migrations:/docker-entrypoint-initdb.d:ro
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}"]
      interval: 5s
      timeout: 3s
      retries: 10

  redis:
    image: redis:7-alpine
    restart: unless-stopped
    volumes:
      - redis_data:/data

  api:
    build: ./backend
    restart: unless-stopped
    env_file:
      - .env
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started

  frontend:
    build: ./frontend
    restart: unless-stopped
    depends_on:
      - api
    networks:
      - default
      - web
    labels:
      - traefik.enable=true
      - traefik.docker.network=web
      - traefik.http.routers.goapi-frontend.rule=Host(`app.vrprime.com.br`)
      - traefik.http.routers.goapi-frontend.entrypoints=websecure
      - traefik.http.routers.goapi-frontend.tls.certresolver=letsencrypt
      - traefik.http.services.goapi-frontend.loadbalancer.server.port=80

volumes:
  postgres_data:
  redis_data:

networks:
  web:
    external: true
```

Observacoes:

- `postgres`, `redis` e `api` nao expoem portas para a internet;
- quem fica publico e apenas o `frontend`;
- o `frontend/nginx.conf` atual ja faz proxy de `/api/` para `api:8080`.

## 15. Primeiro deploy

Suba a aplicacao:

```bash
cd /srv/goapi/app
docker compose -f docker-compose.prod.yml up -d --build
```

Verifique:

```bash
docker compose -f docker-compose.prod.yml ps
docker compose -f docker-compose.prod.yml logs -f
```

Teste:

- `https://app.vrprime.com.br`
- `https://app.vrprime.com.br/api/v1/healthz`

## 16. Firewall e seguranca basica

Ative UFW:

```bash
sudo ufw allow OpenSSH
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

Confirme:

```bash
sudo ufw status
```

Nao publique:

- `5432`
- `6379`
- `8080`

## 17. Modelo de dominio customizado no banco

Antes de suportar dominio de cliente, crie uma tabela dedicada.

Exemplo:

```sql
CREATE TABLE store_domains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    domain VARCHAR(255) NOT NULL UNIQUE,
    verified BOOLEAN NOT NULL DEFAULT false,
    dns_target VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Uso sugerido:

- `domain`: `rainha.com`
- `verified`: se o DNS ja aponta para sua infra
- `dns_target`: `dns1.vrprime.com.br` ou `dns2.vrprime.com.br`

## 18. Onboarding de dominio do cliente

Fluxo recomendado no painel:

1. usuario informa o dominio;
2. sistema normaliza para lowercase;
3. sistema salva como `pending`;
4. painel mostra instrucoes de DNS;
5. backend verifica se o dominio aponta para sua infra;
6. se validar, marca como `verified`;
7. sistema escreve a configuracao dinamica do Traefik;
8. Traefik recarrega;
9. certificado e emitido;
10. dominio entra em producao.

## 19. Como pedir DNS para o cliente

### 19.1 Caso simples com `www`

Voce pode pedir:

- `www.rainha.com CNAME dns1.vrprime.com.br`

Esse e o caso mais simples.

### 19.2 Caso com dominio raiz `@`

Voce pode pedir um destes:

- `rainha.com A IP_PUBLICO_DA_VPS`
- `rainha.com ALIAS dns1.vrprime.com.br`
- `rainha.com ANAME dns1.vrprime.com.br`

Depende do provedor DNS do cliente.

### 19.3 Recomendacao operacional

No inicio, suporte primeiro:

- `www.cliente.com`

Depois adicione:

- `cliente.com`

Motivo:

- `www` com `CNAME` e muito mais previsivel;
- apex/root domain varia bastante de provedor.

## 20. Configuracao dinamica do Traefik para dominios de clientes

Quando o cliente `rainha.com` for validado, voce precisa adicionar uma configuracao dinamica.

Exemplo de arquivo:

`/srv/goapi/traefik/dynamic/rainha.com.yml`

```yaml
http:
  routers:
    store-rainha:
      rule: "Host(`rainha.com`) || Host(`www.rainha.com`)"
      entryPoints:
        - websecure
      tls:
        certResolver: letsencrypt
      service: goapi-frontend

  services:
    goapi-frontend:
      loadBalancer:
        servers:
          - url: "http://frontend:80"
```

Observacoes importantes:

- a ideia acima funciona melhor se o Traefik enxergar o container `frontend` na rede `web`;
- voce pode preferir manter esse service no provider Docker e escrever apenas routers no provider File;
- o principal e: cada dominio real precisa entrar na configuracao do Traefik para emissao automatica de certificado.

Sugestao pratica:

- gere um arquivo por dominio ou por loja;
- use nomes estaveis, por exemplo `store-<store_id>.yml`;
- quando o dominio for removido, apague o arquivo correspondente.

## 21. Resolucao da loja pelo `Host`

No backend, voce vai precisar de uma estrategia assim:

1. ler `r.Host`;
2. remover a porta;
3. procurar em `store_domains`;
4. se encontrar, carregar a loja pelo `store_id`;
5. se nao encontrar, cair em fallback por `slug`.

Exemplo de logica:

```go
host := r.Host
host = strings.Split(host, ":")[0]
host = strings.ToLower(host)

store, err := services.Stores.GetByDomain(ctx, host)
if err == nil {
    return store
}

slug := r.URL.Query().Get("slug")
return services.Stores.GetBySlug(ctx, slug)
```

## 22. Fluxo recomendado de balanceamento `dns1` e `dns2`

No inicio:

- `dns1` e `dns2` podem apontar para a mesma VPS.

Depois:

- `dns1` aponta para o proxy do servidor 1;
- `dns2` aponta para o proxy do servidor 2.

No painel, voce distribui clientes:

- metade em `dns1`;
- metade em `dns2`.

Esse balanceamento e manual, mas ja resolve o inicio sem inventar cluster antes da hora.

## 23. Wildcard `*.vrprime.com.br`

Se voce tambem quiser suportar:

- `rainha.vrprime.com.br`
- `maria.vrprime.com.br`

entao use wildcard do seu dominio.

Com Traefik, isso normalmente pede `DNS-01 challenge`.

Nesse caso, voce vai:

1. usar um provedor DNS com API;
2. configurar o resolver ACME do Traefik para DNS challenge;
3. emitir:
   - `vrprime.com.br`
   - `*.vrprime.com.br`

Isso e separado dos dominios externos dos clientes.

## 24. Backup

Voce precisa ter pelo menos estes backups:

- dump diario do PostgreSQL;
- copia do `acme.json`;
- copia da pasta `dynamic` do Traefik;
- copia do `.env`;
- backup do volume do Postgres.

Exemplo de dump:

```bash
docker exec -t $(docker ps -qf name=postgres) pg_dump -U ${POSTGRES_USER} ${POSTGRES_DB} > /srv/goapi/backups/db-$(date +%F).sql
```

Agende com `cron` ou outro scheduler.

## 25. Logs e monitoramento

No minimo acompanhe:

- `docker logs traefik`
- `docker compose -f docker-compose.prod.yml logs api`
- uso de CPU e RAM
- espaco em disco
- expiracao de certificado

Mais adiante, adicione:

- Grafana
- Prometheus
- Loki
- uptime checks externos

## 26. Atualizacao da aplicacao

Fluxo seguro:

```bash
cd /srv/goapi/app
git pull
docker compose -f docker-compose.prod.yml up -d --build
```

Se alterou banco:

- inclua migration nova no repositorio;
- aplique com cuidado antes ou durante o deploy, conforme a mudanca.

## 27. Check-list de subida inicial

Antes do primeiro cliente:

- VPS criada
- Docker instalado
- rede `web` criada
- Traefik rodando
- dominio principal apontando
- certificado do dominio principal funcionando
- app rodando
- banco persistindo dados
- Redis funcionando
- `.env` seguro
- firewall ativo
- backup configurado

Antes de liberar dominio customizado:

- tabela `store_domains` criada
- endpoint de cadastro de dominio criado
- verificador DNS implementado
- escrita de config dinamica do Traefik implementada
- resolucao da loja por `Host` implementada
- fluxo de renovacao e remocao de dominio testado

## 28. O que eu recomendo voce fazer primeiro no codigo

Antes de ir para a VPS, implemente estas partes no projeto:

1. lojas com `slug`;
2. tabela `store_domains`;
3. cadastro de dominio customizado;
4. verificacao DNS;
5. busca da loja por dominio;
6. fallback por `slug`;
7. tela admin para copiar instrucoes DNS;
8. rotina que escreve arquivos dinamicos do Traefik.

## 29. O que ainda nao vale a pena agora

Nao complique o inicio com:

- cluster Swarm
- Kubernetes
- multiplas VPS com failover automatico
- autoscaling
- CDN complexa
- multi-regiao

Comece com:

- 1 VPS boa
- 1 Traefik
- 1 Postgres
- 1 Redis
- 1 app
- 1 rotina de backup

## 30. Resumo tecnico final

A arquitetura recomendada para voce hoje e:

- VPS limpa com Ubuntu
- Docker + Compose
- Traefik na frente
- frontend publico atras do Traefik
- backend privado atras do frontend
- Postgres e Redis internos
- dominio de cliente mapeado no banco
- configuracao dinamica do Traefik por dominio validado
- certificado emitido por dominio

Esse e o caminho mais simples que continua correto quando o produto crescer.

## 31. Referencias tecnicas

- Documentacao do Traefik sobre ACME e TLS:
  - https://doc.traefik.io/traefik/https/overview/
  - https://doc.traefik.io/traefik/v3.3/https/acme/
- Instalacao Docker no Ubuntu:
  - https://docs.docker.com/engine/install/ubuntu/
