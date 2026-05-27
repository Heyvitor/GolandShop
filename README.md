# 🚀 GoApi - High-Performance Secure System

Este repositório contém uma aplicação Full-Stack de nível *Enterprise*, arquitetada para suportar alta escalabilidade (High Throughput) e garantir segurança de nível bancário contra ataques web comuns. 

O sistema é composto por uma **API em Go (Golang)** ultra-rápida, e um **Frontend em React (Vite + TypeScript)** leve, ambos orquestrados via **Docker** e servidos por trás de um Reverse Proxy **Nginx**.

---

## 🏗️ Arquitetura e Tecnologias

### ⚙️ Backend (Go)
*   **Linguagem:** Go 1.22+
*   **Roteamento:** Padrão (`net/http` ServeMux - sem frameworks pesados).
*   **Banco de Dados:** PostgreSQL com `pgx/v5` (Driver de altíssima performance) usando Pool de conexões.
*   **Cache e Proteção:** Redis (`go-redis/v9`).
*   **Segurança:** 
    *   Sessões JWT validadas criptograficamente, anexadas em UUIDs (JTI).
    *   Compressão *Gzip* via Middleware para economia de banda em listas JSON grandes.
*   **Notificações:** Disparo de e-mails em *background* via SMTP nativo.

### 🎨 Frontend (React)
*   **Ferramental:** Vite (Rápido HMR e Build).
*   **Linguagem:** React com TypeScript.
*   **Roteamento:** `wouter` (Leve, sem o peso do *react-router*).
*   **Estilização:** CSS Vanilla puro com variáveis e semântica limpa.

### 🐳 Infraestrutura (Docker)
Todo o ecossistema roda de forma unificada usando `docker compose`, permitindo o desenvolvimento e o deploy integrados.

---

## 🛡️ Camadas de Segurança Implementadas

1. **Cookies HttpOnly:** O Token JWT não é entregue ao React, mitigando totalmente o risco de ataques XSS que tentam roubar tokens via LocalStorage. O navegador faz o envio via cookie com regras rígidas de segurança (`SameSite=Strict` e `Secure`).
2. **Reverse Proxy Nginx:** Elimina erros de CORS fundindo o Backend e o Frontend na porta 80. O frontend aponta para a mesma raiz (`/`) e o Nginx faz o bypass das chamadas `/api/v1` escondendo o container do Go.
3. **Limitação de Payload (OOM Protection):** A leitura de requisições JSON possui um limitador (`http.MaxBytesReader`) restrito a 1MB para evitar Ataques de Esgotamento de Memória / DoS em envios massivos.
4. **Rate Limiting no Redis:** Rotas de *CPU-bound* pesadas usando Bcrypt (Login e Register) são protegidas através de um *Sliding Window Rate Limiter* via Redis.
5. **JWT Blacklist (Logout):** Como JWTs são *stateless*, implementamos uma rota segura de Logout. O Cookie no cliente é destruído e o ID do token ganha um registro no Redis (expirando no seu limite natural) impedindo a sua reutilização.

---

## 🚀 Como Rodar o Projeto

É incrivelmente fácil rodar tudo em qualquer máquina com apenas o **Docker** instalado.

### 1. Inicie a stack
Clone o projeto, abra um terminal na raiz (`/`) e execute:

```bash
docker compose up --build
```

Isso fará o *build* limpo do Go, o *build* final do React/Nginx e subirá todas as instâncias (Postgres, Redis, API, Frontend).

### 2. Acesse a Aplicação
*   **Frontend (Site):** [http://localhost](http://localhost)
*   **Backend (API Interna via Proxy):** [http://localhost/api/v1/healthz](http://localhost/api/v1/healthz)

### 3. Serviços em Backstage (Mapeamento de Portas Locais)
Você pode utilizar clientes visuais (DBeaver, RedisInsight, etc) na sua própria máquina (Host) para interagir com o estado atualizado:
*   **Postgres:** `localhost:5432` *(User: goapi / Pass: goapi)*
*   **Redis:** `localhost:6379`

### 4. Configurar E-mails (SMTP)
No `docker-compose.yml`, o envio de e-mails já está preparado e roda assincronamente (background) sempre que um usuário se cadastra. Recomendamos usar o **Mailtrap** para testes locais, trocando as variáveis:
*   `SMTP_USERNAME`
*   `SMTP_PASSWORD`

---

## 🧱 Como Desenvolver o Visual (Frontend Dev)
Para que as páginas do site não precisem ser geradas em um *build do Docker* para cada pequena alteração visual durante o dia-a-dia, você pode iniciar os serviços essenciais no terminal A:
```bash
docker compose up postgres redis api
```

E em outro terminal, inicializar o Vite:
```bash
cd frontend
npm run dev
```
*(Atenção: Ao rodar isolado assim para fins visuais, as chamadas na porta 5173 para a porta 8080 poderão gerar bloqueios de CORS devido às restrições restritas dos cookies HttpOnly)*.

---
*Feito com foco extremo em Performance e Simplicidade.*
