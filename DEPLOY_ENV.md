# 🛑 IMPORTANTE: Variáveis de Ambiente para DEPLOY

Para enviar essa aplicação para um servidor em produção (ex: AWS, DigitalOcean), **NÃO** utilize as credenciais padrão do `docker-compose.yml`.

Você deve criar um arquivo `.env` verdadeiro no seu servidor e preenchê-lo com os dados abaixo. **Lembre-se: O arquivo `.env` NUNCA deve ser comitado no Git.**

Copie os dados abaixo, crie um arquivo chamado `.env` no servidor de produção e modifique os valores com senhas reais:

```env
# ==========================================
# CONFIGURAÇÕES DO SISTEMA (PRODUÇÃO)
# ==========================================

# Mude para "production" para desativar logs de debug do Go
APP_ENV=production

# ==========================================
# SEGURANÇA (EXTREMAMENTE IMPORTANTE)
# ==========================================

# Gere uma string muito longa e aleatória (mínimo 64 caracteres)
# Você pode gerar com: `openssl rand -base64 64`
JWT_SECRET=mude_isso_para_uma_chave_enorme_e_impossivel_de_adivinhar

# Tempo de expiração da sessão do usuário
JWT_TTL=24h

# Custo da criptografia de senhas. (Não passe de 14 para não fritar a CPU)
BCRYPT_COST=12

# ==========================================
# BANCO DE DADOS (POSTGRES)
# ==========================================

# Defina a senha do Root do Postgres aqui (Evite caracteres estranhos como @/?)
POSTGRES_USER=app_prod_user
POSTGRES_PASSWORD=senha_muito_forte_aqui_123
POSTGRES_DB=goapi_prod

# A URL que a API do Go vai usar para se conectar.
# DEVE conter exatamente o mesmo usuário e senha definidos acima.
DATABASE_URL=postgres://app_prod_user:senha_muito_forte_aqui_123@postgres:5432/goapi_prod?sslmode=disable

# Controle de Fluxo para escalar. Num cluster, diminua as conexões.
DB_MAX_CONNS=40
DB_MIN_CONNS=5

# ==========================================
# CACHE E SESSÃO (REDIS)
# ==========================================

# A URL do redis interno no Docker
REDIS_URL=redis://redis:6379/0

# ==========================================
# SERVIÇO DE E-MAIL (SMTP)
# ==========================================

# Utilize um serviço como Resend, AWS SES, Mailgun ou SendGrid
SMTP_HOST=smtp.seudominio.com
SMTP_PORT=587
SMTP_USERNAME=api_key_gerada_no_provedor
SMTP_PASSWORD=senha_secreta_do_smtp
SMTP_FROM=nao-responda@suaempresa.com.br
```

---

### ⚠️ DICA DE SEGURANÇA PARA O `docker-compose.yml` NO SERVIDOR:
Lembre-se de remover a exposição das portas do Postgres e do Redis no servidor antes de subir o Docker, para evitar ataques de força bruta no seu banco.

No seu `docker-compose.yml` do servidor, as portas devem ficar apenas no frontend:
```yaml
  frontend:
    ports:
      - "80:80"
```
