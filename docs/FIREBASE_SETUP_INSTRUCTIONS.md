# Firebase Admin SDK Setup Instructions

## Problem
Google Sign-In doesn't work because the backend needs Firebase Admin SDK credentials.

## Solution Steps

### 1. Go to Firebase Console
Visit: https://console.firebase.google.com

### 2. Select Your Project
- Click on "go-pro-platform" project
- Or create a new project if it doesn't exist

### 3. Enable Google Authentication
1. In the left sidebar, click **Authentication**
2. Click **Get Started** (if first time)
3. Click the **Sign-in method** tab
4. Under **Sign-in providers**, find **Google**
5. Click **Google** → Enable toggle → Save

### 4. Download Admin SDK Credentials
1. Click the **Settings gear icon** (⚙️) in the left sidebar
2. Select **Project settings**
3. Go to the **Service accounts** tab
4. Click **Generate new private key**
5. Click **Generate key** in the popup
6. A JSON file will download (e.g., `go-pro-platform-firebase-adminsdk-xxxxx.json`)

### 5. Install the Credentials
```bash
# Move the downloaded file to your backend config directory
cd ~/Desktop/FUN/go-pro/backend
mv ~/Downloads/go-pro-platform-firebase-adminsdk-*.json ./config/firebase-admin-sdk.json

# Verify it's there
ls -la config/firebase-admin-sdk.json
```

### 6. Verify Backend Configuration
Check that `backend/.env` contains:
```bash
FIREBASE_PROJECT_ID=go-pro-platform
FIREBASE_CREDENTIALS_PATH=./config/firebase-admin-sdk.json
```

### 7. Start Backend Server
```bash
cd backend
go run ./cmd/server
```

You should see logs indicating Firebase initialized successfully.

### 8. Test Google Sign-In
1. Start frontend: `cd frontend && bun run dev`
2. Go to http://localhost:3000/signin
3. Click "Google" button
4. Sign in with your Google account
5. You should be redirected to the dashboard

## Verification

### Test the Auth Flow
```bash
# 1. Get a Firebase ID token from the frontend (after signing in)
# Check browser console for the token

# 2. Test the backend verification endpoint
curl -X POST http://localhost:8080/api/v1/auth/verify \
  -H "Authorization: Bearer YOUR_FIREBASE_ID_TOKEN"

# Expected response:
{
  "user": {
    "id": "...",
    "firebase_uid": "...",
    "email": "you@example.com",
    "role": "admin"  // First user becomes admin
  }
}
```

## Security Notes

⚠️ **IMPORTANT**:
- The `firebase-admin-sdk.json` file contains sensitive credentials
- It's already in `.gitignore` - **NEVER commit it to git**
- Keep it secure and don't share it publicly
- Rotate credentials if accidentally exposed

## Troubleshooting

### "FIREBASE_PROJECT_ID not set"
- Ensure `backend/.env` has `FIREBASE_PROJECT_ID=go-pro-platform`
- Restart the backend server after editing .env

### "error initializing Firebase app: failed to create credentials"
- Verify `config/firebase-admin-sdk.json` exists
- Check file permissions: `chmod 600 config/firebase-admin-sdk.json`
- Ensure the path in .env is correct: `FIREBASE_CREDENTIALS_PATH=./config/firebase-admin-sdk.json`

### "Invalid or expired Firebase token"
- Frontend and backend must use the same Firebase project
- Check that `NEXT_PUBLIC_FIREBASE_PROJECT_ID` in frontend matches `FIREBASE_PROJECT_ID` in backend
- Both should be `go-pro-platform`

### Google Sign-In Popup Closes Immediately
- Ensure Google authentication is enabled in Firebase Console
- Check that your local URL is authorized in Firebase Console:
  - Go to Authentication → Settings → Authorized domains
  - Add `localhost` if not present

## Alternative: Use Environment Variables (Cloud Deployments)

For Cloud Run, GKE, or other Google Cloud environments, you can use default credentials:

```bash
# backend/.env
FIREBASE_PROJECT_ID=go-pro-platform
# Leave FIREBASE_CREDENTIALS_PATH empty or don't set it
```

The Firebase SDK will automatically use Application Default Credentials.

## Next Steps

After completing setup:
1. ✅ Sign in with Google should work
2. ✅ First user becomes admin automatically
3. ✅ Subsequent users get "student" role
4. ✅ Progress tracking syncs to backend
5. ✅ Admin can manage users at `/admin/users` (needs frontend route)

For questions, see:
- `backend/FIREBASE_AUTH_STATUS.md` - Complete auth status
- `backend/FIREBASE_AUTH_IMPLEMENTATION.md` - Detailed implementation guide
