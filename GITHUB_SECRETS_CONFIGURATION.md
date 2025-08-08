# 🔐 GitHub Repository Secrets Configuration

**Repository**: `MejonaTechnology/blog-crm-microservice`  
**URL**: https://github.com/MejonaTechnology/blog-crm-microservice/settings/secrets/actions

## Required Secrets for CI/CD Pipeline

### 🌐 AWS Infrastructure Secrets

#### `PRODUCTION_HOST`
- **Value**: `65.1.94.25`
- **Description**: Production AWS EC2 server IP address for blog service deployment

#### `STAGING_HOST` 
- **Value**: `65.1.94.25` (same server, different deployment path)
- **Description**: Staging environment AWS EC2 server IP address

#### `EC2_USERNAME`
- **Value**: `ubuntu`
- **Description**: SSH username for AWS EC2 instance access

#### `EC2_SSH_KEY`
- **Description**: Private SSH key for AWS EC2 access
- **Format**: Complete private key including `-----BEGIN OPENSSH PRIVATE KEY-----` headers
- **Note**: This is the private key corresponding to the EC2 key pair

### 🗄️ Database Configuration Secrets

#### `DB_HOST`
- **Value**: `65.1.94.25`
- **Description**: MySQL database server host (same as EC2 instance)

#### `DB_USER`
- **Value**: `phpmyadmin`
- **Description**: MySQL database username (from production environment)

#### `DB_PASSWORD`
- **Value**: `mFVarH2LCrQK`
- **Description**: MySQL database password (from production .env)

#### `DB_NAME`
- **Value**: `mejona_unified`
- **Description**: Primary database name with enhanced blogs table (78 columns)

### 🔑 Authentication & Security Secrets

#### `JWT_SECRET`
- **Description**: JWT secret key for token signing and validation
- **Recommendation**: Generate a secure 256-bit key
- **Example**: `your-super-secure-jwt-secret-key-here-256-bits`

### 🐳 Docker Registry Secrets (Optional)

#### `DOCKER_USERNAME`
- **Value**: `mejonatechnology`
- **Description**: Docker Hub username for image publishing

#### `DOCKER_PASSWORD`  
- **Description**: Docker Hub access token or password
- **Note**: Required only if pushing Docker images to registry

### 📢 Notification Secrets (Optional)

#### `SLACK_WEBHOOK`
- **Description**: Slack webhook URL for deployment notifications
- **Format**: `https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX`
- **Note**: Required only if Slack notifications are desired

## 🚀 Configuration Steps

### Step 1: Navigate to Repository Settings
Go to: https://github.com/MejonaTechnology/blog-crm-microservice/settings/secrets/actions

### Step 2: Add Required Secrets
Click "New repository secret" for each secret listed above

### Step 3: Priority Order
Configure in this order for immediate deployment:
1. `PRODUCTION_HOST` = `65.1.94.25`
2. `EC2_USERNAME` = `ubuntu`  
3. `EC2_SSH_KEY` = (private key content)
4. `DB_HOST` = `65.1.94.25`
5. `DB_USER` = `phpmyadmin`
6. `DB_PASSWORD` = `mFVarH2LCrQK`
7. `DB_NAME` = `mejona_unified`
8. `JWT_SECRET` = (generate secure key)

### Step 4: Verify Configuration
After adding secrets, trigger a new workflow run or push a commit to test the deployment

## ✅ Current Status

- **Repository**: ✅ Created successfully
- **Code**: ✅ Pushed to GitHub (47 files)
- **CI/CD Pipeline**: ✅ Workflow file deployed
- **Secrets**: ⏳ **REQUIRES MANUAL CONFIGURATION**
- **Deployment**: ⏳ Pending secrets configuration

## 🔍 Verification

After configuring secrets, the workflow will:
1. ✅ Build and test the Go microservice
2. ✅ Create Docker image (if Docker secrets provided)
3. ✅ Deploy to AWS EC2 production server
4. ✅ Configure systemd service on port 8082
5. ✅ Set up nginx reverse proxy
6. ✅ Run health checks and smoke tests
7. ✅ Send deployment notifications (if Slack configured)

## 🏁 Expected Results

Once secrets are configured:
- **Service URL**: http://65.1.94.25:8082/health
- **API Endpoint**: http://65.1.94.25:8082/api/v1/test
- **Admin Integration**: Port 8082 ready for admin dashboard
- **Database**: Connected to enhanced blogs table (78 columns)
- **CRM Features**: Lead generation, analytics, revenue attribution

---

**Next Action**: Configure the secrets manually at the GitHub repository settings page, then monitor the workflow execution in the Actions tab.