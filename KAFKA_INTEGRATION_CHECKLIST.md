# Kafka Real-Time Integration Checklist

## Pre-Setup Verification
- [ ] Docker and Docker Compose installed
- [ ] Go 1.23+ installed
- [ ] At least 4GB RAM available
- [ ] Ports 9092, 8080, 6379, 5432 free

## Phase 1: Environment Setup
- [ ] Clone/navigate to go-pro repository
- [ ] Review `KAFKA_IMPLEMENTATION_SUMMARY.md`
- [ ] Review `KAFKA_REALTIME_GUIDE.md`
- [ ] Verify all new files present:
  - [ ] `backend/internal/messaging/kafka/analytics.go`
  - [ ] `backend/internal/messaging/kafka/notifications.go`
  - [ ] `backend/internal/messaging/kafka/dlq.go`
  - [ ] `backend/internal/messaging/realtime/websocket.go`
  - [ ] `backend/internal/messaging/kafka/example_integration.go`
  - [ ] `docker-compose.kafka.yml`
  - [ ] `KAFKA_QUICK_START.sh`

## Phase 2: Docker Setup
- [ ] Run: `./KAFKA_QUICK_START.sh start`
- [ ] Verify all services running: `./KAFKA_QUICK_START.sh status`
- [ ] Check service health: `./KAFKA_QUICK_START.sh health`
- [ ] Access Kafka UI: `http://localhost:8080`
- [ ] Create topics: `./KAFKA_QUICK_START.sh create-topics`
- [ ] Verify topics: `./KAFKA_QUICK_START.sh list-topics`

## Phase 3: Backend Configuration
- [ ] Update `backend/.env`:
  ```env
  KAFKA_BROKERS=localhost:9092
  KAFKA_GROUP_ID=go-pro-consumer-group
  MESSAGING_ENABLED=true
  ```
- [ ] Add gorilla/websocket: `go get github.com/gorilla/websocket`
- [ ] Run migrations if needed
- [ ] Build backend: `go build ./cmd/server`

## Phase 4: Integration Testing
- [ ] Start backend: `go run ./cmd/server`
- [ ] Check logs for Kafka connection success
- [ ] Produce test event: `./KAFKA_QUICK_START.sh produce`
- [ ] Verify in Kafka UI
- [ ] Check consumer lag: `./KAFKA_QUICK_START.sh lag`

## Phase 5: Feature Validation

### Analytics
- [ ] Check `analyticsConsumer` starts without errors
- [ ] Monitor metrics in real-time
- [ ] Verify aggregations are accurate
- [ ] Test course metrics calculation

### Notifications
- [ ] Complete a lesson (triggers progress event)
- [ ] Verify notification handler is called
- [ ] Check notification output (print/log)
- [ ] Test different score thresholds

### DLQ
- [ ] Deliberately fail a handler
- [ ] Verify message goes to DLQ topic
- [ ] Check DLQ message in Kafka UI
- [ ] Verify retry logic triggers
- [ ] Monitor retry exponential backoff

### WebSocket
- [ ] Implement simple test client
- [ ] Connect to WebSocket endpoint (when added)
- [ ] Receive analytics updates
- [ ] Receive notifications
- [ ] Verify disconnection handling

## Phase 6: Monitoring Setup
- [ ] Access Prometheus: `http://localhost:9090`
- [ ] Add Kafka broker targets
- [ ] Access Grafana: `http://localhost:3000`
- [ ] Import Kafka dashboards
- [ ] Create custom application dashboards
- [ ] Set up alert rules for lag threshold

## Phase 7: Documentation & Handoff
- [ ] Update team docs with:
  - [ ] Event publishing guide
  - [ ] Notification handler registration
  - [ ] WebSocket endpoint documentation
  - [ ] DLQ monitoring instructions
- [ ] Create runbook for common operations:
  - [ ] Restart Kafka
  - [ ] Reset consumer group
  - [ ] Replay messages from DLQ
  - [ ] Monitor consumer lag
- [ ] Document troubleshooting steps

## Phase 8: Production Readiness
- [ ] Increase replication factor for HA
- [ ] Configure persistent volumes
- [ ] Set up log aggregation (ELK/Loki)
- [ ] Configure resource limits
- [ ] Set up alerting for:
  - [ ] Broker down
  - [ ] High consumer lag (>10s)
  - [ ] DLQ messages accumulating
  - [ ] Producer failures
- [ ] Plan for backup/recovery
- [ ] Load test with expected throughput

## Operation Checklist

### Daily Checks
- [ ] Kafka cluster healthy: `./KAFKA_QUICK_START.sh health`
- [ ] Consumer lag <5s: `./KAFKA_QUICK_START.sh lag`
- [ ] No messages in DLQ or investigated
- [ ] Grafana dashboards green

### Weekly Checks
- [ ] Review DLQ for patterns
- [ ] Verify retention policies are working
- [ ] Check disk usage on Kafka
- [ ] Review logs for warnings/errors

