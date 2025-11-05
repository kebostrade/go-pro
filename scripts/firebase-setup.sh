#!/bin/bash

# Firebase Setup Script for Go Pro Platform
# This script helps with initial Firebase setup and deployment

set -e

YELLOW='\033[1;33m'
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Firebase Setup for Go Pro Platform${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Check if Firebase CLI is installed
if ! command -v firebase &> /dev/null; then
    echo -e "${RED}Error: Firebase CLI is not installed${NC}"
    echo "Install it with: npm install -g firebase-tools"
    exit 1
fi

# Check if user is logged in
if ! firebase projects:list &> /dev/null; then
    echo -e "${YELLOW}You need to log in to Firebase${NC}"
    firebase login
fi

# Verify project
echo -e "${YELLOW}Current Firebase project:${NC}"
firebase use

echo ""
echo -e "${YELLOW}What would you like to do?${NC}"
echo "1. Deploy Firestore security rules"
echo "2. Deploy everything (rules + hosting)"
echo "3. Setup environment file"
echo "4. Enable authentication providers"
echo "5. View project info"
echo "6. Exit"
echo ""

read -p "Enter your choice (1-6): " choice

case $choice in
    1)
        echo -e "${YELLOW}Deploying Firestore security rules...${NC}"
        firebase deploy --only firestore:rules
        echo -e "${GREEN}✓ Firestore rules deployed successfully!${NC}"
        ;;
    2)
        echo -e "${YELLOW}Building frontend...${NC}"
        cd frontend
        npm run build
        cd ..

        echo -e "${YELLOW}Deploying to Firebase...${NC}"
        firebase deploy

        echo -e "${GREEN}✓ Deployment complete!${NC}"
        echo -e "${GREEN}Your app is live at: https://go-pro-platform.web.app${NC}"
        ;;
    3)
        echo -e "${YELLOW}Setting up environment file...${NC}"
        if [ ! -f frontend/.env.local ]; then
            cp frontend/.env.local.example frontend/.env.local
            echo -e "${GREEN}✓ Created frontend/.env.local from template${NC}"
            echo -e "${YELLOW}Note: All values are already configured for go-pro-platform${NC}"
        else
            echo -e "${YELLOW}frontend/.env.local already exists${NC}"
        fi
        ;;
    4)
        echo -e "${YELLOW}Opening Firebase Console to enable authentication providers...${NC}"
        echo ""
        echo "To enable authentication providers:"
        echo "1. Go to: https://console.firebase.google.com/project/go-pro-platform/authentication/providers"
        echo "2. Enable Google (click, toggle on, save)"
        echo "3. Enable GitHub:"
        echo "   - Create GitHub OAuth App: https://github.com/settings/developers"
        echo "   - Homepage URL: https://go-pro-platform.firebaseapp.com"
        echo "   - Callback URL: https://go-pro-platform.firebaseapp.com/__/auth/handler"
        echo "   - Copy Client ID and Secret to Firebase"
        echo ""
        read -p "Press Enter to open Firebase Console..."

        # Try to open in browser
        if command -v xdg-open &> /dev/null; then
            xdg-open "https://console.firebase.google.com/project/go-pro-platform/authentication/providers"
        elif command -v open &> /dev/null; then
            open "https://console.firebase.google.com/project/go-pro-platform/authentication/providers"
        else
            echo "Please manually open: https://console.firebase.google.com/project/go-pro-platform/authentication/providers"
        fi
        ;;
    5)
        echo -e "${YELLOW}Project Information:${NC}"
        echo ""
        echo "Project ID: go-pro-platform"
        echo "Project Number: 434643680939"
        echo "Web App ID: 1:434643680939:web:0ce27b5f6cda53789781ee"
        echo ""
        echo "URLs:"
        echo "  - Hosting: https://go-pro-platform.web.app"
        echo "  - Auth Domain: https://go-pro-platform.firebaseapp.com"
        echo "  - Console: https://console.firebase.google.com/project/go-pro-platform"
        echo ""
        echo "Enabled Services:"
        echo "  ✓ Firebase Authentication"
        echo "  ✓ Cloud Firestore"
        echo "  ✓ Firebase Hosting"
        echo "  ✓ Cloud Storage"
        ;;
    6)
        echo -e "${GREEN}Goodbye!${NC}"
        exit 0
        ;;
    *)
        echo -e "${RED}Invalid choice${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}Done!${NC}"
