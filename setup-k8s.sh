#!/bin/bash

# Local Kubernetes Cluster Setup Script

set -e

echo "=== Setting up Local Kubernetes Cluster ==="
echo

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    echo "   Visit: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

echo "✅ Docker is running"
echo

# Install kubectl if not present
if ! command -v kubectl &> /dev/null; then
    echo "📦 Installing kubectl..."
    
    # Detect OS
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
        chmod +x kubectl
        sudo mv kubectl /usr/local/bin/
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/amd64/kubectl"
        chmod +x kubectl
        sudo mv kubectl /usr/local/bin/
    else
        echo "❌ Unsupported OS. Please install kubectl manually."
        echo "   Visit: https://kubernetes.io/docs/tasks/tools/install-kubectl/"
        exit 1
    fi
    
    echo "✅ kubectl installed successfully"
else
    echo "✅ kubectl is already installed"
fi
echo

# Install minikube if not present
if ! command -v minikube &> /dev/null; then
    echo "📦 Installing minikube..."
    
    # Detect OS
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
        sudo install minikube-linux-amd64 /usr/local/bin/minikube
        rm minikube-linux-amd64
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-amd64
        sudo install minikube-darwin-amd64 /usr/local/bin/minikube
        rm minikube-darwin-amd64
    else
        echo "❌ Unsupported OS. Please install minikube manually."
        echo "   Visit: https://minikube.sigs.k8s.io/docs/start/"
        exit 1
    fi
    
    echo "✅ minikube installed successfully"
else
    echo "✅ minikube is already installed"
fi
echo

# Start minikube cluster
echo "🚀 Starting minikube cluster..."
if minikube status | grep -q "Running"; then
    echo "✅ Minikube cluster is already running"
else
    minikube start --driver=docker --memory=4096 --cpus=2
    echo "✅ Minikube cluster started successfully"
fi
echo

# Enable addons
echo "🔧 Enabling minikube addons..."
minikube addons enable ingress
minikube addons enable dashboard
echo "✅ Addons enabled successfully"
echo

# Show cluster status
echo "📊 Cluster Status:"
minikube status
echo

# Show cluster info
echo "🌐 Cluster Information:"
kubectl cluster-info
echo

# Show nodes
echo "🖥️  Nodes:"
kubectl get nodes
echo

echo "🎉 Local Kubernetes cluster setup completed!"
echo
echo "Next steps:"
echo "1. Deploy the air quality monitor:"
echo "   cd k8s && ./deploy.sh"
echo
echo "2. Access the application:"
echo "   http://$(minikube ip):30080"
echo
echo "3. Access the Kubernetes dashboard:"
echo "   minikube dashboard"
echo
echo "4. View cluster logs:"
echo "   minikube logs"
echo
echo "5. Stop the cluster:"
echo "   minikube stop"
echo
echo "6. Delete the cluster:"
echo "   minikube delete"
