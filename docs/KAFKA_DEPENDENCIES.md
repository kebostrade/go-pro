# Kafka Real-Time Implementation - Dependencies

## Go Module Dependencies

Add these to your `backend/go.mod` file:

### Required (Already Installed)
```
github.com/segmentio/kafka-go v0.4.x
```

### New Requirement
```bash
cd backend
go get github.com/gorilla/websocket@latest
```

### Complete go.mod Entry
```
require (
    github.com/gorilla/websocket v1.5.x
    github.com/segmentio/kafka-go v0.4.x
    // ... other existing dependencies
)
```

## Docker Images Used

The `docker-compose.kafka.yml` uses these images (pre-configured):

```
Kafka:                  confluentinc/cp-kafka:7.5.0
Zookeeper:             confluentinc/cp-zookeeper:7.5.0
Kafka UI:              provectuslabs/kafka-ui:latest
Redis:                 redis:7-alpine
PostgreSQL:            postgres:15-alpine
Prometheus:            prom/prometheus:latest
Grafana:               grafana/grafana:latest
Qdrant:                qdrant/qdrant:latest
Adminer:               adminer:latest
```

## System Requirements

### Minimum (Development)
- **CPU**: 2 cores
- **Memory**: 4GB RAM
- **Disk**: 5GB free space
- **Docker**: 20.10+
- **Docker Compose**: 1.29+

### Recommended (Development)
- **CPU**: 4+ cores
- **Memory**: 8GB RAM
- **Disk**: 10GB free space
- **Docker**: Latest stable
- **Docker Compose**: Latest stable

### Production (HA Setup)
- **CPU**: 8+ cores (per broker)
- **Memory**: 16GB+ RAM
- **Disk**: 100GB+ (depends on retention)
- **Network**: 1Gbps+
- **Replication Factor**: 3+
- **Min In-Sync Replicas**: 2

## Installation Steps

### 1. Update Go Dependencies
```bash
cd backend
go get github.com/gorilla/websocket@latest
go mod tidy
```

### 2. Verify Installation
```bash
# Check that imports work
grep "gorilla/websocket" backend/go.mod
grep "segmentio/kafka-go" backend/go.mod

# Should show:
# require (
#     github.com/gorilla/websocket v1.5.x
#     github.com/segmentio/kafka-go v0.4.x
#     ...
# )
```

### 3. Docker Verification
```bash
# Ensure Docker is running
docker --version
docker-compose --version

# Pull required images (optional, docker-compose will do this)
docker-compose -f docker-compose.kafka.yml pull
```

## Verification Checklist

- [ ] Go 1.23+ installed: `go version`
- [ ] Gorilla WebSocket installed: `grep gorilla backend/go.mod`
- [ ] Kafka-go installed: `grep kafka-go backend/go.mod`
- [ ] Docker installed: `docker --version`
- [ ] Docker Compose installed: `docker-compose --version`
- [ ] Sufficient disk space: `df -h`
- [ ] Sufficient memory: `free -h`

## Common Installation Issues

### Issue: "go get" fails for gorilla/websocket
**Solution**:
```bash
cd backend
go get -u github.com/gorilla/websocket
go mod tidy
```

### Issue: Docker images won't pull
**Solution**:
```bash
# Check internet connectivity
ping hub.docker.com

# Try manually pulling an image
docker pull confluentinc/cp-kafka:7.5.0

# If behind proxy, configure Docker:
# Edit ~/.docker/config.json
```

### Issue: Insufficient disk space for Docker
**Solution**:
```bash
# Check available space
df -h /var/lib/docker

# Clean up unused Docker resources
docker system prune -a
docker volume prune
```

### Issue: Port already in use
**Solution**:
```bash
# Find process using port 9092
lsof -i :9092

# Or for docker-compose, modify docker-compose.kafka.yml:
# Change ports to available alternatives
# e.g., 9093:9092 instead of 9092:9092
```

## Version Compatibility Matrix

| Component | Version | Status |
|-----------|---------|--------|
| Go | 1.23+ | ✅ Supported |
| Kafka | 7.5.0+ | ✅ Supported |
| Zookeeper | 7.5.0+ | ✅ Supported |
| gorilla/websocket | 1.5+ | ✅ Supported |
| segmentio/kafka-go | 0.4+ | ✅ Supported |
| Docker | 20.10+ | ✅ Supported |
| Docker Compose | 1.29+ | ✅ Supported |

## Upgrading Dependencies

To upgrade to latest versions:

```bash
cd backend

# Check current versions
go list -m all | grep -E "websocket|kafka"

# Update to latest
go get -u github.com/gorilla/websocket
go get -u github.com/segmentio/kafka-go

# Verify compatibility
go mod tidy
go test ./...
```

## Security Considerations

### For Development
- Default credentials are fine
- Plaintext SASL is acceptable
- Single broker is okay

### For Production
- Enable SASL/SSL authentication
- Use multi-broker setup
- Configure network policies
- Regular security updates
- Monitor access logs

```yaml
# Production Kafka config addition
environment:
  KAFKA_SECURITY_PROTOCOL: SASL_SSL
  KAFKA_SASL_MECHANISM: PLAIN
  KAFKA_SASL_USERNAME: admin
  KAFKA_SASL_PASSWORD: ${KAFKA_PASSWORD}
```

## Performance Tuning

### Connection Pooling
```go
// Kafka producer pooling (already optimized in implementation)
producer := kafka.NewProducer(config, topics)
// Reuse single producer instance for multiple operations
```

### Memory Usage
- Kafka broker: ~2GB base + 4GB buffer
- Prometheus: ~1GB
- Grafana: ~512MB
- Total stack: ~8-10GB recommended

## Backup & Recovery

### Backup Kafka Data
```bash
# Docker volumes contain Kafka data
docker-compose -f docker-compose.kafka.yml exec kafka \
  cp -r /var/lib/kafka-logs /backup/

# Or backup volume directly
docker run --rm -v go-pro-kafka_data:/data \
  -v /backup:/backup \
  busybox tar czf /backup/kafka-backup.tar.gz /data
```

### Restore from Backup
```bash
docker-compose -f docker-compose.kafka.yml down -v
docker-compose -f docker-compose.kafka.yml up -d
# Restore volume files before starting consumer
```

## Support & Troubleshooting

For installation issues:
1. Check `./KAFKA_QUICK_START.sh health`
2. Review `KAFKA_REALTIME_GUIDE.md` troubleshooting section
3. Check Docker logs: `docker-compose -f docker-compose.kafka.yml logs`
4. Check application logs: `go run ./cmd/server 2>&1 | grep -i error`

## Next Steps

Once dependencies are verified:
1. Start Kafka stack: `./KAFKA_QUICK_START.sh start`
2. Build backend: `cd backend && go build ./cmd/server`
3. Run backend: `go run ./cmd/server`
4. Monitor: `./KAFKA_QUICK_START.sh urls`
