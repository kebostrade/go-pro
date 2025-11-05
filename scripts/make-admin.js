#!/usr/bin/env node

/**
 * Make a user admin in Firebase
 *
 * Usage: node scripts/make-admin.js gcp.inspiration@gmail.com
 *
 * Requirements:
 * - Firebase Admin SDK installed: npm install firebase-admin
 * - Service account key (download from Firebase Console)
 */

const admin = require('firebase-admin');
const readline = require('readline');

// Configuration
const EMAIL = process.argv[2] || 'gcp.inspiration@gmail.com';

// Check if Firebase Admin is initialized
if (!admin.apps.length) {
  try {
    // Try to initialize with service account
    const serviceAccount = require('../firebase-service-account.json');
    admin.initializeApp({
      credential: admin.credential.cert(serviceAccount),
    });
    console.log('✓ Firebase Admin SDK initialized');
  } catch (error) {
    console.error('✗ Failed to initialize Firebase Admin SDK');
    console.error('');
    console.error('You need a service account key file.');
    console.error('');
    console.error('To get it:');
    console.error('1. Go to: https://console.firebase.google.com/project/go-pro-platform/settings/serviceaccounts/adminsdk');
    console.error('2. Click "Generate new private key"');
    console.error('3. Save as: firebase-service-account.json in project root');
    console.error('4. Run this script again');
    console.error('');
    process.exit(1);
  }
}

async function makeUserAdmin(email) {
  try {
    console.log(`\nLooking up user: ${email}...`);

    // Get user by email
    const userRecord = await admin.auth().getUserByEmail(email);
    console.log(`✓ Found user: ${userRecord.uid}`);

    // Create/Update user document in Firestore
    const userRef = admin.firestore().collection('users').doc(userRecord.uid);
    const userDoc = await userRef.get();

    if (userDoc.exists) {
      // Update existing user
      await userRef.update({
        role: 'admin',
        updatedAt: admin.firestore.FieldValue.serverTimestamp(),
      });
      console.log(`✓ Updated existing user to admin role`);
    } else {
      // Create new user document
      await userRef.set({
        uid: userRecord.uid,
        email: userRecord.email,
        displayName: userRecord.displayName || null,
        photoURL: userRecord.photoURL || null,
        emailVerified: userRecord.emailVerified,
        phoneNumber: userRecord.phoneNumber || null,
        role: 'admin',
        createdAt: admin.firestore.FieldValue.serverTimestamp(),
        lastLoginAt: admin.firestore.FieldValue.serverTimestamp(),
        progress: {
          completedLessons: [],
          xp: 0,
          level: 1,
        },
        preferences: {
          theme: 'system',
          notifications: true,
          language: 'en',
        },
        security: {
          mfaEnabled: false,
          phoneNumberVerified: false,
          linkedProviders: userRecord.providerData.map(p => p.providerId),
          accountLockout: false,
          failedLoginAttempts: 0,
        },
      });
      console.log(`✓ Created new user profile with admin role`);
    }

    // Optionally set custom claims
    const rl = readline.createInterface({
      input: process.stdin,
      output: process.stdout,
    });

    rl.question('\nDo you want to set admin custom claims? (y/n): ', async (answer) => {
      if (answer.toLowerCase() === 'y') {
        await admin.auth().setCustomUserClaims(userRecord.uid, {
          admin: true,
          role: 'admin',
        });
        console.log('✓ Set custom admin claims');
        console.log('  (User needs to sign out and back in for claims to take effect)');
      }

      console.log('\n✅ Done! User is now an admin.');
      console.log(`\nUser details:`);
      console.log(`  UID: ${userRecord.uid}`);
      console.log(`  Email: ${userRecord.email}`);
      console.log(`  Role: admin`);
      console.log(`\nNext steps:`);
      console.log(`  1. User should sign out and back in`);
      console.log(`  2. Access admin dashboard at: /admin/users`);
      console.log('');

      rl.close();
      process.exit(0);
    });

  } catch (error) {
    if (error.code === 'auth/user-not-found') {
      console.error(`\n✗ User not found: ${email}`);
      console.error('\nThe user needs to sign up first:');
      console.error('  1. Go to your app: http://localhost:3000/auth/signup');
      console.error('  2. Sign up with: ' + email);
      console.error('  3. Run this script again');
      console.error('');
    } else {
      console.error('\n✗ Error:', error.message);
    }
    process.exit(1);
  }
}

// Run
console.log('Firebase Admin - Make User Admin');
console.log('=================================');
makeUserAdmin(EMAIL);
