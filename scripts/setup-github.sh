#!/bin/bash

# =================================================================
# GitHub Repository Setup Script for Blog CRM Management Microservice
# =================================================================

set -e

# Configuration
REPO_NAME="blog-crm-microservice"
GITHUB_ORG="MejonaTechnology"
REPO_URL="https://github.com/$GITHUB_ORG/$REPO_NAME.git"

echo "🚀 Setting up GitHub repository for Blog CRM Management Microservice..."

# Function to print colored output
print_status() {
    echo -e "\n\033[1;34m==>\033[0m $1"
}

print_success() {
    echo -e "\033[1;32m✓\033[0m $1"
}

print_error() {
    echo -e "\033[1;31m✗\033[0m $1"
}

# Check if git is installed
if ! command -v git &> /dev/null; then
    print_error "Git is not installed. Please install Git first."
    exit 1
fi

# Check if GitHub CLI is available
if command -v gh &> /dev/null; then
    print_status "GitHub CLI detected. Checking authentication..."
    if gh auth status &> /dev/null; then
        USE_GH_CLI=true
        print_success "GitHub CLI authenticated"
    else
        print_error "GitHub CLI not authenticated. Run 'gh auth login' first."
        USE_GH_CLI=false
    fi
else
    print_status "GitHub CLI not found. Will use manual setup."
    USE_GH_CLI=false
fi

# Initialize git repository if not already initialized
if [ ! -d ".git" ]; then
    print_status "Initializing Git repository..."
    git init
    print_success "Git repository initialized"
else
    print_success "Git repository already exists"
fi

# Create or update .gitignore if it doesn't exist
if [ ! -f ".gitignore" ]; then
    print_status "Creating .gitignore file..."
    cat > .gitignore << 'EOF'
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
*-linux
*-darwin
*-windows

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
vendor/

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Environment files
.env
.env.local
.env.production
.env.staging

# Log files
*.log
logs/

# Database files
*.db
*.sqlite
*.sqlite3

# Temporary files
tmp/
temp/

# Build directories
build/
dist/

# Node modules (if any JS tooling)
node_modules/

# Service files
service.log
*.pid

# AWS credentials
.aws/

# Docker files
.docker/

# Test coverage
coverage.out
coverage.html

# Documentation build
docs/_build/

