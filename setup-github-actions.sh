#!/bin/bash

# GitHub Actions Setup Script for Air Quality Monitor

echo "=== GitHub Actions Setup for Air Quality Monitor ==="
echo

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Not in a git repository. Please run this script from your project root."
    exit 1
fi

# Check if we have a remote origin
if ! git remote get-url origin > /dev/null 2>&1; then
    echo "❌ No remote origin found. Please add your GitHub repository as origin:"
    echo "   git remote add origin https://github.com/username/repository.git"
    exit 1
fi

echo "✅ Git repository detected"
echo

# Get repository information
REPO_URL=$(git remote get-url origin)
REPO_NAME=$(basename -s .git "$REPO_URL")
echo "Repository: $REPO_NAME"
echo

echo "📋 GitHub Actions Setup Instructions:"
echo "====================================="
echo

echo "1. Enable GitHub Actions:"
echo "   - Go to your repository on GitHub"
echo "   - Click on 'Actions' tab"
echo "   - Click 'Enable Actions' if prompted"
echo

echo "2. Configure Repository Secrets (if deploying to remote clusters):"
echo "   - Go to Settings > Secrets and variables > Actions"
echo "   - Add the following secrets:"
echo

echo "   For Development/Testing:"
echo "   - KUBE_CONFIG: Base64 encoded kubeconfig for your cluster"
echo "     (Generate with: kubectl config view --raw | base64 -w 0)"
echo

echo "   For Staging:"
echo "   - STAGING_KUBE_CONFIG: Base64 encoded kubeconfig for staging cluster"
echo

echo "   For Production:"
echo "   - PRODUCTION_KUBE_CONFIG: Base64 encoded kubeconfig for production cluster"
echo

echo "3. For Local Minikube Deployment:"
echo "   - No secrets required"
echo "   - Use the 'Deploy to Minikube' workflow"
echo "   - Download artifacts and deploy locally"
echo

echo "4. Available Workflows:"
echo "   - test.yml: Runs tests and builds the application"
echo "   - deploy.yml: Full CI/CD pipeline with multiple environments"
echo "   - minikube-deploy.yml: Creates deployment package for local minikube"
echo

echo "5. Workflow Triggers:"
echo "   - Push to main/master: Automatic test and build"
echo "   - Pull Request: Automatic test and build"
echo "   - Manual trigger: Deploy to specific environment"
echo

echo "6. Container Registry:"
echo "   - Uses GitHub Container Registry (ghcr.io)"
echo "   - Images are automatically tagged with:"
echo "     - Branch name"
echo "     - Commit SHA"
echo "     - Semantic version (if using tags)"
echo

echo "📦 Generate Kubeconfig for Local Cluster:"
echo "=========================================="
echo

# Check if kubectl is available
if command -v kubectl &> /dev/null; then
    echo "Current kubectl context:"
    kubectl config current-context 2>/dev/null || echo "No context set"
    echo
    
    echo "To generate kubeconfig for GitHub Actions:"
    echo "kubectl config view --raw | base64 -w 0"
    echo
    
    echo "Copy the output and add it as KUBE_CONFIG secret in GitHub."
else
    echo "kubectl not found. Install kubectl first."
fi

echo
echo "🚀 Next Steps:"
echo "=============="
echo

echo "1. Push your code to GitHub:"
echo "   git add ."
echo "   git commit -m 'Add GitHub Actions workflows'"
echo "   git push origin main"
echo

echo "2. Check the Actions tab on GitHub to see the workflows running"
echo

echo "3. For local deployment:"
echo "   - Go to Actions > Deploy to Minikube"
echo "   - Click 'Run workflow'"
echo "   - Enter your PurpleAir device URL"
echo "   - Download the deployment package from artifacts"
echo

echo "4. For remote cluster deployment:"
echo "   - Set up the required secrets"
echo "   - Use the 'Build and Deploy Air Quality Monitor' workflow"
echo

echo "📚 Documentation:"
echo "================="
echo " - README.md: Application documentation"
echo " - KUBERNETES.md: Kubernetes deployment guide"
echo " - .github/workflows/: GitHub Actions workflows"
echo

echo "🎉 Setup instructions completed!"
echo
echo "Need help? Check the GitHub Actions documentation:"
echo "https://docs.github.com/en/actions"
