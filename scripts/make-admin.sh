#!/bin/bash

# Make a user admin in Firebase
# Usage: ./scripts/make-admin.sh <email>

set -e

EMAIL="${1:-gcp.inspiration@gmail.com}"

echo "Making $EMAIL an admin user..."

# This script helps create/update admin users
# Since we don't have direct Admin SDK access, we'll use Firebase CLI

# Option 1: Using Firebase Auth to get UID, then update Firestore
echo ""
echo "Method 1: Direct Firestore Update"
echo "=================================="
echo ""
echo "Since the user is already authenticated, you can update their role manually:"
echo ""
echo "1. Go to Firebase Console:"
echo "   https://console.firebase.google.com/project/go-pro-platform/firestore/data"
echo ""
echo "2. Find/Create the 'users' collection"
echo ""
echo "3. Find the document with email: $EMAIL"
echo "   OR create a new document if they haven't signed in yet"
echo ""
echo "4. Set the 'role' field to: 'admin'"
echo ""
echo "=================================="
echo ""
echo "Method 2: Sign in to the app first"
echo "=================================="
echo ""
echo "1. Have $EMAIL sign in to your app"
echo "2. This will automatically create their Firestore profile"
echo "3. Then run this script again to update their role"
echo ""
echo "OR use the admin dashboard once you're admin!"
echo ""
echo "Method 3: Use Firebase Admin SDK (Backend)"
echo "=========================================="
echo ""
echo "Create a Cloud Function or backend script:"
echo ""
cat << 'FUNCTION'
import * as admin from 'firebase-admin';

admin.initializeApp();

async function makeAdmin(email: string) {
  const user = await admin.auth().getUserByEmail(email);
  await admin.firestore().collection('users').doc(user.uid).set({
    uid: user.uid,
    email: user.email,
    displayName: user.displayName,
    photoURL: user.photoURL,
    emailVerified: user.emailVerified,
    role: 'admin',
    createdAt: admin.firestore.FieldValue.serverTimestamp(),
    lastLoginAt: admin.firestore.FieldValue.serverTimestamp(),
  }, { merge: true });

  console.log(`${email} is now an admin!`);
}

makeAdmin('gcp.inspiration@gmail.com');
FUNCTION

echo ""
echo "For now, the easiest way is Method 1 (Firebase Console)."
echo ""