# Database migration temp files
migrations/*.tmp
EOF
    print_success ".gitignore created"
fi

# Set up initial commit
print_status "Preparing initial commit..."

# Add all files except ignored ones
git add .

# Check if there are any changes to commit
if git diff --staged --quiet; then
    print_status "No changes to commit"
else
    # Commit initial files
    git commit -m "Initial commit: Blog CRM Management Microservice

- Complete Go microservice with 30+ API endpoints
- Advanced CRM features for blog lead generation
- Real-time analytics and performance tracking
- Revenue attribution and ROI calculation
- SEO optimization and content analysis
- JWT authentication with role-based access
- MySQL database with 78-column enhanced blogs table
- Docker and CI/CD deployment ready
- AWS EC2 production deployment scripts
- Comprehensive API documentation

Features:
✅ Blog CRUD operations with CRM integration
✅ Lead generation tracking and qualification
✅ Advanced analytics and reporting
✅ SEO performance analysis and recommendations
✅ Revenue attribution and conversion tracking
✅ Campaign management and A/B testing
✅ Social media integration and tracking
✅ Content optimization scoring
✅ Real-time performance monitoring
✅ Professional admin dashboard integration
✅ Health monitoring and metrics
✅ Authentication and authorization
✅ Database migrations and indexing
✅ Structured logging and audit trails
✅ CORS and security middleware
✅ Production-ready configuration

Architecture:
🏗️ Clean Architecture with Go microservice
📊 Enhanced database with 58 CRM fields
🔌 30+ RESTful API endpoints
💻 React admin dashboard integration
🚀 Docker containerization
⚙️ GitHub Actions CI/CD pipeline
☁️ AWS EC2 deployment automation
📈 Comprehensive monitoring and alerting

Built with ❤️ by Mejona Technology"
    print_success "Initial commit created"
fi

# Set up remote repository
if [ "$USE_GH_CLI" = true ]; then
    print_status "Creating GitHub repository using GitHub CLI..."
    
    # Create repository
    gh repo create "$GITHUB_ORG/$REPO_NAME" \
        --description "Enterprise-grade blog CRM microservice with advanced analytics, lead generation, and revenue attribution built with Go for Mejona Technology Admin Dashboard" \
        --homepage "https://mejona.in" \
        --public \
        --add-readme=false \
        --clone=false
    
    print_success "GitHub repository created"
    
    # Add remote
    git remote add origin "$REPO_URL" 2>/dev/null || git remote set-url origin "$REPO_URL"
    print_success "Remote origin added"
    
else
    print_status "Setting up remote manually..."
    echo ""
    echo "📋 Manual GitHub Setup Required:"
    echo "   1. Go to https://github.com/new"
    echo "   2. Repository name: $REPO_NAME"
    echo "   3. Description: Enterprise-grade blog CRM microservice with advanced analytics, lead generation, and revenue attribution built with Go for Mejona Technology Admin Dashboard"
    echo "   4. Set as Public repository"
    echo "   5. Do NOT initialize with README, .gitignore, or license"
    echo "   6. Create repository"
    echo ""
    read -p "Press Enter after creating the repository on GitHub..."
    
    # Add remote
    git remote add origin "$REPO_URL" 2>/dev/null || git remote set-url origin "$REPO_URL"
    print_success "Remote origin configured"
fi

# Set up main branch
print_status "Setting up main branch..."
git branch -M main
print_success "Main branch configured"

# Push to GitHub
print_status "Pushing to GitHub..."
git push -u origin main
print_success "Code pushed to GitHub"

# Create additional branches
print_status "Creating development branches..."
git checkout -b develop
git push -u origin develop

git checkout -b staging
git push -u origin staging

git checkout main
print_success "Development branches created"

# Set up repository topics and settings (if using GitHub CLI)
if [ "$USE_GH_CLI" = true ]; then
    print_status "Configuring repository settings..."
    
    # Add topics
    gh repo edit "$GITHUB_ORG/$REPO_NAME" \
        --add-topic go \
        --add-topic microservice \
        --add-topic api \
        --add-topic blog-crm \
        --add-topic analytics \
        --add-topic lead-generation \
        --add-topic revenue-attribution \
        --add-topic gin \
        --add-topic mysql \
        --add-topic jwt \
        --add-topic docker \
        --add-topic aws \
        --add-topic cicd \
        --add-topic mejona-technology
    
    print_success "Repository topics added"
fi

# Create initial GitHub issues
if [ "$USE_GH_CLI" = true ]; then
    print_status "Creating initial project issues..."
    
    gh issue create \
        --title "Set up production database with CRM enhancements" \
        --body "Configure production MySQL database with enhanced blogs table (78 columns) and CRM integration." \
        --label "enhancement,database,crm"
    
    gh issue create \
        --title "Implement comprehensive blog analytics monitoring" \
        --body "Set up advanced analytics dashboard with lead generation, conversion tracking, and revenue attribution metrics." \
        --label "enhancement,analytics,monitoring"
    
    gh issue create \
        --title "Configure SSL and domain for blog CRM service" \
        --body "Set up HTTPS with SSL certificates and configure domain routing for blog CRM microservice on port 8082." \
        --label "enhancement,security,deployment"
    
    gh issue create \
        --title "Implement automated testing for CRM features" \
        --body "Create comprehensive test suites for lead generation, analytics, and CRM functionality with CI integration." \
        --label "enhancement,testing,crm"
    
    gh issue create \
        --title "Set up blog lead nurturing automation" \
        --body "Implement automated lead qualification, scoring, and follow-up workflows based on blog engagement." \
        --label "feature,automation,crm"
    
    gh issue create \
        --title "Integrate with marketing automation tools" \
        --body "Connect blog CRM with email marketing platforms and CRM systems for seamless lead handoff." \
        --label "feature,integration,marketing"
    
    print_success "Initial issues created"
fi

echo ""
print_success "🎉 GitHub repository setup completed!"
echo ""
echo "📋 Repository Information:"
echo "   • Repository URL: $REPO_URL"
echo "   • Main Branch: main"
echo "   • Development Branch: develop"
echo "   • Staging Branch: staging"
echo ""
echo "🔗 Quick Links:"
echo "   • Repository: https://github.com/$GITHUB_ORG/$REPO_NAME"
echo "   • Issues: https://github.com/$GITHUB_ORG/$REPO_NAME/issues"
echo "   • Actions: https://github.com/$GITHUB_ORG/$REPO_NAME/actions"
echo "   • Releases: https://github.com/$GITHUB_ORG/$REPO_NAME/releases"
echo ""
echo "📝 Next Steps:"
echo "   1. Configure repository secrets for CI/CD deployment"
echo "   2. Set up branch protection rules for main/staging branches"
echo "   3. Configure AWS deployment credentials and environment variables"
echo "   4. Set up production database with enhanced CRM schema"
echo "   5. Configure blog CRM analytics dashboards and monitoring"
echo ""
echo "🚀 Blog CRM Features Ready for Deployment:"
echo "   • 30+ RESTful API endpoints with comprehensive CRM functionality"
echo "   • Advanced lead generation and qualification system"
echo "   • Revenue attribution and conversion tracking"
echo "   • SEO performance analysis and optimization"
echo "   • Real-time analytics and reporting dashboards"
echo "   • Campaign management and A/B testing framework"
echo "   • Professional React admin dashboard integration"
echo ""