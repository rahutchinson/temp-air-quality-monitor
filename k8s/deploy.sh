#!/bin/bash

# Air Quality Monitor Kubernetes Deployment Script

set -e

echo "=== Air Quality Monitor Kubernetes Deployment ==="
echo

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is not installed. Please install kubectl first."
    exit 1
fi

# Check if minikube is running
if ! minikube status | grep -q "Running"; then
    echo "❌ Minikube is not running. Please start minikube first:"
    echo "   minikube start"
    exit 1
fi

echo "✅ Prerequisites check passed"
echo

# Build Docker image
echo "🔨 Building Docker image..."
docker build -t air-quality-monitor:latest .

# Load image into minikube
echo "📦 Loading image into minikube..."
minikube image load air-quality-monitor:latest

echo "✅ Image built and loaded successfully"
echo

# Create namespace
echo "📁 Creating namespace..."
kubectl apply -f namespace.yaml

# Apply Kubernetes manifests
echo "🚀 Deploying application..."
kubectl apply -f configmap.yaml
kubectl apply -f persistent-volume.yaml
kubectl apply -f persistent-volume-claim.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

echo "✅ Application deployed successfully"
echo

# Wait for deployment to be ready
echo "⏳ Waiting for deployment to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/air-quality-monitor -n air-quality-monitor

echo "✅ Deployment is ready!"
echo

# Get service information
echo "🌐 Service Information:"
echo "   NodePort: $(kubectl get service air-quality-monitor-service -n air-quality-monitor -o jsonpath='{.spec.ports[0].nodePort}')"
echo "   Access URL: http://$(minikube ip):$(kubectl get service air-quality-monitor-service -n air-quality-monitor -o jsonpath='{.spec.ports[0].nodePort}')"
echo

# Show pod status
echo "📊 Pod Status:"
kubectl get pods -n air-quality-monitor

echo
echo "🎉 Deployment completed successfully!"
echo
echo "To access the application:"
echo "   http://$(minikube ip):$(kubectl get service air-quality-monitor-service -n air-quality-monitor -o jsonpath='{.spec.ports[0].nodePort}')"
echo
echo "To view logs:"
echo "   kubectl logs -f deployment/air-quality-monitor -n air-quality-monitor"
echo
echo "To delete the deployment:"
echo "   ./k8s/cleanup.sh"
