# Database Configuration
export DB_HOST="${DB_HOST:-localhost}"
export DB_PORT="${DB_PORT:-5433}"
export DB_USERNAME="${DB_USERNAME:-my_root}"
export DB_PASSWORD="${DB_PASSWORD:-987654321}"
export DB_NAME="${DB_NAME:-my_DB}"
export DB_SCHEMA="${SCHEMA:-public}"
export SSLMODE="${SSLMODE:-disable}"

# Connection Pool Configuration
export MAX_OPEN_CONNS="${MAX_OPEN_CONNS:-25}"
export MAX_IDLE_CONNS="${MAX_IDLE_CONNS:-5}"
export MAX_LIFETIME="${MAX_LIFETIME:-5m}"
export MAX_IDLE_TIME="${MAX_IDLE_TIME:-10m}"