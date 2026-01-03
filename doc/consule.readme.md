============================================================
CONSUL SERVICE DISCOVERY GUIDE FOR GO GATEWAY
============================================================

1. ARCHITECTURE OVERVIEW
------------------------------------------------------------
Consul acts as the "Source of Truth" for your microservices.
- Service Registration: Your Go Gateway registers itself on startup.
- Health Checking: Consul pings your /health endpoint.
- KV Store: Holds dynamic configuration (e.g., JWT keys, Rate Limits).
- DNS/API: Other services find the Gateway via Consul instead of static IPs.



2. DOCKER SETUP (Development Mode)
------------------------------------------------------------
Use this command to start a local instance for development:

docker run -d --name consul-dev -p 8500:8500 -p 8600:8600/udp hashicorp/consul agent -dev -client=0.0.0.0 -ui

UI Access: http://localhost:8500


3. CLI CHEAT SHEET (macOS / Docker)
------------------------------------------------------------
# Enter the container shell:
docker exec -it consul-dev /bin/sh

# --- Inside Shell or using 'docker exec -it consul-dev <command>' ---

# Service Discovery
consul catalog services          # List all registered services
consul catalog nodes -service=X  # Show instances of service X
consul health -state=critical    # List failing services

# Key-Value (KV) Store
consul kv put path/key "val"     # Save a config value
consul kv get path/key           # Read a config value
consul kv export "prefix/"       # Export all keys under a prefix

# Membership & Maintenance
consul members                   # List all cluster nodes
consul monitor                   # Stream logs for debugging

5. ENVIRONMENT CONFIGURATION (.env)
------------------------------------------------------------
CONSUL_ADDRESS=127.0.0.1:8500
SERVICE_NAME=gateway-service
SERVICE_CHECK_URL=http://localhost:8080/health
ENV=development

============================================================