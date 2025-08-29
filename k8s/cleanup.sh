#!/bin/bash

# Air Quality Monitor Kubernetes Cleanup Script

echo "=== Cleaning up Air Quality Monitor Kubernetes Deployment ==="
echo

# Delete Kubernetes resources
echo "🗑️  Deleting Kubernetes resources..."
kubectl delete -f service.yaml --ignore-not-found=true
kubectl delete -f deployment.yaml --ignore-not-found=true
kubectl delete -f persistent-volume-claim.yaml --ignore-not-found=true
kubectl delete -f persistent-volume.yaml --ignore-not-found=true
kubectl delete -f configmap.yaml --ignore-not-found=true
kubectl delete -f namespace.yaml --ignore-not-found=true

echo "✅ Cleanup completed successfully!"
echo

# Show remaining resources
echo "📊 Remaining resources in air-quality-monitor namespace:"
kubectl get all -n air-quality-monitor 2>/dev/null || echo "   Namespace air-quality-monitor not found"

echo
echo "🎉 Cleanup completed!"