### Monthly Tasks
- [ ] Update Prometheus/Grafana configs
- [ ] Review and optimize batch sizes
- [ ] Analyze throughput trends
- [ ] Plan capacity for growth

## Common Commands Reference

**Start/Stop:**
```bash
./KAFKA_QUICK_START.sh start          # Start stack
./KAFKA_QUICK_START.sh stop           # Stop stack
./KAFKA_QUICK_START.sh restart        # Restart
./KAFKA_QUICK_START.sh rebuild        # Clean rebuild
```

**Monitoring:**
```bash
./KAFKA_QUICK_START.sh status         # See running services
./KAFKA_QUICK_START.sh health         # Health check
./KAFKA_QUICK_START.sh logs kafka     # View broker logs
./KAFKA_QUICK_START.sh lag             # Check consumer lag
```

**Topics:**
```bash
./KAFKA_QUICK_START.sh list-topics    # List all topics
./KAFKA_QUICK_START.sh describe <topic> # Topic details
./KAFKA_QUICK_START.sh create-topics  # Create all topics
./KAFKA_QUICK_START.sh clean          # Delete all topics
```

**Testing:**
```bash
./KAFKA_QUICK_START.sh produce        # Send test event
./KAFKA_QUICK_START.sh consume <topic> # Read events
./KAFKA_QUICK_START.sh run-backend    # Start backend
```

**URLs:**
```bash
./KAFKA_QUICK_START.sh urls           # Show all URLs
```

## Service Ports
| Service | Port | URL |
|---------|------|-----|
| Kafka | 9092 | `localhost:9092` |
| Kafka UI | 8080 | `http://localhost:8080` |
| Prometheus | 9090 | `http://localhost:9090` |
| Grafana | 3000 | `http://localhost:3000` |
| Adminer | 8081 | `http://localhost:8081` |
| Qdrant | 6333 | `http://localhost:6333` |
| PostgreSQL | 5432 | `localhost:5432` |
| Redis | 6379 | `localhost:6379` |

## Troubleshooting Quick Reference

**Kafka won't start:**
```bash
# Check logs
./KAFKA_QUICK_START.sh logs kafka

# Verify Zookeeper is running
./KAFKA_QUICK_START.sh logs zookeeper

# Rebuild from scratch
./KAFKA_QUICK_START.sh rebuild
```

**No events in topic:**
```bash
# Check if messages are being produced
./KAFKA_QUICK_START.sh produce

# Verify in Kafka UI
# Go to http://localhost:8080 → Topics → go-pro.progress.events

# Check producer logs in backend
grep -i "published\|error" backend/logs.txt
```

**Consumer lagging:**
```bash
# Check lag
./KAFKA_QUICK_START.sh lag

# Reset to earliest
./KAFKA_QUICK_START.sh reset go-pro-consumer-group

# Check handler performance
grep -i "processing\|error" backend/logs.txt
```

**DLQ accumulating messages:**
```bash
# View DLQ topic
./KAFKA_QUICK_START.sh consume go-pro.dlq

# Check errors in logs
./KAFKA_QUICK_START.sh logs backend

# Once fixed, replay:
./KAFKA_QUICK_START.sh reset go-pro-dlq-consumer
```

## Performance Baseline

After setup, verify these metrics:

**Producer Performance:**
- [ ] ~10,000 events/sec single producer
- [ ] Message size: ~500 bytes average
- [ ] Latency: <50ms p99 from publish to broker acknowledge

**Consumer Performance:**
- [ ] Processing latency: <100ms p99
- [ ] Throughput: matches producer rate
- [ ] Lag: increasing then catching up (not stuck)

**Analytics Aggregation:**
- [ ] Metrics update: <500ms from event
- [ ] Active users: accurate count
- [ ] Scores: proper averaging

## Sign-Off

- [ ] All checks passed
- [ ] Team trained on operations
- [ ] Documentation complete
- [ ] Ready for production deployment

**Date Completed**: ___________
**Verified By**: ___________
**Notes**:

---

## Next Steps After Integration

1. **Connect Frontend**:
   - Add WebSocket endpoint to backend
   - Create dashboard component
   - Subscribe to real-time updates

2. **Add More Consumers**:
   - Real-time search indexing
   - ML pipeline triggers
   - Email/SMS notifications

3. **Scale Out**:
   - Multiple Kafka brokers
   - Topic partitioning
   - Consumer group scaling

4. **Advanced Features**:
   - Exactly-once semantics
   - Stream processing (Flink/Kafka Streams)
   - Schema registry integration

## Resources
- **Summary**: `KAFKA_IMPLEMENTATION_SUMMARY.md`
- **Guide**: `KAFKA_REALTIME_GUIDE.md`
- **Quick Start**: `./KAFKA_QUICK_START.sh`
- **Source Code**: `backend/internal/messaging/`
