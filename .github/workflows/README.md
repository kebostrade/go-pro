# GO-PRO CI/CD Pipeline Documentation

This directory contains GitHub Actions workflows for continuous integration and deployment of the GO-PRO platform.

## Overview

The CI/CD pipeline is designed to:
- Ensure code quality and security
- Run automated tests
- Build and publish Docker images
- Deploy to multiple environments
- Monitor infrastructure drift
- Keep dependencies up to date

## Workflows

### 1. Backend CI/CD (`backend-ci.yml`)

**Triggers**:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop`
- Changes in `backend/**` directory

**Jobs**:
1. **Lint** - Code formatting and linting
   - golangci-lint
   - gofmt check
   - go vet

2. **Security** - Security scanning
   - Gosec (SAST)
   - Trivy (vulnerability scanning)
   - SARIF upload to GitHub Security

3. **Test** - Unit testing
   - Race condition detection
   - Code coverage
   - Codecov upload

4. **Integration Test** - Integration testing
   - PostgreSQL service
   - Redis service
   - Full integration test suite

5. **Build** - Docker image build
   - Multi-platform (amd64, arm64)
   - Push to GitHub Container Registry
   - Layer caching

6. **Deploy Dev** - Deploy to development
   - Triggered on `develop` branch
   - AWS EKS deployment
   - Smoke tests

7. **Deploy Prod** - Deploy to production
   - Triggered on `main` branch
   - Blue-green deployment
   - Production smoke tests
   - Slack notifications

**Required Secrets**:
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `SLACK_WEBHOOK`

### 2. Microservices CI/CD (`microservices-ci.yml`)

**Triggers**:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop`
- Changes in `services/**` directory

**Jobs**:
1. **Detect Changes** - Identify changed services
   - Uses path filtering
   - Outputs changed services

2. **Shared Quality** - Test shared libraries
   - Lint and format
   - Unit tests

3. **Service CI** - Per-service testing
   - API Gateway
   - User Service
   - Course Service
   - Progress Service
   - Each with PostgreSQL and Redis

4. **Build Images** - Build Docker images
   - Matrix build for all services
   - Push to GitHub Container Registry
   - Layer caching

**Features**:
- Only builds changed services
- Parallel execution
- Service isolation

### 3. Terraform CI/CD (`terraform-ci.yml`)

**Triggers**:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop`
- Changes in `terraform/**` directory
- Schedule (drift detection)

**Jobs**:
1. **Validate** - Terraform validation
   - Format check
   - Validation
   - Initialization

2. **Security** - Infrastructure security
   - tfsec scanning
   - Checkov policy checks
   - SARIF upload

3. **Cost** - Cost estimation
   - Infracost analysis
   - PR comments with cost changes

4. **Plan Dev/Prod** - Terraform planning
   - Environment-specific plans
   - Plan artifacts

5. **Apply Dev/Prod** - Infrastructure deployment   - Auto-approve on main/dev
   - Output artifacts
   - Slack notifications

6. **Drift Detection** - Daily drift checks
   - Scheduled runs
   - Alerts on drift

**Required Secrets**:
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `INFRACOST_API_KEY`
- `SLACK_WEBHOOK`

### 4. Frontend CI/CD (`frontend-ci.yml`)

**Triggers**:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop`
- Changes in `frontend/**` directory

**Jobs**:
1. **Lint** - Code quality
   - ESLint
   - Prettier
   - TypeScript type checking

2. **Security** - Security scanning
   - npm audit
   - Snyk scanning

3. **Test** - Unit testing
   - Jest tests
   - Coverage reports
   - Codecov upload

4. **E2E** - End-to-end testing
   - Playwright tests
   - Test artifacts

5. **Build** - Application build
   - Next.js build
   - Build artifacts

6. **Docker Build** - Container image
   - Multi-stage build
   - Push to registry

7. **Deploy Vercel** - Vercel deployment
   - Development and production
   - Lighthouse CI
   - Performance monitoring

**Required Secrets**:
- `VERCEL_TOKEN`
- `VERCEL_ORG_ID`
- `VERCEL_PROJECT_ID`
- `SNYK_TOKEN`
- `API_URL`
- `SLACK_WEBHOOK`

### 5. Security Scanning (`security.yml`)

**Triggers**:
- Daily schedule (2 AM UTC)
- Push to `main` or `develop`
- Pull requests
- Manual dispatch

**Jobs**:
1. **CodeQL** - Code analysis
   - Go and JavaScript
   - Security queries
   - Quality queries

2. **Dependency Review** - Dependency scanning
   - PR dependency review
   - Vulnerability detection

3. **Secret Scan** - Secret detection
   - TruffleHog scanning
   - Historical commits

4. **Container Scan** - Image scanning
   - Trivy vulnerability scanner
   - All service images

5. **Semgrep** - SAST scanning
   - Multiple rulesets
   - Language-specific rules

6. **License Check** - License compliance
   - Go licenses
   - npm licenses

7. **SBOM** - Software Bill of Materials
   - Backend SBOM
   - Frontend SBOM

8. **Scorecard** - OpenSSF Scorecard
   - Security best practices
   - Supply chain security

9. **Notify** - Security alerts
   - Slack notifications
   - GitHub issues

**Required Secrets**:
- `SLACK_SECURITY_WEBHOOK`

### 6. Dependency Updates (`dependency-update.yml`)

**Triggers**:
- Weekly schedule (Monday 9 AM UTC)
- Manual dispatch

**Jobs**:
1. **Update Go Deps** - Go dependencies
   - Matrix for all modules
   - Automated PRs

2. **Update npm Deps** - npm dependencies
   - Frontend updates
   - Security fixes

3. **Update Terraform** - Terraform providers
   - Provider updates
   - Validation

4. **Update Actions** - GitHub Actions
   - Renovate bot
   - Action updates

**Features**:
- Automated pull requests
- Test validation
- Dependency labels

## Environment Configuration

### Development Environment

**Name**: `development`  
**URL**: https://dev-api.gopro.com  
**Branch**: `develop`

**Secrets**:
- AWS credentials for dev account
- Development database credentials
- Development API keys

### Production Environment

**Name**: `production`  
**URL**: https://api.gopro.com  
**Branch**: `main`

**Secrets**:
- AWS credentials for prod account
- Production database credentials
- Production API keys
- Slack webhook for notifications

**Protection Rules**:
- Required reviewers: 2
- Deployment approval required
- Branch protection enabled

## Required Secrets

### AWS Credentials
```
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
```

### Container Registry
```
GITHUB_TOKEN (automatically provided)
```

### Notifications
```
SLACK_WEBHOOK
SLACK_SECURITY_WEBHOOK
```

### External Services
```
VERCEL_TOKEN
VERCEL_ORG_ID
VERCEL_PROJECT_ID
SNYK_TOKEN
INFRACOST_API_KEY
```

### Application Configuration
```
API_URL
JWT_SECRET
DB_PASSWORD
REDIS_PASSWORD
```

## Workflow Best Practices

### 1. Branch Strategy
- `main` - Production deployments
- `develop` - Development deployments
- `feature/*` - Feature branches (no deployment)
- `hotfix/*` - Hotfix branches (fast-track to prod)

### 2. Pull Request Workflow
1. Create feature branch
2. Make changes
3. Push to GitHub
4. CI runs automatically
5. Review required
6. Merge to dev
7. Auto-deploy to dev
8. Test in dev
9. Merge to main
10. Auto-deploy to prod

### 3. Deployment Strategy
- **Development**: Automatic on merge to `develop`
- **Production**: Automatic on merge to `main` with approval
- **Rollback**: Revert commit and push

### 4. Security
- All secrets in GitHub Secrets
- SARIF results uploaded to Security tab
- Daily security scans
- Automated dependency updates

### 5. Monitoring
- Codecov for test coverage
- Lighthouse for performance
- Infracost for infrastructure costs
- Slack notifications for deployments

## Troubleshooting

### Workflow Fails on Lint
```bash
# Fix locally
cd backend
golangci-lint run --fix
go fmt ./...
```

### Tests Fail
```bash
# Run tests locally
cd backend
go test ./...

# With coverage
go test -cover ./...
```

### Docker Build Fails
```bash
# Test build locally
docker build -t test ./backend

# Check Dockerfile
docker build --no-cache -t test ./backend
```

### Deployment Fails
1. Check AWS credentials
2. Verify EKS cluster access
3. Check deployment logs
4. Verify image exists in registry

### Security Scan Fails
1. Review security findings
2. Update dependencies
3. Fix vulnerabilities
4. Re-run scan

## Maintenance

### Weekly Tasks
- Review dependency update PRs
- Check security scan results
- Monitor deployment success rate

### Monthly Tasks
- Review and update workflows
- Update action versions
- Review secret rotation
- Check cost reports

### Quarterly Tasks
- Security audit
- Performance review
- Infrastructure optimization
- Documentation updates

## Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Terraform CI/CD](https://www.terraform.io/docs/cloud/guides/recommended-practices/index.html)
- [Go Testing](https://golang.org/pkg/testing/)
- [Next.js Deployment](https://nextjs.org/docs/deployment)

