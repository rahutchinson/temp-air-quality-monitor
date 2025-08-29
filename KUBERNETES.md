# Air Quality Monitor - Kubernetes Deployment

This guide will help you deploy the Air Quality Monitor application to a local Kubernetes cluster.

## Prerequisites

- Docker installed and running
- At least 4GB RAM available for the cluster
- At least 2 CPU cores available

## Quick Start

### 1. Set up Local Kubernetes Cluster

```bash
# Run the setup script to install and configure minikube
./setup-k8s.sh
```

This script will:
- Install kubectl (if not present)
- Install minikube (if not present)
- Start a local Kubernetes cluster
- Enable necessary addons (ingress, dashboard)

### 2. Deploy the Application

```bash
# Navigate to the k8s directory
cd k8s

# Deploy the application
./deploy.sh
```

This script will:
- Build the Docker image
- Load it into minikube
- Create the namespace and all Kubernetes resources
- Deploy the application

### 3. Access the Application

After deployment, you can access the application at:
```
http://<minikube-ip>:30080
```

To get the minikube IP:
```bash
minikube ip
```

## Manual Deployment Steps

If you prefer to deploy manually, follow these steps:

### 1. Build and Load Docker Image

```bash
# Build the image
docker build -t air-quality-monitor:latest .

# Load into minikube
minikube image load air-quality-monitor:latest
```

### 2. Apply Kubernetes Manifests

```bash
# Create namespace
kubectl apply -f k8s/namespace.yaml

# Apply resources
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/persistent-volume.yaml
kubectl apply -f k8s/persistent-volume-claim.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

### 3. Verify Deployment

```bash
# Check pod status
kubectl get pods -n air-quality-monitor

# Check service
kubectl get services -n air-quality-monitor

# View logs
kubectl logs -f deployment/air-quality-monitor -n air-quality-monitor
```

## Kubernetes Resources

### Namespace
- `air-quality-monitor`: Isolates the application resources

### ConfigMap
- `air-quality-config`: Contains application configuration
  - Device URL
  - Server address
  - Database path

### PersistentVolume & PersistentVolumeClaim
- Provides persistent storage for the SQLite database
- Uses hostPath storage (local filesystem)
- 1GB storage allocation

### Deployment
- Runs 1 replica of the application
- Resource limits: 128Mi memory, 100m CPU
- Health checks (liveness and readiness probes)
- Environment variables from ConfigMap

### Service
- Type: NodePort
- Exposes port 8080 on NodePort 30080
- Allows external access to the application

## Configuration

### Environment Variables

The application can be configured using environment variables:

- `DEVICE_URL`: URL of the PurpleAir device (default: http://192.168.0.249/json)
- `SERVER_ADDR`: Server address (default: :8080)
- `DATABASE_PATH`: Database file path (default: /app/data/air_quality.db)

### Updating Configuration

To update the device URL or other settings:

1. Edit the ConfigMap:
```bash
kubectl edit configmap air-quality-config -n air-quality-monitor
```

2. Restart the deployment:
```bash
kubectl rollout restart deployment/air-quality-monitor -n air-quality-monitor
```

## Monitoring and Logs

### View Application Logs

```bash
# Follow logs in real-time
kubectl logs -f deployment/air-quality-monitor -n air-quality-monitor

# View recent logs
kubectl logs deployment/air-quality-monitor -n air-quality-monitor --tail=100
```

### Check Application Health

```bash
# Check pod status
kubectl get pods -n air-quality-monitor

# Check service endpoints
kubectl get endpoints -n air-quality-monitor

# Test health endpoint
curl http://$(minikube ip):30080/health
```

### Access Kubernetes Dashboard

```bash
# Open the dashboard
minikube dashboard
```

## Scaling

### Scale the Application

```bash
# Scale to 3 replicas
kubectl scale deployment air-quality-monitor --replicas=3 -n air-quality-monitor

# Check scaling status
kubectl get pods -n air-quality-monitor
```

**Note**: Multiple replicas will share the same database file, which may cause conflicts. For production use, consider using a proper database like PostgreSQL.

## Troubleshooting

### Common Issues

1. **Pod not starting**
   ```bash
   # Check pod events
   kubectl describe pod <pod-name> -n air-quality-monitor
   
   # Check pod logs
   kubectl logs <pod-name> -n air-quality-monitor
   ```

2. **Cannot access application**
   ```bash
   # Check service
   kubectl get service air-quality-monitor-service -n air-quality-monitor
   
   # Check if port is accessible
   curl -v http://$(minikube ip):30080/health
   ```

3. **Database issues**
   ```bash
   # Check persistent volume
   kubectl get pv,pvc -n air-quality-monitor
   
   # Check volume mounts
   kubectl describe pod <pod-name> -n air-quality-monitor
   ```

### Reset Everything

```bash
# Clean up application
cd k8s && ./cleanup.sh

# Reset minikube (optional)
minikube delete
minikube start
```

## Production Considerations

For production deployment, consider:

1. **Database**: Use a proper database (PostgreSQL, MySQL) instead of SQLite
2. **Ingress**: Configure proper ingress controller with SSL
3. **Monitoring**: Add Prometheus and Grafana for monitoring
4. **Logging**: Use centralized logging (ELK stack, Fluentd)
5. **Security**: Configure RBAC, network policies, and secrets
6. **Backup**: Set up database backups and disaster recovery
7. **Scaling**: Use Horizontal Pod Autoscaler (HPA)

## Useful Commands

```bash
# Get minikube IP
minikube ip

# Open minikube dashboard
minikube dashboard

# View cluster status
minikube status

# Stop cluster
minikube stop

# Delete cluster
minikube delete

# View all resources
kubectl get all -n air-quality-monitor

# Port forward (alternative to NodePort)
kubectl port-forward service/air-quality-monitor-service 8080:8080 -n air-quality-monitor
```
