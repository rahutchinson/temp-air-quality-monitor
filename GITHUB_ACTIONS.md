# GitHub Actions Deployment Guide

This guide explains how to use GitHub Actions to automatically build, test, and deploy your Air Quality Monitor application to Kubernetes.

## Overview

The GitHub Actions workflows provide:

- **Automated Testing**: Run tests on every push and pull request
- **Docker Image Building**: Build and push images to GitHub Container Registry
- **Multi-Environment Deployment**: Deploy to development, staging, and production
- **Local Deployment Package**: Create packages for local minikube deployment

## Available Workflows

### 1. Test and Build (`test.yml`)

**Triggers**: Push to main/master/develop, Pull Requests

**What it does**:
- Runs Go tests
- Builds the application
- Builds Docker image
- Uploads build artifacts

**No configuration required** - runs automatically on code changes.

### 2. Build and Deploy (`deploy.yml`)

**Triggers**: Push to main/master, Manual trigger

**What it does**:
- Runs tests
- Builds and pushes Docker image to GitHub Container Registry
- Deploys to Kubernetes clusters (development, staging, production)

**Requires**: Kubernetes cluster configuration secrets

### 3. Deploy to Minikube (`minikube-deploy.yml`)

**Triggers**: Manual trigger only

**What it does**:
- Builds the application
- Creates a deployment package
- Generates deployment instructions
- Uploads everything as artifacts

**No secrets required** - perfect for local deployment

## Setup Instructions

### Step 1: Enable GitHub Actions

1. Push your code to GitHub
2. Go to your repository
3. Click on the "Actions" tab
4. Click "Enable Actions" if prompted

### Step 2: Configure Secrets (for remote deployment)

If you want to deploy to remote Kubernetes clusters, you'll need to set up secrets:

1. Go to your repository Settings
2. Navigate to "Secrets and variables" > "Actions"
3. Add the following secrets:

#### For Development/Testing:
```
KUBE_CONFIG: <base64-encoded-kubeconfig>
```

#### For Staging:
```
STAGING_KUBE_CONFIG: <base64-encoded-kubeconfig>
```

#### For Production:
```
PRODUCTION_KUBE_CONFIG: <base64-encoded-kubeconfig>
```

### Step 3: Generate Kubeconfig

To generate the base64-encoded kubeconfig:

```bash
# For your local cluster
kubectl config view --raw | base64 -w 0

# Copy the output and add it as a secret in GitHub
```

## Usage Scenarios

### Scenario 1: Local Development with Minikube

1. **Trigger the workflow**:
   - Go to Actions > "Deploy to Minikube"
   - Click "Run workflow"
   - Enter your PurpleAir device URL (e.g., `http://192.168.0.249/json`)
   - Click "Run workflow"

2. **Download and deploy**:
   - Wait for the workflow to complete
   - Download the "deployment-package" artifact
   - Extract and follow the instructions in "deployment-instructions.md"

3. **Deploy locally**:
   ```bash
   # Extract the package
   tar -xzf deployment-package.tar.gz
   cd deployment-package
   
   # Set up minikube
   ./setup-k8s.sh
   
   # Deploy the application
   cd k8s
   ./deploy.sh
   ```

### Scenario 2: Automated Deployment to Remote Cluster

1. **Set up secrets** (see Step 2 above)

2. **Deploy automatically**:
   - Push to main/master branch triggers automatic deployment
   - Or manually trigger the "Build and Deploy Air Quality Monitor" workflow

3. **Monitor deployment**:
   - Check the Actions tab for deployment status
   - Use kubectl to verify the deployment

### Scenario 3: Multi-Environment Deployment

1. **Set up all required secrets**

2. **Deploy to specific environment**:
   - Go to Actions > "Build and Deploy Air Quality Monitor"
   - Click "Run workflow"
   - Select the target environment (development, staging, production)
   - Click "Run workflow"

## Workflow Details

### Test Workflow (`test.yml`)

```yaml
# Triggers on:
- push to main/master/develop
- pull requests

# Jobs:
- test: Runs Go tests and builds application
```

### Deploy Workflow (`deploy.yml`)

```yaml
# Triggers on:
- push to main/master
- manual trigger

# Jobs:
- test: Runs tests
- build: Builds and pushes Docker image
- deploy-development: Deploys to development cluster
- deploy-staging: Deploys to staging cluster (manual)
- deploy-production: Deploys to production cluster (manual)
```

