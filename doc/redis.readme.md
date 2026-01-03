============================================================
REDIS CACHING & RATE LIMITING GUIDE FOR GO GATEWAY
============================================================

1. ARCHITECTURE OVERVIEW
------------------------------------------------------------
In your Gateway, Redis serves three primary roles:
- Rate Limiting: Tracking request counts per IP.
- Session Storage: Storing JWT blacklists or temporary user sessions.
- Distributed Locking: Ensuring single-execution for background tasks.

2. QUICK START (DOCKER)
------------------------------------------------------------
# Start Redis
docker run -d --name redis-dev -p 6379:6379 redis:alpine

# Check if it's alive
docker exec -it redis-dev redis-cli ping
# Response should be: PONG

3. CLI CHEAT SHEET (MacOS / Docker)
------------------------------------------------------------
# Enter the Redis CLI:
docker exec -it redis-dev redis-cli

# --- Basic Key Operations ---
SET user:1 "John Doe"           # Set a key
GET user:1                      # Get a key
DEL user:1                      # Delete a key
EXISTS user:1                   # Check if key exists (1 = yes, 0 = no)

# --- Expiration & TTL ---
EXPIRE user:1 60                # Set key to expire in 60 seconds
TTL user:1                      # See remaining time (in seconds)

# --- Debugging & Monitoring ---
MONITOR                         # Stream all commands (Use for debugging)
FLUSHALL                        # Delete EVERYTHING (Use with caution!)
INFO replication                # Check Master/Slave status

4. GO USAGE SNIPPET
------------------------------------------------------------
ctx := context.Background()
rdb, _ := singleton.GetRedisClient(config)

// Setting a key with 1 hour expiration
err := rdb.Set(ctx, "key", "value", time.Hour).Err()

// Getting a key
val, err := rdb.Get(ctx, "key").Result()

5. ENVIRONMENT VARIABLES (.env)
------------------------------------------------------------
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_USER=default
REDIS_PASS=
REDIS_DB=0
============================================================