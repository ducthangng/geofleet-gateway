# declare the providers for the complete infrastructure
# each package will have a provider, install it before proceeding
terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0.1"
    }

    kafka = {
      source  = "Mongey/kafka"
      version = "~> 0.7"
    }

    postgresql = {
      source = "cyrilgdn/postgresql"
      version = "1.26.0"
    }

    time = {
      source = "hashicorp/time"
      version = "0.13.1"
    }
  }
}

provider "docker" {}

provider "kafka" {
  bootstrap_servers = ["127.0.0.1:9092"]
}

# --- NETWORKS ---
resource "docker_network" "gateway_network" {
  name = "gateway-network"
}

# --- VOLUMES ---
resource "docker_volume" "postgres_data" {
  name = "postgres_data"
}

# --- REDIS ---
resource "docker_container" "redis" {
  name  = "redis-dev"
  image = "redis:7-alpine"
  restart = "always"
  command = ["redis-server", "--save", "60", "1", "--loglevel", "warning"]
  
  ports {
    internal = 6379
    external = 6379
  }
}

# --- KAFKA (KRaft mode) ---
resource "docker_container" "kafka" {
  name  = "kafka-dev"
  image = "apache/kafka:3.8.0"
  
  ports {
    internal = 9092
    external = 9092
  }

  env = [
    "KAFKA_PROCESS_ROLES=broker,controller",
    "KAFKA_NODE_ID=1",
    "KAFKA_CONTROLLER_QUORUM_VOTERS=1@localhost:9093", # Dùng localhost
    "KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER",
    "KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093",
    "KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://127.0.0.1:9092",
    "KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT",
    "KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT",
    "KAFKA_CLUSTER_ID=6bc8e63a-44c1-4d10-8b40-97500096677a",
    
    # --- CÁC DÒNG QUAN TRỌNG CHO SINGLE NODE ---
    "KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1",
    "KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR=1",
    "KAFKA_TRANSACTION_STATE_LOG_MIN_ISR=1",
    "KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS=0"
  ]

  # Bảo vệ Kafka khỏi việc vô tình bị xóa
  # lifecycle {
  #   prevent_destroy = true
  # }
}

# initialization
resource "kafka_topic" "raw_coordinates" {
  # Provider của Monzo sẽ gửi lệnh đến localhost:9092 để tạo cái này
  name = "tracking.raw-coordinates"
  replication_factor = 1 
  partitions = 6

  config = {
    "retention.ms" = "86400000"
  }

  depends_on         = [time_sleep.wait_for_containers]
}

# --- CONSUL ---
resource "docker_container" "consul" {
  name  = "consul-server"
  image = "consul:1.15.4"
  restart = "always"
  command = ["agent", "-dev", "-client=0.0.0.0", "-ui"]
  
  networks_advanced {
    name = docker_network.gateway_network.name
  }

  ports {
    internal = 8500
    external = 8500
  }

  ports {
    internal = 8600
    external = 8600
    protocol = "udp"
  }
}

# --- POSTGIS DB ---
resource "docker_container" "db" {
  name  = "postgis_db"
  image = "kartoza/postgis:18-3.6--v2025.11.24"
  restart = "always"
  
  env = [
    "POSTGRES_USER=root",
    "POSTGRES_PASS=secret",
    "POSTGRES_DBNAME=geofleet"
  ]

  ports {
    internal = 5432
    external = 5432
  }

  mounts {
    target = "/var/lib/postgresql/data"
    source = docker_volume.postgres_data.name
    type   = "volume"
  }

  # lifecycle {
  #   prevent_destroy = true
  # }
}


# Provider default
provider "postgresql" {
  host            = "localhost"
  port            = 5432
  database        = "postgres"
  username        = "root"
  password        = "secret"
  sslmode         = "disable"
}

# wait for 20 second before doing anything
resource "time_sleep" "wait_for_containers" {
  depends_on = [docker_container.db, docker_container.kafka]
  create_duration = "8s"
}

# # Provider cho Ride Tracking Service
# provider "postgresql" {
#   alias           = "ride_tracking"
#   host            = "localhost"
#   port            = 5432
#   database        = "geofleet_ride_tracking_service"
#   username        = "root"
#   password        = "secret"
#   sslmode         = "disable"
# }

resource "postgresql_database" "user_db" {
  name       = "geofleet_user_service"
  depends_on = [time_sleep.wait_for_containers, docker_container.db]
}

resource "postgresql_database" "ride_db" {
  name       = "geofleet_ride_tracking_service"
  depends_on = [time_sleep.wait_for_containers, docker_container.db]
}

provider "postgresql" {
  alias           = "user_service"
  host            = "localhost"
  database        = postgresql_database.user_db.name
  username        = "root"
  password        = "secret"
}

provider "postgresql" {
  alias           = "ride_tracking"
  host            = "localhost"
  database        = postgresql_database.ride_db.name
  username        = "root"
  password        = "secret"
}

# Tạo schema cho User Service
resource "postgresql_schema" "user_schema" {
  provider = postgresql.user_service # Trỏ đúng vào alias ở trên
  name     = "user_domain"
  depends_on = [postgresql_database.user_db]
}

# Tạo schema cho Ride Service
resource "postgresql_schema" "ride_tracking_schema" {
  provider = postgresql.ride_tracking # Trỏ đúng vào alias ở trên
  name     = "ride_tracking_domain"
  depends_on = [postgresql_database.ride_db]
}

# --- SWAGGER UI ---
resource "docker_container" "swagger_ui" {
  name  = "geofleet-swagger-ui"
  image = "swaggerapi/swagger-ui"
  
  ports {
    internal = 8080
    external = 8081
  }

  env = [
    "SWAGGER_JSON=/app/geofleet.swagger.json"
  ]

  host {
    host = "host.docker.internal"
    ip   = "host-gateway"
  }

  # Lưu ý: Terraform cần đường dẫn tuyệt đối cho volumes local
  # Hãy thay đổi path bên dưới cho đúng với máy của bạn
  volumes {
    host_path      = "${abspath(path.module)}/geofleet-proto/gen/openapiv2/geofleet.swagger.json"
    container_path = "/app/geofleet.swagger.json"
  }
}