### Minikube Deploy Workflow (`minikube-deploy.yml`)

```yaml
# Triggers on:
- manual trigger only

# Jobs:
- deploy: Creates deployment package and instructions
```

## Container Registry

The workflows use GitHub Container Registry (ghcr.io) to store Docker images.

**Image naming**: `ghcr.io/username/repository:tag`

**Automatic tags**:
- `main` or `master`: Latest from main branch
- `sha-abc123`: Specific commit
- `v1.0.0`: Semantic version tags

## Environment Variables

The workflows automatically set these environment variables:

- `REGISTRY`: `ghcr.io`
- `IMAGE_NAME`: Your repository name
- `DEVICE_URL`: From workflow input or config

## Troubleshooting

### Common Issues

1. **Workflow fails on test**:
   - Check that all Go tests pass locally
   - Ensure all dependencies are in go.mod

2. **Docker build fails**:
   - Check Dockerfile syntax
   - Ensure all required files are present

3. **Deployment fails**:
   - Verify Kubernetes cluster is accessible
   - Check that secrets are properly configured
   - Ensure cluster has sufficient resources

4. **Image pull fails**:
   - Check that GitHub Container Registry is accessible
   - Verify image tags are correct

### Debugging

1. **Check workflow logs**:
   - Go to Actions tab
   - Click on the failed workflow
   - Check the logs for each step

2. **Test locally**:
   - Use the minikube deployment workflow
   - Test the deployment package locally

3. **Verify secrets**:
   - Ensure kubeconfig is properly base64 encoded
   - Test cluster connectivity manually

## Best Practices

### Security

1. **Use secrets for sensitive data**:
   - Never commit kubeconfig files
   - Use GitHub secrets for all credentials

2. **Limit permissions**:
   - Use least-privilege access for Kubernetes clusters
   - Regularly rotate credentials

### Reliability

1. **Test before deploy**:
   - All workflows run tests first
   - Use staging environment for testing

2. **Monitor deployments**:
   - Check pod status after deployment
   - Monitor application logs

3. **Rollback capability**:
   - Keep previous image versions
   - Use Kubernetes rollback commands if needed

### Performance

1. **Optimize Docker builds**:
   - Use multi-stage builds
   - Leverage build cache

2. **Parallel jobs**:
   - Workflows run jobs in parallel where possible
   - Use dependencies to ensure proper order

## Advanced Configuration

### Customizing Workflows

You can customize the workflows by editing the YAML files:

1. **Add new environments**:
   - Copy the deploy job
   - Update environment names and secrets

2. **Modify build process**:
   - Add additional build steps
   - Change Docker build arguments

3. **Add notifications**:
   - Configure Slack/email notifications
   - Add status checks

### Conditional Deployment

Use workflow conditions to control deployment:

```yaml
# Only deploy on main branch
if: github.ref == 'refs/heads/main'

# Only deploy on specific tags
if: startsWith(github.ref, 'refs/tags/v')
```

### Environment Protection

For production deployments:

1. **Add environment protection rules**:
   - Require manual approval
   - Add status checks
   - Restrict who can deploy

2. **Use deployment strategies**:
   - Rolling updates
   - Blue-green deployment
   - Canary deployments

## Support

For issues with GitHub Actions:

1. Check the [GitHub Actions documentation](https://docs.github.com/en/actions)
2. Review workflow logs for error messages
3. Test workflows locally when possible
4. Use the minikube deployment for local testing

## Examples

### Example 1: Quick Local Deployment

```bash
# 1. Push code to GitHub
git push origin main

# 2. Trigger minikube deployment workflow
# 3. Download deployment package
# 4. Deploy locally
./setup-k8s.sh
cd k8s && ./deploy.sh
```

### Example 2: Production Deployment

```bash
# 1. Set up production secrets
# 2. Create a release tag
git tag v1.0.0
git push origin v1.0.0

# 3. Trigger production deployment workflow
# 4. Monitor deployment status
```

### Example 3: Continuous Deployment

```bash
# 1. Set up all environments
# 2. Push to main branch
git push origin main

# 3. Automatic deployment to development
# 4. Manual promotion to staging/production
```
