# Use the official HashiCorp Consul image
FROM hashicorp/consul:1.15

# Optional: Add a custom configuration directory
# This allows you to pre-define services or checks in JSON
COPY ./consul-config /consul/config

# Optional: Set the working directory
WORKDIR /consul

# Expose necessary ports:
# 8500: HTTP API & UI
# 8600: DNS Server
EXPOSE 8500 8600/udp

# Run consul in agent/development mode by default
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["agent", "-dev", "-client", "0.0.0.0"]