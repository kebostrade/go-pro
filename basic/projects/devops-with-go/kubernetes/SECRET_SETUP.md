# Kubernetes Secrets Setup

## Security Notice

**NEVER commit actual secrets to Git!**

The `secret.yaml` file has been removed and added to `.gitignore` for security.

## Creating Secrets

### Option 1: Using kubectl (Recommended)

```bash
# Generate strong passwords
POSTGRES_PASSWORD=$(openssl rand -base64 32)
JWT_SECRET=$(openssl rand -base64 64)

# Create secret directly in Kubernetes
kubectl create secret generic app-secrets \
  --from-literal=POSTGRES_USER=devops_user \
  --from-literal=POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
  --from-literal=JWT_SECRET="$JWT_SECRET" \
  --namespace=devops-demo

# Verify (values will be hidden)
kubectl get secret app-secrets -n devops-demo
```

### Option 2: Using secret.yaml (Local Development Only)

```bash
# Copy the example file
cp kubernetes/secret.yaml.example kubernetes/secret.yaml

# Generate base64 values
echo -n "devops_user" | base64
echo -n "$(openssl rand -base64 32)" | base64
echo -n "$(openssl rand -base64 64)" | base64

# Edit kubernetes/secret.yaml with the generated values
# DO NOT commit this file!
```

### Option 3: Using Sealed Secrets (Production)

```bash
# Install Sealed Secrets controller
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/controller.yaml

# Install kubeseal CLI
wget https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/kubeseal-0.24.0-linux-amd64.tar.gz
tar xfz kubeseal-0.24.0-linux-amd64.tar.gz
sudo install -m 755 kubeseal /usr/local/bin/kubeseal

# Create sealed secret
kubectl create secret generic app-secrets \
  --from-literal=POSTGRES_PASSWORD="$(openssl rand -base64 32)" \
  --from-literal=JWT_SECRET="$(openssl rand -base64 64)" \
  --dry-run=client -o yaml | \
  kubeseal -o yaml > kubernetes/sealed-secret.yaml

# Commit sealed-secret.yaml (it's encrypted and safe)
git add kubernetes/sealed-secret.yaml
```

## Rotating Secrets

```bash
# Update existing secret
kubectl create secret generic app-secrets \
  --from-literal=POSTGRES_PASSWORD="NEW_PASSWORD" \
  --from-literal=JWT_SECRET="NEW_JWT_SECRET" \
  --namespace=devops-demo \
  --dry-run=client -o yaml | kubectl apply -f -

# Restart pods to pick up new secrets
kubectl rollout restart deployment/devops-go-app -n devops-demo
```

## Security Best Practices

1. **Never commit unencrypted secrets to Git**
2. **Use strong, randomly generated passwords** (minimum 32 characters)
3. **Rotate secrets regularly** (every 90 days minimum)
4. **Use Sealed Secrets or External Secrets Operator** for production
5. **Limit RBAC permissions** to secrets
6. **Enable audit logging** for secret access
